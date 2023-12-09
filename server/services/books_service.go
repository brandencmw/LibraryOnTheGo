package services

import (
	"context"
	"errors"
	"fmt"
	"libraryonthego/server/data"
	"libraryonthego/server/files"
	"strings"
	"sync"
	"time"
)

// AuthorRepository will allow the service to perform operations in order to interact with stored author information
type BookRepository interface {
	CreateBook(ctx context.Context, book *data.Book, commit chan bool) error
	GetAllBooks(ctx context.Context) ([]*data.Book, error)
}

type Book struct {
	Title       *string
	Synopsis    *string
	PublishDate *time.Time
	PageCount   *int
	Authors     []Author
	Categories  []string
	Cover       *Image
}

// BookOption defines a function that modifies part of the book struct and could fail
type BookOption func(book *Book) error

func WithTitle(title string) BookOption {
	return func(book *Book) error {
		book.Title = &title
		return nil
	}
}

func WithSynopsis(desc string) BookOption {
	return func(book *Book) error {
		book.Synopsis = &desc
		return nil
	}
}

func WithPublishDate(date string) BookOption {
	return func(book *Book) error {
		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			return err
		}
		book.PublishDate = &date
		return nil
	}
}

func WithPageCount(count int) BookOption {
	return func(book *Book) error {
		if !(count > 0) {
			return errors.New("Page count must be greater than 0")
		}
		book.PageCount = &count
		return nil
	}
}

func WithAuthors(authorNames ...string) BookOption {
	return func(book *Book) error {
		if len(authorNames) == 0 {
			return errors.New("Must have at least one author in list")
		}
		authors := make([]Author, 0, len(authorNames))
		for _, name := range authorNames {
			nameParts := strings.SplitN(name, " ", 2)
			nameOptions := make([]AuthorOption, 0, len(nameParts))
			if len(nameParts) > 0 {
				nameOptions = append(nameOptions, WithFirstName(nameParts[0]))
			}
			if len(nameParts) > 1 {
				nameOptions = append(nameOptions, WithLastName(nameParts[1]))
			}
			author, _ := NewAuthor(nameOptions...)
			fmt.Printf("First name:%v\n", *author.FirstName)
			fmt.Printf("Last name:%v\n", *author.LastName)
			authors = append(authors, *author)
		}
		book.Authors = authors
		return nil
	}
}

func WithCategories(categories ...string) BookOption {
	return func(book *Book) error {
		if len(categories) == 0 {
			return errors.New("Must have at least one category in list")
		}
		book.Categories = categories
		return nil
	}
}

func WithCover(img Image) BookOption {
	return func(book *Book) error {
		if img.Content == nil {
			return errors.New("Image must have content")
		} else if img.Name == "" {
			return errors.New("Image must have a title")
		}
		book.Cover = &img
		return nil
	}
}

// NewAuthor creates a new author object using the functional options pattern and returns
// a reference to the author which may be partially completed if an error occurs part way
// through handling the options
func NewBook(options ...BookOption) (*Book, error) {
	book := &Book{}
	for _, option := range options {
		err := option(book)
		if err != nil { // Return partially built book if error is encountered
			return book, err
		}
	}
	return book, nil
}

type BookOutput struct {
	ID             string
	Title          string
	Synopsis       string
	PublishDate    time.Time
	PageCount      int
	Authors        []AuthorOutput
	Categories     []string
	CoverObjectKey string
}

// AuthorsService defines a concrete implementation for performing business logic
// on author objects
type BooksService struct {
	ImageRepo ImageRepository
	DataRepo  BookRepository
}

func NewBooksService(imageRepo ImageRepository, dataRepo BookRepository) *BooksService {
	return &BooksService{
		ImageRepo: imageRepo,
		DataRepo:  dataRepo,
	}
}

