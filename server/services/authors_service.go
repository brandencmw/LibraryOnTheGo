package services

import (
	"context"
	"errors"
	"libraryonthego/server/data"
	"libraryonthego/server/files"
	"sync"
)

// AuthorRepository will allow the service to perform operations in order to interact with stored author information
type AuthorRepository interface {
	CreateAuthor(ctx context.Context, author *data.Author, commit chan bool) error
	GetAuthor(ctx context.Context, ID string) (*data.Author, error)
	GetAllAuthors(ctx context.Context) ([]*data.Author, error)
	UpdateAuthor(ctx context.Context, author *data.Author, commit chan bool) error
	DeleteAuthor(ctx context.Context, ID string, commit chan bool) error
	SearchByName(ctx context.Context, name string, maxResults uint) ([]*data.Author, error)
}

// ImageRepository will allow the service to perform operations on an image associated with a particular author
type ImageRepository interface {
	AddImage(ctx context.Context, image data.AddImageJSON) error
	GetImageReference(ctx context.Context, name string) (string, error)
	DeleteImage(ctx context.Context, name string) error
	ReplaceImage(ctx context.Context, updatedImage data.UpdateImageJSON) error
}

// Image represents data for an image that needs an operation to be performed on it
type Image struct {
	Content []byte
	Name    string
}

// AuthorSearchParams represents information that can be used in search queries for authors
type AuthorSearchParams struct {
	MaxResults    uint
	SearchTerms   string
	IncludeImages bool
}

// author represents all of an author's possible input fields to the
// service layer. It is not exported to force use of the constructor
// function and functional options pattern
// Using functional options pattern for input so variable number
// of fields can be specified for updating author.
// Important to distinguish difference between empty string provided
// and no value provided.
type author struct {
	FirstName *string
	LastName  *string
	Bio       *string
	Headshot  *Image
}

// AuthorOption defines a function that modifies part of the author struct and could fail
type AuthorOption func(author *author) error

// WithLastName adds a last name field to the given author as an option function
func WithFirstName(fName string) AuthorOption {
	return func(author *author) error {
		author.FirstName = &fName
		return nil
	}
}

// WithLastName adds a last name field to the given author as an option function
func WithLastName(lName string) AuthorOption {
	return func(author *author) error {
		author.LastName = &lName
		return nil
	}
}

// WithImage adds a bio field to the given author as an option function
func WithBio(bio string) AuthorOption {
	return func(author *author) error {
		author.Bio = &bio
		return nil
	}
}

// WithImage adds an image struct to the given author as an option function
func WithImage(img *Image) AuthorOption {
	return func(author *author) error {
		if img.Name == "" {
			return errors.New("Image must have a name")
		} else if len(img.Content) == 0 {
			return errors.New("Image must have content")
		}
		author.Headshot = img
		return nil
	}
}

// NewAuthor creates a new author object using the functional options pattern and returns
// a reference to the author which may be partially completed if an error occurs part way
// through handling the options
func NewAuthor(options ...AuthorOption) (*author, error) {
	author := &author{}
	for _, option := range options {
		err := option(author)
		if err != nil { // Return partially built author if error is encountered
			return author, err
		}
	}
	return author, nil
}

// AuthorOutput defines all fields that may be passed out of the service for a retrieved author
// Unlike author struct, fields are not pointers as it is not important to distinguish between nil and empty
// string values
type AuthorOutput struct {
	ID                string
	FirstName         string
	LastName          string
	Bio               string
	HeadshotObjectKey string
}

// AuthorsService defines a concrete implementation for performing business logic
// on author objects
type AuthorsService struct {
	ImageRepo ImageRepository
	DataRepo  AuthorRepository
}

// NewAuthorsService creates a new AuthorsService object using an ImageRepository and Author Repository.
// It returns a reference to the newly built service
func NewAuthorsService(imageRepo ImageRepository, dataRepo AuthorRepository) *AuthorsService {
	return &AuthorsService{
		ImageRepo: imageRepo,
		DataRepo:  dataRepo,
	}
}

