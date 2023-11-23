package services

import (
	"errors"
	"fmt"
	"libraryonthego/server/data"
	"libraryonthego/server/files"
	"sync"
)

type AuthorRepository interface {
	CreateAuthor(author data.Author) error
	GetAuthor(ID uint) (data.Author, error)
	GetAllAuthors() (map[uint]data.Author, error)
	UpdateAuthor(ID uint, author data.Author) error
	DeleteAuthor(ID uint, proceed chan bool, result chan data.AuthorWithErr)
}

type ImageRepository interface {
	AddImage(data.ImageJSON) error
	GetImage(string) (*data.ImageJSON, error)
	DeleteImage(string, chan bool) error
	ReplaceImage(data.ImageJSON) error
}

type Image struct {
	Content []byte
	Name    string
}

// Using functional options pattern for input so variable number
// of fields can be specified for updating author
type author struct {
	FirstName *string
	LastName  *string
	Bio       *string
	Headshot  *Image
}

type AuthorOption func(author *author) error

func WithFirstName(fName string) AuthorOption {
	return func(author *author) error {
		author.FirstName = &fName
		return nil
	}
}

func WithLastName(lName string) AuthorOption {
	return func(author *author) error {
		author.LastName = &lName
		return nil
	}
}

func WithBio(bio string) AuthorOption {
	return func(author *author) error {
		author.Bio = &bio
		return nil
	}
}

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

func NewAuthor(options ...AuthorOption) (*author, error) {
	author := &author{}
	for _, option := range options {
		err := option(author)
		if err != nil {
			return author, err
		}
	}
	return author, nil
}

type AuthorOutput struct {
	ID uint
	author
}

type AuthorsService struct {
	ImageRepo ImageRepository
	DataRepo  AuthorRepository
}

func NewAuthorsService(imageRepo ImageRepository, dataRepo AuthorRepository) *AuthorsService {
	return &AuthorsService{
		ImageRepo: imageRepo,
		DataRepo:  dataRepo,
	}
}

func (s *AuthorsService) AddAuthor(a author) error {

	imageName := files.CreateFriendlyFileName(a.Headshot.Name, *a.FirstName, *a.LastName)
	image := data.ImageJSON{
		Image:     a.Headshot.Content,
		ImageName: imageName,
	}
	if err := s.ImageRepo.AddImage(image); err != nil {
		return err
	}

	err := s.DataRepo.CreateAuthor(data.Author{
		FirstName: *a.FirstName,
		LastName:  *a.LastName,
		Bio:       *a.Bio,
	})

	return err
}

func (s *AuthorsService) getAuthorImageAsync(a *AuthorOutput, wg *sync.WaitGroup, c chan error) {
	if wg != nil {
		defer wg.Done()
	}
	fileName := files.CreateFriendlyFileName("", *a.FirstName, *a.LastName)
	content, err := s.ImageRepo.GetImage(fileName)
	if err != nil {
		c <- err
		return
	}
	a.Headshot = &Image{
		Content: content.Image,
		Name:    content.ImageName,
	}
	c <- nil
}

func (s *AuthorsService) fetchAuthorImages(authors *[]*AuthorOutput) error {
	var wg sync.WaitGroup
	errChannel := make(chan error, len(*authors))
	for _, author := range *authors {
		wg.Add(1)
		go s.getAuthorImageAsync(author, &wg, errChannel)
	}

	wg.Wait()
	close(errChannel)

	var errs []error
	for err := range errChannel {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (s *AuthorsService) GetAllAuthors(includeImages bool) ([]*AuthorOutput, error) {

	authorData, err := s.DataRepo.GetAllAuthors()
	authors := make([]*AuthorOutput, 0)
	for id, author := range authorData {
		a, _ := NewAuthor(WithFirstName(author.FirstName), WithLastName(author.LastName), WithBio(author.Bio))
		ao := &AuthorOutput{
			ID:     id,
			author: *a,
		}
		NewAuthor(
			WithFirstName(author.FirstName),
			WithLastName(author.LastName),
			WithBio(author.Bio),
		)
		authors = append(authors, ao)
	}

	if includeImages && err == nil {
		err = s.fetchAuthorImages(&authors)
	}

	return authors, err
}

func (s *AuthorsService) GetAuthor(ID uint, includeImage bool) (AuthorOutput, error) {
	author, err := s.DataRepo.GetAuthor(ID)
	if err != nil {
		return AuthorOutput{}, err
	}

	a, err := NewAuthor(WithFirstName(author.FirstName), WithLastName(author.LastName), WithBio(author.Bio))
	ao := AuthorOutput{
		ID:     ID,
		author: *a,
	}
	if includeImage {
		errChan := make(chan error)
		go s.getAuthorImageAsync(&ao, nil, errChan)
		err = <-errChan
		if err != nil {
			return AuthorOutput{}, err
		}
		close(errChan)
	}
	return ao, nil
}

func (s *AuthorsService) DeleteAuthor(ID uint) error {
	proceed := make(chan bool)
	authorChan := make(chan data.AuthorWithErr)
	go s.DataRepo.DeleteAuthor(ID, proceed, authorChan)
	author := <-authorChan
	if author.Err != nil {
		return fmt.Errorf("Failed to delete author with ID %v", author.Err.Error())
	}

	imageName := files.CreateFriendlyFileName("", author.FirstName, author.LastName)

	err := s.ImageRepo.DeleteImage(imageName, proceed)
	close(authorChan)
	close(proceed)
	return err
}

func UpdateAuthor(ID uint, a author) error {
	return nil
}
