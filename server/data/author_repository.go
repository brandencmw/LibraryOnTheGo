package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateAuthorData struct {
	FirstName string
	LastName  string
	Bio       string
}

type AuthorRepository interface {
	CreateAuthor(CreateAuthorData) error
	GetAuthor()
	UpdateAuthor()
	DeleteAuthor()
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

func (r *PostgresAuthorRepository) CreateAuthor(author CreateAuthorData) error {
	conn, _ := r.connPool.Acquire(context.TODO())
	defer conn.Release()
	conn.Exec(context.TODO(), "INSERT INTO authors(first_name, last_name, bio) VALUES ($1, $2, $3)", author.FirstName, author.LastName, author.Bio)
	return nil
}

func (r *PostgresAuthorRepository) GetAuthor() {
}

func (r *PostgresAuthorRepository) UpdateAuthor() {
}

func (r *PostgresAuthorRepository) DeleteAuthor() {
}
