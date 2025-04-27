package postgres

import (
	"Crypto_Bot/MainServer/custom_errors"
	"Crypto_Bot/MainServer/storage"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	ADD_CHAT_REPO_RECORD         = `INSERT INTO CHAT_REPO_RECORD(CHAT, REPO, TAGS, EVENTS) VALUES($1, $2, $3, $4) RETURNING ID;`
	REMOVE_CHAT_REPO_RECORD      = `DELETE FROM CHAT_REPO_RECORD WHERE CHAT=$1 and REPO = $2;`
	GET_BY_CHAT_CHAT_REPO_RECORD = `SELECT cr.ID AS chat_repo_record_id,
       c.CHAT_ID,
       c.TYPE AS chat_type,
       r.ID AS repo_id,
       r.NAME AS repo_name,
       r.OWNER AS repo_owner,
       r.LINK,
       r.LAST_COMMIT,
       r.LAST_ISSUE,
       r.LAST_PULL_REQUEST,
       cr.TAGS,
       cr.EVENTS
FROM CHAT_REPO_RECORD cr
JOIN CHAT c ON cr.chat = c.CHAT_ID
JOIN REPO r ON cr.REPO = r.ID
WHERE c.CHAT_ID = $1;`
	GET_BY_ID_CHAT_REPO_RECORD = `SELECT cr.ID AS chat_repo_record_id,
       c.CHAT_ID,
       c.TYPE AS chat_type,
       r.ID AS repo_id,
       r.NAME AS repo_name,
       r.OWNER AS repo_owner,
       r.LINK,
       r.LAST_COMMIT,
       r.LAST_ISSUE,
       r.LAST_PULL_REQUEST,
       cr.TAGS,
       cr.EVENTS
FROM CHAT_REPO_RECORD cr
JOIN CHAT c ON cr.chat = c.CHAT_ID
JOIN REPO r ON cr.REPO = r.ID
WHERE cr.ID = $1;`
	GET_BY_CHAT_AND_LINK_CHAT_REPO_RECORD = `SELECT cr.ID AS chat_repo_record_id,
       c.CHAT_ID,
       c.TYPE AS chat_type,
       r.ID AS repo_id,
       r.NAME AS repo_name,
       r.OWNER AS repo_owner,
       r.LINK,
       r.LAST_COMMIT,
       r.LAST_ISSUE,
       r.LAST_PULL_REQUEST,
       cr.TAGS,
       cr.EVENTS
FROM CHAT_REPO_RECORD cr
JOIN CHAT c ON cr.chat = c.CHAT_ID
JOIN REPO r ON cr.REPO = r.ID
WHERE  c.CHAT_ID = $1 AND r.ID=$2;`
	GET_BY_LINK_ID_CHAT_REPO_RECORD = `SELECT cr.ID AS chat_repo_record_id,
       c.CHAT_ID,
       c.TYPE AS chat_type,
       r.ID AS repo_id,
       r.NAME AS repo_name,
       r.OWNER AS repo_owner,
       r.LINK,
       r.LAST_COMMIT,
       r.LAST_ISSUE,
       r.LAST_PULL_REQUEST,
       cr.TAGS,
       cr.EVENTS
FROM CHAT_REPO_RECORD cr
JOIN CHAT c ON cr.chat = c.CHAT_ID
JOIN REPO r ON cr.REPO = r.ID
WHERE r.ID = $1;`
	GET_OFFSET_CHAT_REPO_RECORD = `SELECT cr.ID AS chat_repo_record_id,
       c.CHAT_ID,
       c.TYPE AS chat_type,
       r.ID AS repo_id,
       r.NAME AS repo_name,
       r.OWNER AS repo_owner,
       r.LINK,
       r.LAST_COMMIT,
       r.LAST_ISSUE,
       r.LAST_PULL_REQUEST,
       cr.TAGS,
       cr.EVENTS
FROM CHAT_REPO_RECORD cr
JOIN CHAT c ON cr.chat = c.CHAT_ID
JOIN REPO r ON cr.REPO = r.ID
ORDER BY chat_repo_record_id
OFFSET $1 LIMIT $2;`
	GET_NUMBER_CHAT_REPO_RECORD = `
SELECT COUNT(*) AS TOTAL_ROWS FROM CHAT_REPO_RECORD;
`
	UPDATE_CHAT_REPO_RECORD = `
UPDATE CHAT_REPO_RECORD
SET chat = $1,
    repo = $2,
    tags = $3,
    events = $4
WHERE id = $5;
`
)

type PostgresChatRepoRecordStore struct {
	pool    *pgxpool.Pool
	timeout int
}

func NewPostgresChatRepoRecordStore(timeout int, dbUrl string) (*PostgresChatRepoRecordStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return &PostgresChatRepoRecordStore{pool: pool, timeout: timeout}, nil
}

func (ps *PostgresChatRepoRecordStore) AddNewRecord(record *storage.ChatRepoRecord) (int, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var id int
	err = tx.QueryRow(ctx, ADD_CHAT_REPO_RECORD, record.Chat.ChatID, record.Repo.ID, record.Tags, record.Events).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return id, err
}

