package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	ID          *string
	Title       *string
	PublishDate *time.Time
	PageCount   *int
	Synopsis    *string
	Authors     []*Author
	Categories  []string
}

type BookOption func(book *Book)

func WithBookID(ID string) BookOption {
	return func(book *Book) {
		book.ID = &ID
	}
}

func WithTitle(title string) BookOption {
	return func(book *Book) {
		book.Title = &title
	}
}

func WithPublishDate(date time.Time) BookOption {
	return func(book *Book) {
		book.PublishDate = &date
	}
}

func WithPageCount(count int) BookOption {
	return func(book *Book) {
		book.PageCount = &count
	}
}

func WithSynopsis(synopsis string) BookOption {
	return func(book *Book) {
		book.Synopsis = &synopsis
	}
}

func WithAuthors(authors []*Author) BookOption {
	return func(book *Book) {
		book.Authors = authors
	}
}

func WithCategories(categories []string) BookOption {
	return func(book *Book) {
		book.Categories = categories
	}
}

func NewBook(options ...BookOption) *Book {
	book := &Book{}
	for _, option := range options {
		option(book)
	}
	return book
}

type Category struct {
	ID   *string
	Name *string
}

type PostgresBookRepository struct {
	connPool *pgxpool.Pool
}

func NewPostgresBookRespository(pool *pgxpool.Pool) *PostgresBookRepository {
	return &PostgresBookRepository{
		connPool: pool,
	}
}

func extractCategoryNames(categories ...Category) []string {
	names := make([]string, len(categories))
	for _, category := range categories {
		names = append(names, *category.Name)
	}
	return names
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
	categoryRows := createListOfInsertRows(b.Categories...)
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

	fmt.Printf("Authors:%+v\n", b.Authors)
	names := extractAuthorNames(b.Authors...)
	fmt.Printf("1:%v\n", names[0])
	fmt.Printf("Author names: %v\n", names)
	authorIDs, err := getIDsForBookAuthors(ctx, tx, names...)
	if err != nil {
		return err
	}
	fmt.Printf("IDs: %v\n", authorIDs)

	fmt.Println("Inserting to book_authors")
	for _, authorID := range authorIDs {
		_, err = tx.Exec(ctx, "INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)", bookID, authorID)
		if err != nil {
			return err
		}
	}
	fmt.Println("Inserted to book_authors")

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
		fmt.Printf("Looking for author: %v\n", name)
		row := tx.QueryRow(ctx, "SELECT id FROM authors WHERE name_search_vector @@ plainto_tsquery('english', $1) LIMIT 1", name)
		err := row.Scan(&ID)
		if err != nil {
			return IDs, err
		}
		IDs = append(IDs, ID)
	}
	return IDs, nil
}

func (r *PostgresBookRepository) GetAllBooks(ctx context.Context) ([]*Book, error) {
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		`SELECT 
			b.id, b.title, b.synopsis, b.publish_date, b.page_count, 
			array_agg(DISTINCT(a.id, a.first_name, a.last_name)) AS authors, 
			array_agg(c.category_name) AS categories FROM books AS b
		JOIN book_authors AS ba ON b.id = ba.book_id
		JOIN authors AS a ON ba.author_id = a.id
		JOIN book_categories AS bc ON b.id = bc.book_id
		JOIN categories AS c ON bc.category_id = c.id
		GROUP BY b.id`,
	)
	if err != nil {
		return nil, err
	}

	var rowID, rowPageCount int
	var rowTitle, rowSynopsis string
	var rowPublishDate time.Time
	var rowAuthors []*Author
	var rowCategories []string
	var books []*Book
	for rows.Next() {
		err = rows.Scan(&rowID, &rowTitle, &rowSynopsis, &rowPublishDate, &rowPageCount, &rowAuthors, &rowCategories)
		if err != nil {
			return books, err
		}
		books = append(books, NewBook(
			WithBookID(fmt.Sprint(rowID)),
			WithTitle(rowTitle),
			WithSynopsis(rowSynopsis),
			WithPublishDate(rowPublishDate),
			WithPageCount(rowPageCount),
			WithAuthors(rowAuthors),
			WithCategories(rowCategories),
		))
	}
	return books, nil
}

func (r *PostgresBookRepository) GetBook(ctx context.Context, ID string) (*Book, error) {
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	row := conn.QueryRow(ctx,
		`SELECT 
			b.id, b.title, b.synopsis, b.publish_date, b.page_count, 
			array_agg(DISTINCT(a.id, a.first_name, a.last_name)) AS authors, 
			array_agg(c.category_name) AS categories FROM books AS b
		JOIN book_authors AS ba ON b.id = ba.book_id
		JOIN authors AS a ON ba.author_id = a.id
		JOIN book_categories AS bc ON b.id = bc.book_id
		JOIN categories AS c ON bc.category_id = c.id
		WHERE b.id=$1 GROUP BY b.id`, ID,
	)

	var rowID, rowPageCount int
	var rowTitle, rowSynopsis string
	var rowPublishDate time.Time
	var rowAuthors []*Author
	var rowCategories []string
	err = row.Scan(&rowID, &rowTitle, &rowSynopsis, &rowPublishDate, &rowPageCount, &rowAuthors, &rowCategories)
	if err != nil {
		return nil, err
	}
	return NewBook(
		WithBookID(ID),
		WithTitle(rowTitle),
		WithSynopsis(rowSynopsis),
		WithPublishDate(rowPublishDate),
		WithPageCount(rowPageCount),
		WithAuthors(rowAuthors),
		WithCategories(rowCategories),
	), nil
}

func (r *PostgresBookRepository) UpdateBook(ctx context.Context, book Book, commitChan chan bool) error {
	intID, err := strconv.ParseUint(*book.ID, 10, 64)
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
		`UPDATE books
		SET
			title=COALESCE($1, title),
			synopsis=COALESCE($2, synopsis)
		WHERE
			id=$3`,
		book.Title, book.Synopsis, intID)
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
