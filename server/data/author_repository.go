package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Author struct {
	FirstName string
	LastName  string
	Bio       string
}

type AuthorWithErr struct {
	Author
	Err error
}

type PostgresAuthorRepository struct {
	connPool *pgxpool.Pool
}

func CreatePostgresConnectionPool() (*pgxpool.Pool, error) {
	connString := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DB") + "?sslmode=verify-full&sslrootcert=./certificates/root-ca.crt&sslcert=./certificates/client/backend-client.crt&sslkey=./certificates/client/backend-client.key"

	pool, err := pgxpool.New(context.TODO(), connString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func NewPostgresAuthorRepository(pool *pgxpool.Pool) *PostgresAuthorRepository {
	return &PostgresAuthorRepository{
		connPool: pool,
	}
}

func (r *PostgresAuthorRepository) CreateAuthor(author Author) error {
	conn, err := r.connPool.Acquire(context.TODO())
	defer conn.Release()
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.TODO(), "INSERT INTO authors(first_name, last_name, bio) VALUES ($1, $2, $3)", author.FirstName, author.LastName, author.Bio)
	return err
}

func (r *PostgresAuthorRepository) GetAuthor(ID uint) (Author, error) {
	conn, err := r.connPool.Acquire(context.TODO())
	if err != nil {
		return Author{}, err
	}

	defer conn.Release()
	row := conn.QueryRow(context.TODO(), "SELECT first_name, last_name, bio FROM authors WHERE id=$1", ID)
	var firstName, lastName, bio string
	err = row.Scan(&firstName, &lastName, &bio)
	if err != nil {
		return Author{}, err
	}

	return Author{FirstName: firstName, LastName: lastName, Bio: bio}, nil
}

func (r *PostgresAuthorRepository) GetAllAuthors() (map[uint]Author, error) {
	conn, err := r.connPool.Acquire(context.TODO())
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), "SELECT id, first_name, last_name, bio FROM authors")
	if err != nil {
		return nil, err
	}

	var rowID uint
	var rowFirstName, rowLastName, rowBio string
	authors := make(map[uint]Author, 0)
	for rows.Next() {
		err := rows.Scan(&rowID, &rowFirstName, &rowLastName, &rowBio)
		if err != nil {
			return authors, err
		}
		authors[rowID] = Author{FirstName: rowFirstName, LastName: rowLastName, Bio: rowBio}
	}
	return authors, nil
}

func (r *PostgresAuthorRepository) UpdateAuthor(ID uint, author Author) error {
	return nil
}

func (r *PostgresAuthorRepository) DeleteAuthor(ID uint, proceed chan bool, result chan AuthorWithErr) {
	a := AuthorWithErr{}
	tx, err := r.connPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		tx.Rollback(context.TODO())
		a.Err = err
		result <- a
	}

	row := tx.QueryRow(context.TODO(), "DELETE FROM authors WHERE id=$1 RETURNING first_name, last_name", ID)
	err = row.Scan(&a.FirstName, &a.LastName)
	if err != nil {
		tx.Rollback(context.TODO())
		a.Err = err
		result <- a
	}

	result <- a

	commit := <-proceed
	if commit {
		tx.Commit(context.TODO())
	} else {
		tx.Rollback(context.TODO())
	}
}
