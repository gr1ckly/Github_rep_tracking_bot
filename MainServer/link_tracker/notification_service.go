package link_tracker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"time"
)

type KafkaNotificationService struct {
	writer  *kafka.Writer
	timeout int
}

func NewNotificationService(timeout int, network string, addr string, topicName string, topicPartition int, topicReplicationFactor int) (*KafkaNotificationService, error) {
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
	writer := &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaNotificationService{writer, timeout}, nil
}

func (ns *KafkaNotificationService) Notify(msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ns.timeout)*time.Second)
	defer cancel()
	return ns.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Message"),
		Value: data,
	})
}

func (ns *KafkaNotificationService) Close() error {
	if ns.writer != nil {
		return ns.writer.Close()
	}
	return nil
}