func (ps *PostgresChatRepoRecordStore) RemoveRecord(chat_id int64, repo_id int) (int64, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	res, err := tx.Exec(ctx, REMOVE_CHAT_REPO_RECORD, chat_id, repo_id)
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return res.RowsAffected(), err
}

func (ps *PostgresChatRepoRecordStore) GetRecordByChat(chat_id int64) ([]storage.ChatRepoRecord, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_BY_CHAT_CHAT_REPO_RECORD, chat_id)
	if err != nil {
		return nil, err
	}
	answer := []storage.ChatRepoRecord{}
	for row.Next() {
		record := storage.ChatRepoRecord{Chat: &storage.Chat{}, Repo: &storage.Repo{}}
		err = row.Scan(&record.ID, &record.Chat.ChatID, &record.Chat.Type, &record.Repo.ID, &record.Repo.Name, &record.Repo.Owner, &record.Repo.Link,
			&record.Repo.LastCommit, &record.Repo.LastIssue, &record.Repo.LastPR, &record.Tags, &record.Events)
		if err != nil {
			return nil, err
		}
		answer = append(answer, record)
	}
	err = tx.Commit(ctx)
	return answer, err
}

func (ps *PostgresChatRepoRecordStore) GetRecordByLink(link_id int) ([]storage.ChatRepoRecord, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_BY_LINK_ID_CHAT_REPO_RECORD, link_id)
	if err != nil {
		return nil, err
	}
	answer := []storage.ChatRepoRecord{}
	for row.Next() {
		record := storage.ChatRepoRecord{Chat: &storage.Chat{}, Repo: &storage.Repo{}}
		err = row.Scan(&record.ID, &record.Chat.ChatID, &record.Chat.Type, &record.Repo.ID, &record.Repo.Name, &record.Repo.Owner, &record.Repo.Link,
			&record.Repo.LastCommit, &record.Repo.LastIssue, &record.Repo.LastPR, &record.Tags, &record.Events)
		if err != nil {
			return nil, err
		}
		answer = append(answer, record)
	}
	err = tx.Commit(ctx)
	return answer, err
}

func (ps *PostgresChatRepoRecordStore) GetRecordById(id int) (*storage.ChatRepoRecord, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	record := storage.ChatRepoRecord{Chat: &storage.Chat{}, Repo: &storage.Repo{}}
	err = tx.QueryRow(ctx, GET_BY_ID_CHAT_REPO_RECORD, id).Scan(&record.ID, &record.Chat.ChatID, &record.Chat.Type, &record.Repo.ID, &record.Repo.Name, &record.Repo.Owner, &record.Repo.Link,
		&record.Repo.LastCommit, &record.Repo.LastIssue, &record.Repo.LastPR, &record.Tags, &record.Events)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &record, err
}

func (ps *PostgresChatRepoRecordStore) GetRecordOffset(start int, limit int) ([]storage.ChatRepoRecord, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_OFFSET_CHAT_REPO_RECORD, start, limit)
	if err != nil {
		return nil, err
	}
	answer := []storage.ChatRepoRecord{}
	for row.Next() {
		record := storage.ChatRepoRecord{Chat: &storage.Chat{}, Repo: &storage.Repo{}}
		err = row.Scan(&record.ID, &record.Chat.ChatID, &record.Chat.Type, &record.Repo.ID, &record.Repo.Name, &record.Repo.Owner, &record.Repo.Link,
			&record.Repo.LastCommit, &record.Repo.LastIssue, &record.Repo.LastPR, &record.Tags, &record.Events)
		if err != nil {
			return nil, err
		}
		answer = append(answer, record)
	}
	err = tx.Commit(ctx)
	return answer, err
}

func (ps *PostgresChatRepoRecordStore) GetRecordByChatAndLink(chatId int64, linkId int) (*storage.ChatRepoRecord, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	record := storage.ChatRepoRecord{Chat: &storage.Chat{}, Repo: &storage.Repo{}}
	err = tx.QueryRow(ctx, GET_BY_CHAT_AND_LINK_CHAT_REPO_RECORD, chatId, linkId).Scan(&record.ID, &record.Chat.ChatID, &record.Chat.Type, &record.Repo.ID, &record.Repo.Name, &record.Repo.Owner, &record.Repo.Link,
		&record.Repo.LastCommit, &record.Repo.LastIssue, &record.Repo.LastPR, &record.Tags, &record.Events)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &record, err
}

func (ps *PostgresChatRepoRecordStore) GetRecordNumber() (int, error) {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var number int
	err = tx.QueryRow(ctx, GET_NUMBER_CHAT_REPO_RECORD).Scan(&number)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return number, err
}

func (ps *PostgresChatRepoRecordStore) UpdateRecord(record *storage.ChatRepoRecord) error {
	conn, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ps.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, UPDATE_CHAT_REPO_RECORD, record.Chat.ChatID, record.Repo.ID, record.Tags, record.Events, record.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (ps *PostgresChatRepoRecordStore) Close() {
	if ps.pool != nil {
		ps.pool.Close()
		ps.pool = nil
	}
}
