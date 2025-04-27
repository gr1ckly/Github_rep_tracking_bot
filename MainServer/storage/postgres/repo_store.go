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
	ADD_REPO                   = `INSERT INTO REPO(NAME, OWNER, LINK, LAST_COMMIT, LAST_ISSUE, LAST_PULL_REQUEST) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID;`
	REMOVE_REPO                = `DELETE FROM REPO WHERE ID = $1`
	GET_BY_ID_REPO             = `SELECT * FROM REPO WHERE ID = $1`
	GET_BY_OWNER_AND_NAME_REPO = `SELECT * FROM REPO WHERE NAME = $1 AND OWNER = $2;`
	GET_REPOS_OFFSET           = `SELECT * FROM REPO order by ID OFFSET $1 LIMIT $2;`
	GET_REPO_NUMBER            = `SELECT COUNT(*) AS TOTAL_ROWS FROM REPO;`
	UPDATE_REPO                = `UPDATE REPO
SET NAME = $1,
    OWNER = $2,
    LINK = $3,
    LAST_COMMIT = $4,
    LAST_ISSUE = $5,
    LAST_PULL_REQUEST = $6
WHERE id = $7;`
)

type PostgresRepoStore struct {
	pool    *pgxpool.Pool
	timeout int
}

func NewPostgresRepoStore(timeout int, dbUrl string) (*PostgresRepoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return &PostgresRepoStore{pool: pool, timeout: timeout}, nil
}

func (pr *PostgresRepoStore) AddNewRepo(repo *storage.Repo) (int, error) {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var id int
	err = tx.QueryRow(ctx, ADD_REPO, repo.Name, repo.Owner, repo.Link, repo.LastCommit.Format(time.RFC3339), repo.LastIssue.Format(time.RFC3339), repo.LastPR.Format(time.RFC3339)).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return id, err
}

func (pr *PostgresRepoStore) RemoveRepo(id int) error {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, REMOVE_REPO, id)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (pr *PostgresRepoStore) GetRepoByID(id int) (*storage.Repo, error) {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	repo := storage.Repo{}
	err = tx.QueryRow(ctx, GET_BY_ID_REPO, id).Scan(&repo.ID, &repo.Name, &repo.Owner, &repo.Link, &repo.LastCommit, &repo.LastIssue, &repo.LastPR)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &repo, err
}

func (pr *PostgresRepoStore) GetRepoByOwnerAndName(owner string, name string) (*storage.Repo, error) {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	repo := storage.Repo{}
	err = tx.QueryRow(ctx, GET_BY_OWNER_AND_NAME_REPO, name, owner).Scan(&repo.ID, &repo.Name, &repo.Owner, &repo.Link, &repo.LastCommit, &repo.LastIssue, &repo.LastPR)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &repo, err
}

func (pr *PostgresRepoStore) GetReposOffset(start int, limit int) ([]storage.Repo, error) {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_REPOS_OFFSET, start, limit)
	if err != nil {
		return nil, err
	}
	answer := []storage.Repo{}
	for row.Next() {
		repo := storage.Repo{}
		err = row.Scan(&repo.ID, &repo.Name, &repo.Owner, &repo.Link, &repo.LastCommit, &repo.LastIssue, &repo.LastPR)
		if err != nil {
			return nil, err
		}
		answer = append(answer, repo)
	}
	err = tx.Commit(ctx)
	return answer, err
}

func (pr *PostgresRepoStore) GetRepoNumber() (int, error) {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var number int
	err = tx.QueryRow(ctx, GET_REPO_NUMBER).Scan(&number)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, custom_errors.NewNoValuesError(err)
	}
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return number, err
}

func (pr *PostgresRepoStore) UpdateRepo(repo *storage.Repo) error {
	conn, err := pr.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pr.timeout)*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, UPDATE_REPO, repo.Name, repo.Owner, repo.Link, repo.LastCommit.Format(time.RFC3339), repo.LastIssue.Format(time.RFC3339), repo.LastPR.Format(time.RFC3339), repo.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (pr *PostgresRepoStore) Close() {
	if pr.pool != nil {
		pr.pool.Close()
		pr.pool = nil
	}
}