func (s *BooksService) AddBook(parent context.Context, b *Book) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	errChan := make(chan error, 2) // Collects errors from goroutines
	commitChan := make(chan bool)  // Signals to data store that it does/does not need to rollback operation
	var wg sync.WaitGroup

	fmt.Printf("Book authors fName:%v Lname:%v\n", *b.Authors[0].FirstName, *b.Authors[0].LastName)
	wg.Add(1)
	// This go routine adds new book data to the data store and may return an error to the channel
	go func() {
		defer wg.Done()
		newBook := serviceBookToDataBook(b)
		err := s.DataRepo.CreateBook(ctx, newBook, commitChan)
		if err != nil {
			errChan <- err
			cancel()
		}
	}()

	wg.Add(1)
	// This goroutine adds a new book image to the image store
	go func() {
		defer wg.Done()
		defer close(commitChan) // This is the only routine that writes to the channel so we can close it when it is done
		imageName := files.CreateFriendlyFileName(b.Cover.Name, *b.Title)
		if err := s.ImageRepo.AddImage(ctx, data.AddImageJSON{ImageName: imageName, Image: b.Cover.Content}); err != nil {
			errChan <- err
			commitChan <- false
		} else {
			commitChan <- true
		}
	}()

	wg.Wait()      // Wait for all data operations to complete or fail before continuing
	close(errChan) // No more errors will be written once goroutines finish

	return collectErrors(errChan)
}

func (s *AuthorsService) getCoverObjectKey(ctx context.Context, a *AuthorOutput) error {
	fileName := files.CreateFriendlyFileName("", a.FirstName, a.LastName)
	reference, err := s.ImageRepo.GetImageReference(ctx, fileName)
	if err != nil {
		return err
	}
	a.HeadshotObjectKey = reference
	return nil
}

func dataBookToServiceBook(dataBook *data.Book) *BookOutput {
	bo := &BookOutput{}
	if dataBook.ID != nil {
		bo.ID = *dataBook.ID
	}
	if dataBook.Title != nil {
		bo.Title = *dataBook.Title
	}
	if dataBook.Synopsis != nil {
		bo.Synopsis = *dataBook.Synopsis
	}
	if dataBook.PublishDate != nil {
		bo.PublishDate = *dataBook.PublishDate
	}
	if dataBook.PageCount != nil {
		bo.PageCount = *dataBook.PageCount
	}
	authorOutputs := make([]AuthorOutput, 0, len(dataBook.Authors))
	for _, author := range dataBook.Authors {
		authorOutputs = append(authorOutputs, *dataAuthorToServiceAuthor(author))
	}
	bo.Authors = authorOutputs
	bo.Categories = dataBook.Categories
	return bo
}

func serviceBookToDataBook(serviceBook *Book) *data.Book {
	authors := make([]*data.Author, 0, len(serviceBook.Authors))
	for _, author := range serviceBook.Authors {
		authors = append(authors, serviceAuthorToDataAuthor(&author))
	}
	fmt.Printf("FName:%v, LName%v\n", *authors[0].FirstName, *authors[0].LastName)
	return data.NewBook(
		data.WithTitle(*serviceBook.Title),
		data.WithSynopsis(*serviceBook.Synopsis),
		data.WithSynopsis(*serviceBook.Synopsis),
		data.WithPublishDate(*serviceBook.PublishDate),
		data.WithPageCount(*serviceBook.PageCount),
		data.WithAuthors(authors),
		data.WithCategories(serviceBook.Categories),
	)
}

func (s *BooksService) getHeadshotObjectKey(ctx context.Context, b *BookOutput) error {
	fileName := files.CreateFriendlyFileName("", b.Title)
	reference, err := s.ImageRepo.GetImageReference(ctx, fileName)
	if err != nil {
		return err
	}
	b.CoverObjectKey = reference
	return nil
}

func (s *BooksService) GetAllBooks(parent context.Context, includeImages bool) ([]*BookOutput, error) {
	dataBooks, err := s.DataRepo.GetAllBooks(parent)
	if err != nil {
		return nil, err
	}

	books := make([]*BookOutput, 0, len(dataBooks))
	var wg sync.WaitGroup
	errs := make([]error, 0, len(dataBooks))
	for _, book := range dataBooks {
		bo := dataBookToServiceBook(book)
		books = append(books, bo)
		if includeImages {
			wg.Add(1)
			go func(book *BookOutput) {
				defer wg.Done()
				errs = append(errs, s.getHeadshotObjectKey(parent, bo))
			}(bo)
		}
	}

	return books, err
}

