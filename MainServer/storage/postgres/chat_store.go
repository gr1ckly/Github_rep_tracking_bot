package postgres

import (
	"Crypto_Bot/MainServer/custom_errors"
	"Crypto_Bot/MainServer/storage"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	return &PostgresChatStore{pool: pool, timeout: timeout}, nil
}

func (pc *PostgresChatStore) AddNewChat(chat *storage.Chat) (int64, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, ADD_CHAT, chat.ChatID, chat.Type)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return -1, custom_errors.NewAlreadyExistsError(err)
		}
	}
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return chat.ChatID, err
}

func (pc *PostgresChatStore) RemoveChat(id int64) error {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, REMOVE_CHAT, id)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (pc *PostgresChatStore) GetChatByID(id int64) (*storage.Chat, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	chat := storage.Chat{}
	err = tx.QueryRow(ctx, GET_BY_ID_CHAT, id).Scan(&chat.ChatID, &chat.Type)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &chat, err
}

func (pc *PostgresChatStore) GetChatsOffset(start int, limit int) ([]storage.Chat, error) {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_CHATS_OFFSET, start, limit)
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
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var number int
	err = tx.QueryRow(ctx, GET_CHAT_NUMBER).Scan(&number)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return number, err
}

func (pc *PostgresChatStore) UpdateChat(chat *storage.Chat) error {
	conn, err := pc.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pc.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, UPDATE_CHAT, chat.Type, chat.ChatID)
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
