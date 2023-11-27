package data

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Author struct {
	ID        *uint
	FirstName *string
	LastName  *string
	Bio       *string
}

type AuthorOption func(author *Author)

func WithID(ID uint) AuthorOption {
	return func(author *Author) {
		author.ID = &ID
	}
}

func WithFirstName(fName string) AuthorOption {
	return func(author *Author) {
		author.FirstName = &fName
	}
}

func WithLastName(lName string) AuthorOption {
	return func(author *Author) {
		author.LastName = &lName
	}
}

func WithBio(bio string) AuthorOption {
	return func(author *Author) {
		author.Bio = &bio
	}
}

func NewAuthor(options ...AuthorOption) *Author {
	author := &Author{}
	for _, option := range options {
		option(author)
	}
	return author
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

func (r *PostgresAuthorRepository) CreateAuthor(author Author, proceed chan bool, errChan chan error) {
	tx, err := r.connPool.BeginTx(context.TODO(), pgx.TxOptions{})
	defer tx.Rollback(context.TODO())
	if err != nil {
		errChan <- err
		return
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO authors(first_name, last_name, bio) VALUES ($1, $2, $3)", author.FirstName, author.LastName, author.Bio)
	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil
	if <-proceed {
		tx.Commit(context.TODO())
	}
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

	return Author{ID: &ID, FirstName: &firstName, LastName: &lastName, Bio: &bio}, nil
}

func (r *PostgresAuthorRepository) GetAllAuthors() ([]Author, error) {
	conn, err := r.connPool.Acquire(context.TODO())
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), "SELECT id, first_name, last_name, bio FROM authors ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	var rowID uint
	var rowFirstName, rowLastName, rowBio string
	authors := make([]Author, 0)
	for rows.Next() {
		err := rows.Scan(&rowID, &rowFirstName, &rowLastName, &rowBio)
		if err != nil {
			return authors, err
		}
		authors = append(authors, *NewAuthor(WithID(rowID), WithFirstName(rowFirstName), WithLastName(rowLastName), WithBio(rowBio)))
	}
	return authors, nil
}

func (r *PostgresAuthorRepository) UpdateAuthor(ID uint, author Author, proceed chan bool, result chan AuthorWithErr) {
	a := AuthorWithErr{}
	tx, err := r.connPool.BeginTx(context.TODO(), pgx.TxOptions{})
	defer tx.Rollback(context.TODO())
	if err != nil {
		a.Err = err
		result <- a
		return
	}

	// Get updated first and last name to update image name
	row := tx.QueryRow(context.TODO(),
		`UPDATE authors
		SET
			first_name=COALESCE($1, first_name),
			last_name=COALESCE($2, last_name),
			bio=COALESCE($3, bio)
		WHERE
			id=$4
		RETURNING
			first_name, last_name`,
		author.FirstName, author.LastName, author.Bio, author.ID)

	var rowFName, rowLName string
	err = row.Scan(&rowFName, &rowLName)
	if err != nil {
		a.Err = err
		result <- a
		return
	}

	a.FirstName = &rowFName
	a.LastName = &rowLName

	result <- a

	commit := <-proceed
	if commit {
		tx.Commit(context.TODO())
	}
}

func (r *PostgresAuthorRepository) DeleteAuthor(ID uint, proceed chan bool, result chan AuthorWithErr) {
	a := AuthorWithErr{}
	tx, err := r.connPool.BeginTx(context.TODO(), pgx.TxOptions{})
	defer tx.Rollback(context.TODO())
	if err != nil {
		a.Err = err
		result <- a
		return
	}

	row := tx.QueryRow(context.TODO(), "DELETE FROM authors WHERE id=$1 RETURNING first_name, last_name", ID)
	err = row.Scan(&a.FirstName, &a.LastName)
	if err != nil {
		a.Err = err
		result <- a
		return
	}

	result <- a

	commit := <-proceed
	if commit {
		tx.Commit(context.TODO())
	}
}
