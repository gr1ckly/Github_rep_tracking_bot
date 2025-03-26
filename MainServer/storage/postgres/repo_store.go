package postgres

import (
	"Crypto_Bot/MainServer/storage"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ADD_REPO                   = `INSERT INTO REPO(NAME, OWNER, LAST_COMMIT, LAST_ISSUE, LAST_PULL_REQUEST) VALUES ($1, $2, $3, $4, $5) RETURNING ID;`
	REMOVE_REPO                = `DELETE FROM REPO WHERE ID = $1`
	GET_BY_ID_REPO             = `SELECT * FROM REPO WHERE ID = $1`
	GET_BY_OWNER_AND_NAME_REPO = `SELECT * FROM REPO WHERE NAME = $1 AND OWNER = $2;`
	GET_REPOS_OFFSET           = `SELECT * FROM REPO order by ID OFFSET $1 LIMIT $2;`
	GET_REPO_NUMBER            = `SELECT COUNT(*) AS TOTAL_ROWS FROM REPO;`
	UPDATE_REPO                = `UPDATE REPO
SET NAME = $1,
    OWNER = $2,
    LAST_COMMIT = $3,
    LAST_ISSUE = $4,
    LAST_PULL_REQUEST = $5
WHERE id = $6;`
)

const (
	ADD_REPO_NAME                   = "add_repo"
	REMOVE_REPO_NAME                = "remove_repo"
	GET_BY_ID_REPO_NAME             = "get_by_id_repo"
	GET_BY_OWNER_AND_NAME_REPO_NAME = "get_by_owner_and_name_repo"
	GET_REPOS_OFFSET_NAME           = "get_repos_offset"
	GET_REPO_NUMBER_NAME            = "get_repo_number"
	UPDATE_REPO_NAME                = "update_repo"
)

var repoStatementMap = map[string]string{
	ADD_REPO_NAME:                   ADD_REPO,
	REMOVE_REPO_NAME:                REMOVE_REPO,
	GET_BY_ID_REPO_NAME:             GET_BY_ID_REPO,
	GET_BY_OWNER_AND_NAME_REPO_NAME: GET_BY_OWNER_AND_NAME_REPO,
	GET_REPOS_OFFSET_NAME:           GET_REPOS_OFFSET,
	GET_REPO_NUMBER_NAME:            GET_REPO_NUMBER,
	UPDATE_REPO_NAME:                UPDATE_REPO,
}

type PostgresRepoStore struct {
	pool *pgxpool.Pool
}

func NewPostgresRepoStore(ctx context.Context, dbUrl string) (*PostgresRepoStore, error) {
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		for key, _ := range repoStatementMap {
			_, stmtErr := conn.Conn().Prepare(context.Background(), key, repoStatementMap[key])
			if stmtErr != nil {
				return stmtErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &PostgresRepoStore{pool: pool}, nil
}

func (pr *PostgresRepoStore) AddNewRepo(ctx context.Context, repo *storage.Repo) (int, error) {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return -1, nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var id int
	err = tx.QueryRow(ctx, ADD_REPO_NAME, repo.Name, repo.Owner, repo.LastCommit, repo.LastIssue, repo.LastPR).Scan(&id)
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return id, err
}

func (pr *PostgresRepoStore) RemoveRepo(ctx context.Context, repo *storage.Repo) error {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, REMOVE_REPO_NAME, repo.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (pr *PostgresRepoStore) GetRepoByID(ctx context.Context, id int) (*storage.Repo, error) {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return nil, nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	repo := storage.Repo{}
	err = tx.QueryRow(ctx, GET_BY_ID_REPO_NAME, id).Scan(&repo.ID, &repo.Name, &repo.Owner, repo.LastCommit, repo.LastIssue, repo.LastPR)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &repo, err
}

func (pr *PostgresRepoStore) GetRepoByOwnerAndName(ctx context.Context, owner string, name string) (*storage.Repo, error) {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return nil, nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	repo := storage.Repo{}
	err = tx.QueryRow(ctx, GET_BY_ID_REPO_NAME, name, owner).Scan(&repo.ID, &repo.Name, &repo.Owner, repo.LastCommit, repo.LastIssue, repo.LastPR)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	return &repo, err
}

func (pr *PostgresRepoStore) GetReposOffset(ctx context.Context, start int, limit int) ([]storage.Repo, error) {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return nil, nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	row, err := tx.Query(ctx, GET_CHATS_OFFSET_NAME, start, limit)
	if err != nil {
		return nil, err
	}
	answer := []storage.Repo{}
	for row.Next() {
		repo := storage.Repo{}
		err = row.Scan(&repo.ID, &repo.Name, &repo.Owner, repo.LastCommit, repo.LastIssue, repo.LastPR)
		if err != nil {
			return nil, err
		}
		answer = append(answer, repo)
	}
	err = tx.Commit(ctx)
	return answer, err
}

func (pr *PostgresRepoStore) GetRepoNumber(ctx context.Context) (int, error) {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return -1, nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	var number int
	err = tx.QueryRow(ctx, GET_REPO_NUMBER).Scan(&number)
	if err != nil {
		return -1, err
	}
	err = tx.Commit(ctx)
	return number, err
}

func (pr *PostgresRepoStore) UpdateRepo(ctx context.Context, repo *storage.Repo) error {
	conn, err := pr.pool.Acquire(ctx)
	if err != nil {
		return nil
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, UPDATE_REPO_NAME, repo.Name, repo.Owner, repo.LastCommit, repo.LastIssue, repo.LastPR, repo.ID)
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
