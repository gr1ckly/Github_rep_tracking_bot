package notification_service

import (
	"Common"
	"TeleRequestHandler/bot"
	"TeleRequestHandler/notification_service/converters"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"net"
	"os"
	"strconv"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type KafkaNotificationWaiter struct {
	reader         *kafka.Reader
	bot            bot.Bot[any, string, int64]
	messagePattern string
}

func NewKafkaNotificationWaiter(network string, addr string, topicName string, topicPartition int, topicReplicationFactor int, groupId string, messagePattern string, bot bot.Bot[any, string, int64]) (*KafkaNotificationWaiter, error) {
	conn, err := kafka.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	partition, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}
	needToCreate := true
	for _, part := range partition {
		if part.Topic == topicName {
			needToCreate = false
		}
	}
	if needToCreate {
		controller, err := conn.Controller()
		if err != nil {
			return nil, err
		}
		var controllerConn *kafka.Conn
		controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			return nil, err
		}
		defer controllerConn.Close()
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topicName,
				NumPartitions:     topicPartition,
				ReplicationFactor: topicReplicationFactor,
			},
		}
		err = controllerConn.CreateTopics(topicConfigs...)
		if err != nil {
			return nil, err
		}
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{addr},
		Topic:          topicName,
		GroupID:        groupId,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
		StartOffset:    kafka.FirstOffset,
	})
	return &KafkaNotificationWaiter{reader, bot, messagePattern}, nil
}

func (ns *KafkaNotificationWaiter) WaitNotification(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := ns.reader.FetchMessage(ctx)
			if err != nil {
				logger.Error(err.Error())
				if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
					return
				}
				continue
			}
			var dto Common.ChangingDTO
			err = json.Unmarshal(msg.Value, &dto)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			err = ns.processMessage(dto)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			err = ns.reader.CommitMessages(ctx, msg)
			if err != nil {
				logger.Error(err.Error())
				if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
					return
				}
				continue
			}
		}
	}
}

func (ns *KafkaNotificationWaiter) processMessage(dto Common.ChangingDTO) error {
	msg := converters.ConvertChanging(ns.messagePattern, dto)
	err := ns.bot.SendMessage(dto.ChatId, msg)
	return err
}

func (ns *KafkaNotificationWaiter) Close() error {
	if ns.reader != nil {
		err := ns.reader.Close()
		ns.reader = nil
		return err
	}
	return nil
}
