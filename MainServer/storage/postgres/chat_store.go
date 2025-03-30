package postgres

import (
	"Crypto_Bot/MainServer/storage"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	ADD_CHAT         = `INSERT INTO CHAT(CHAT_ID, TYPE) VALUES ($1, $2);`
	REMOVE_CHAT      = `DELETE FROM CHAT WHERE CHAT_ID = $1`
	GET_BY_ID_CHAT   = `SELECT * FROM CHAT WHERE CHAT_ID = $1`
	GET_CHATS_OFFSET = `SELECT * FROM CHAT order by CHAT_ID OFFSET $1 LIMIT $2;`
	GET_CHAT_NUMBER  = `SELECT COUNT(*) AS TOTAL_ROWS FROM CHAT;`
	UPDATE_CHAT      = `UPDATE CHAT
SET TYPE = $1
WHERE CHAT_ID = $2;`
)

const (
	ADD_CHAT_NAME         = "add_chat"
	REMOVE_CHAT_NAME      = "remove_chat"
	GET_BY_ID_CHAT_NAME   = "get_by_id_chat"
	GET_CHATS_OFFSET_NAME = "get_chats_offset"
	GET_CHAT_NUMBER_NAME  = "get_chat_number"
	UPDATE_CHAT_NAME      = "update_chat"
)

var chatStatementMap = map[string]string{
	ADD_CHAT_NAME:         ADD_CHAT,
	REMOVE_CHAT_NAME:      REMOVE_CHAT,
	GET_BY_ID_CHAT_NAME:   GET_BY_ID_CHAT,
	GET_CHATS_OFFSET_NAME: GET_CHATS_OFFSET,
	GET_CHAT_NUMBER_NAME:  GET_CHAT_NUMBER,
	UPDATE_CHAT_NAME:      UPDATE_CHAT,
}

type PostgresChatStore struct {
	pool    *pgxpool.Pool
	timeout int
}

func NewPostgresChatStore(timeout int, dbUrl string) (*PostgresChatStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	initCtx, initCancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer initCancel()
	err = pool.AcquireFunc(initCtx, func(conn *pgxpool.Conn) error {
		for key, _ := range chatStatementMap {
			_, stmtErr := conn.Conn().Prepare(context.Background(), key, chatStatementMap[key])
			if stmtErr != nil {
				return stmtErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &PostgresChatStore{pool: pool, timeout: timeout}, nil
}

func (pc *PostgresChatStore) AddNewChat(chat *storage.Chat) (int, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return -1, nil
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, ADD_CHAT_NAME, chat.ChatID, chat.Type)
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return chat.ChatID, err
}

func (pc *PostgresChatStore) RemoveChat(id int) error {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return nil
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, REMOVE_CHAT_NAME, id)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (pc *PostgresChatStore) GetChatByID(id int) (*storage.Chat, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return nil, nil
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	chat := storage.Chat{}
	err = tx.QueryRow(ctx, GET_BY_ID_CHAT_NAME, id).Scan(&chat.ChatID, &chat.Type)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &chat, err
}

func (pc *PostgresChatStore) GetChatsOffset(start int, limit int) ([]storage.Chat, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return nil, nil
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_CHATS_OFFSET_NAME, start, limit)
	if err != nil {
		return nil, err
	}
	answer := []storage.Chat{}
	for row.Next() {
		chat := storage.Chat{}
		err = row.Scan(&chat.ChatID, &chat.Type)
		if err != nil {
			return nil, err
		}
		answer = append(answer, chat)
	}
	err = tx.Commit(ctx)
	return answer, err
}

func (pc *PostgresChatStore) GetChatNumber() (int, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return -1, nil
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var number int
	err = tx.QueryRow(ctx, GET_CHAT_NUMBER_NAME).Scan(&number)
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return number, err
}

func (pc *PostgresChatStore) UpdateChat(chat *storage.Chat) error {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return nil
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, UPDATE_CHAT_NAME, chat.Type, chat.ChatID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (pc *PostgresChatStore) Close() {
	if pc.pool != nil {
		pc.pool.Close()
		pc.pool = nil
	}
}