// AddAuthor adds an author to the system, storing the image in its associated image
// repository and the other fields in its data repository
func (s *AuthorsService) AddAuthor(parent context.Context, a *author) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	errChan := make(chan error, 2) // Collects errors from goroutines
	commitChan := make(chan bool)  // Signals to data store that it does/does not need to rollback operation
	var wg sync.WaitGroup

	wg.Add(1)
	// This go routine adds new author data to the data store and may return an error to the channel
	go func() {
		defer wg.Done()
		newAuthor := data.NewAuthor(data.WithFirstName(*a.FirstName), data.WithLastName(*a.LastName), data.WithBio(*a.Bio))
		err := s.DataRepo.CreateAuthor(ctx, newAuthor, commitChan)
		if err != nil {
			errChan <- err
			cancel()
		}
	}()

	wg.Add(1)
	// This goroutine adds a new author image to the image store
	go func() {
		defer wg.Done()
		defer close(commitChan) // This is the only routine that writes to the channel so we can close it when it is done
		imageName := files.CreateFriendlyFileName(a.Headshot.Name, *a.FirstName, *a.LastName)
		if err := s.ImageRepo.AddImage(ctx, data.AddImageJSON{ImageName: imageName, Image: a.Headshot.Content}); err != nil {
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

func (s *AuthorsService) getHeadshotObjectKey(ctx context.Context, a *AuthorOutput) error {
	fileName := files.CreateFriendlyFileName("", a.FirstName, a.LastName)
	reference, err := s.ImageRepo.GetImageReference(ctx, fileName)
	if err != nil {
		return err
	}
	a.HeadshotObjectKey = reference
	return nil
}

// collectErrors is a private method for services that will collect errors from a closed error channel and
// return them as a single joined error. This function may be blocking if the channel is not closed
func collectErrors(errChan chan error) error {
	errs := make([]error, 0)
	for err := range errChan {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// GetAllAuthors retrieves information for every author in the system. Since retrieving image data may be an intensive task
// it can be skipped by the includeImages flag being set to false. Returns a list of pointers to output objects as the objects
// may be several hundred bytes making the whole structure several KB or more so making copies would be a memory intensive
// operation
func (s *AuthorsService) GetAllAuthors(parent context.Context, includeImages bool) ([]*AuthorOutput, error) {
	dataAuthors, err := s.DataRepo.GetAllAuthors(parent)
	if err != nil {
		return nil, err
	}

	authors := make([]*AuthorOutput, 0, len(dataAuthors))
	var wg sync.WaitGroup
	errs := make([]error, 0, len(authors))
	for _, author := range dataAuthors {
		ao := &AuthorOutput{ID: *author.ID, FirstName: *author.FirstName, LastName: *author.LastName, Bio: *author.Bio}
		authors = append(authors, ao)
		if includeImages {
			wg.Add(1)
			go func(author *AuthorOutput) {
				defer wg.Done()
				errs = append(errs, s.getHeadshotObjectKey(parent, author))
			}(ao)
		}
	}
	wg.Wait()

	return authors, errors.Join(errs...)
}

// GetAuthor retrieves data for one author associated with the provided ID. Retrieving image data may be a time intensive operation
// so there is the option to skip retrieving the image if necessary
func (s *AuthorsService) GetAuthor(parent context.Context, ID string, includeImage bool) (AuthorOutput, error) {
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
func (s *AuthorsService) DeleteAuthor(parent context.Context, ID string) error {
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

func (s *AuthorsService) UpdateAuthor(parent context.Context, ID string, a author) error {
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

func (s *AuthorsService) SearchAuthorsByName(parent context.Context, params AuthorSearchParams) ([]*AuthorOutput, error) {

	dataAuthors, err := s.DataRepo.SearchByName(parent, params.SearchTerms, params.MaxResults)
	if err != nil {
		return nil, err
	}

	authors := make([]*AuthorOutput, 0, len(dataAuthors))
	var wg sync.WaitGroup
	errs := make([]error, 0, len(dataAuthors))
	for _, author := range dataAuthors {
		ao := &AuthorOutput{ID: *author.ID, FirstName: *author.FirstName, LastName: *author.LastName, Bio: *author.Bio}
		authors = append(authors, ao)
		if params.IncludeImages {
			wg.Add(1)
			go func(author *AuthorOutput) {
				defer wg.Done()
				errs = append(errs, s.getHeadshotObjectKey(parent, ao))
			}(ao)
		}
	}
	wg.Wait()

	return authors, errors.Join(errs...)
}