// GetAuthor retrieves data for one author associated with the provided ID. Retrieving image data may be a time intensive operation
// so there is the option to skip retrieving the image if necessary
func (s *AuthorsService) GetBook(parent context.Context, ID string, includeImage bool) (AuthorOutput, error) {
	author, err := s.DataRepo.GetAuthor(parent, ID)
	if err != nil {
		return AuthorOutput{}, err
	}

	ao := AuthorOutput{
		ID:        ID,
		FirstName: *author.FirstName,
		LastName:  *author.LastName,
		Bio:       *author.Bio,
	}

	if includeImage {
		err = s.getHeadshotObjectKey(parent, &ao)
		if err != nil {
			return AuthorOutput{}, err
		}
	}
	return ao, nil
}

// DeleteAuthor deletes data for one author associated with the given ID including the base data and image. The deletion of base
// data and image are deleted in parallel to speed up the operation
func (s *AuthorsService) DeleteBook(parent context.Context, ID string) error {
	author, err := s.DataRepo.GetAuthor(parent, ID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	errChan := make(chan error, 2)
	commitChan := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.DataRepo.DeleteAuthor(ctx, ID, commitChan)
		if err != nil {
			errChan <- err
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(commitChan)
		imageName := files.CreateFriendlyFileName("", *author.FirstName, *author.LastName)
		if err := s.ImageRepo.DeleteImage(ctx, imageName); err != nil {
			errChan <- err
			commitChan <- false
		} else {
			commitChan <- true
		}
	}()

	wg.Wait()
	close(errChan)

	err = collectErrors(errChan)
	return err
}

func (s *AuthorsService) UpdateBook(parent context.Context, ID string, a Author) error {
	originalAuthor, err := s.DataRepo.GetAuthor(parent, ID)
	if err != nil {
		return err
	}

	updateAuthorData := &data.Author{
		ID:        &ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Bio:       a.Bio,
	}

	commitChan := make(chan bool, 1)
	defer close(commitChan)
	errChan := make(chan error, 2)

	// If any of these fields have changed then the image needs to be updated with
	// a new title or new content
	if a.FirstName != nil || a.LastName != nil || a.Headshot != nil {
		ctx, cancel := context.WithCancel(parent)
		defer cancel()
		originalFilename := files.CreateFriendlyFileName("", *originalAuthor.FirstName, *originalAuthor.LastName)
		updatedHeadshot := data.UpdateImageJSON{OriginalName: originalFilename}

		if a.Headshot != nil {
			updatedHeadshot.NewContent = a.Headshot.Content
			updatedHeadshot.NewName = files.CreateFriendlyFileName(a.Headshot.Name, originalFilename)
		}

		var wg sync.WaitGroup

		if a.FirstName != nil || a.LastName != nil || a.Bio != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := s.DataRepo.UpdateAuthor(ctx, updateAuthorData, commitChan)
				if err != nil {
					errChan <- err
					cancel()
				}
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if updateAuthorData.FirstName == nil {
				updateAuthorData.FirstName = originalAuthor.FirstName
			}
			if updateAuthorData.LastName == nil {
				updateAuthorData.LastName = originalAuthor.LastName
			}
			updatedHeadshot.NewName = files.CreateFriendlyFileName(updatedHeadshot.NewName, *updateAuthorData.FirstName, *updateAuthorData.LastName)
			err := s.ImageRepo.ReplaceImage(ctx, updatedHeadshot)
			if err != nil {
				errChan <- err
				commitChan <- false
			} else {
				commitChan <- true
			}
		}()
		wg.Wait()
		close(errChan)
		return collectErrors(errChan)
	} else {
		commitChan <- true
		return s.DataRepo.UpdateAuthor(parent, updateAuthorData, commitChan)
	}
}
