package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Author struct {
	ID        *string
	FirstName *string
	LastName  *string
	Bio       *string
}

type AuthorOption func(author *Author)

func WithAuthorID(ID string) AuthorOption {
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

func (r *PostgresAuthorRepository) CreateAuthor(ctx context.Context, author *Author, commitChan chan bool) error {
	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{}) // Open transaction using cancellable context
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // Use background context so tx can still be rolled back when parent context is cancelled

	_, err = tx.Exec(ctx, "INSERT INTO authors(first_name, last_name, bio) VALUES ($1, $2, $3)", author.FirstName, author.LastName, author.Bio)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case commit := <-commitChan:
		if commit {
			return tx.Commit(ctx)
		}
	}

	return nil
}

func (r *PostgresAuthorRepository) GetAuthor(ctx context.Context, ID string) (*Author, error) {
	intID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid ID provided")
	}
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT first_name, last_name, bio FROM authors WHERE id=$1", intID)
	var firstName, lastName, bio string
	err = row.Scan(&firstName, &lastName, &bio)
	if err != nil {
		return nil, err
	}

	return &Author{ID: &ID, FirstName: &firstName, LastName: &lastName, Bio: &bio}, nil
}

func (r *PostgresAuthorRepository) GetAllAuthors(ctx context.Context) ([]*Author, error) {
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	rows, err := conn.Query(ctx, "SELECT id, first_name, last_name, bio FROM authors ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	var rowID uint
	var rowFirstName, rowLastName, rowBio string
	authors := make([]*Author, 0)
	for rows.Next() {
		err := rows.Scan(&rowID, &rowFirstName, &rowLastName, &rowBio)
		if err != nil {
			return authors, err
		}
		authors = append(authors, NewAuthor(WithAuthorID(fmt.Sprint(rowID)), WithFirstName(rowFirstName), WithLastName(rowLastName), WithBio(rowBio)))
	}
	return authors, nil
}

func (r *PostgresAuthorRepository) UpdateAuthor(ctx context.Context, author *Author, commitChan chan bool) error {
	intID, err := strconv.ParseUint(*author.ID, 10, 64)
	if err != nil {
		return err
	}
	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(context.Background())
	if err != nil {
		return err
	}

	// Get updated first and last name to update image name
	_, err = tx.Exec(ctx,
		`UPDATE authors
		SET
			first_name=COALESCE($1, first_name),
			last_name=COALESCE($2, last_name),
			bio=COALESCE($3, bio)
		WHERE
			id=$4`,
		author.FirstName, author.LastName, author.Bio, intID)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case commit := <-commitChan:
		if commit {
			return tx.Commit(ctx)
		}
	}
	return nil
}

func (r *PostgresAuthorRepository) DeleteAuthor(ctx context.Context, ID string, commitChan chan bool) error {
	intID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		return err
	}
	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM authors WHERE id=$1", intID)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case commit := <-commitChan:
		if commit {
			return tx.Commit(ctx)
		}
	}
	return nil
}

func (r *PostgresAuthorRepository) SearchByName(ctx context.Context, name string, maxResults uint) ([]*Author, error) {
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	baseQuery := "SELECT id, first_name, last_name, bio FROM authors WHERE name_search_vector @@ plainto_tsquery('english', $1)"
	var rows pgx.Rows
	if maxResults > 0 {
		rows, err = conn.Query(ctx, baseQuery+" LIMIT $2;", name, maxResults)
	} else {
		rows, err = conn.Query(ctx, baseQuery+";", name)
	}
	if err != nil {
		return nil, err
	}

	var rowID uint
	var rowFirstName, rowLastName, rowBio string
	authors := make([]*Author, 0)
	for rows.Next() {
		err := rows.Scan(&rowID, &rowFirstName, &rowLastName, &rowBio)
		if err != nil {
			return authors, err
		}
		authors = append(authors, NewAuthor(WithAuthorID(fmt.Sprint(rowID)), WithFirstName(rowFirstName), WithLastName(rowLastName), WithBio(rowBio)))
	}
	return authors, nil
}

func extractAuthorNames(authors ...*Author) []string {
	names := make([]string, 0, len(authors))
	for _, author := range authors {
		names = append(names, strings.Join([]string{*author.FirstName, *author.LastName}, " "))
	}
	return names
}
