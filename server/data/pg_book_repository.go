package data

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	BookID        *string
	Title         *string
	PublishDate   *time.Time
	PageCount     *int
	Synopsis      *string
	AuthorNames   []string
	CategoryNames []string
}

type Category struct {
	CategoryID *string
	Name       *string
}

type PostgresBookRepository struct {
	connPool *pgxpool.Pool
}

func NewPostgresBookRespository(pool *pgxpool.Pool) *PostgresBookRepository {
	return &PostgresBookRepository{
		connPool: pool,
	}
}

func createListOfInsertRows(rows ...string) string {
	rowsBuilder := strings.Builder{}
	for _, row := range rows {
		rowsBuilder.WriteString("(" + row + "),")
	}
	return strings.TrimRight(rowsBuilder.String(), ",")
}

func (r *PostgresBookRepository) CreateBook(ctx context.Context, b *Book, commitChan chan bool) error {
	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Insert to books table
	row := tx.QueryRow(ctx,
		`INSERT INTO books 
			(title, synopsis, publish_date, page_count) 
		VALUES 
			($1, $2, $3, $4)
		RETURNING id`,
		b.Title, b.Synopsis, b.PublishDate, b.PageCount,
	)
	var bookID int
	err = row.Scan(&bookID)
	if err != nil {
		return err
	}

	// Insert categories and get IDs
	categoryRows := createListOfInsertRows(b.CategoryNames...)
	rows, err := tx.Query(ctx,
		`INSERT INTO categories 
			(category_name) 
		VALUES ($1) 
		ON CONFLICT (category_name) DO NOTHING 
		RETURNING id`,
		categoryRows,
	)
	if err != nil {
		return err
	}

	var rowID int
	categoryIDs := make([]int, 0)
	for rows.Next() {
		rows.Scan(&rowID)
		categoryIDs = append(categoryIDs, rowID)
	}

	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(ctx, "INSERT INTO book_categories (book_id, category_id) VALUES ($1, $2)", bookID, categoryID)
		if err != nil {
			return err
		}
	}

	authorIDs, err := getIDsForBookAuthors(ctx, tx, b.AuthorNames...)
	if err != nil {
		return err
	}

	for _, authorID := range authorIDs {
		_, err = tx.Exec(ctx, "INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)", bookID, authorID)
		if err != nil {
			return err
		}
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

func getIDsForBookAuthors(ctx context.Context, tx pgx.Tx, names ...string) ([]int, error) {
	IDs := make([]int, 0, len(names))
	var ID int
	for _, name := range names {
		row := tx.QueryRow(ctx, "SELECT id FROM authors WHERE name_search_vector @@ plainto_tsquery('english', $1) LIMIT 1", name)
		err := row.Scan(&ID)
		if err != nil {
			return IDs, err
		}
		IDs = append(IDs, ID)
	}
	return IDs, nil
}
