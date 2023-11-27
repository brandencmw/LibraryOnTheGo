package services

import (
	"errors"
	"fmt"
	"libraryonthego/server/data"
	"libraryonthego/server/files"
	"sync"
)

type AuthorRepository interface {
	CreateAuthor(author data.Author, proceed chan bool, errChan chan error)
	GetAuthor(ID uint) (data.Author, error)
	GetAllAuthors() ([]data.Author, error)
	UpdateAuthor(ID uint, author data.Author, proceed chan bool, result chan data.AuthorWithErr)
	DeleteAuthor(ID uint, proceed chan bool, result chan data.AuthorWithErr)
}

type ImageRepository interface {
	AddImage(image data.AddImageJSON, finished chan bool, errChan chan error)
	GetImageReference(name string) (string, error)
	DeleteImage(name string, finished chan bool) error
	ReplaceImage(updatedImage data.UpdateImageJSON, finished chan bool)
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
	ID                uint
	FirstName         string
	LastName          string
	Bio               string
	HeadshotObjectKey string
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
	proceed := make(chan bool)
	errChan := make(chan error)
	defer close(proceed)
	defer close(errChan)
	go s.DataRepo.CreateAuthor(data.Author{
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Bio:       a.Bio,
	}, proceed, errChan)

	err := <-errChan
	if err != nil {
		return err
	}

	imageName := files.CreateFriendlyFileName(a.Headshot.Name, *a.FirstName, *a.LastName)
	go s.ImageRepo.AddImage(data.AddImageJSON{
		Image:     a.Headshot.Content,
		ImageName: imageName,
	}, proceed, errChan)

	err = <-errChan
	return err
}

func (s *AuthorsService) getHeadshotObjectKeySync(a *AuthorOutput) error {
	fileName := files.CreateFriendlyFileName("", a.FirstName, a.LastName)
	reference, err := s.ImageRepo.GetImageReference(fileName)
	if err != nil {
		return err
	}
	a.HeadshotObjectKey = reference
	return nil
}

func (s *AuthorsService) getHeadshotObjectKeyAsync(a *AuthorOutput, wg *sync.WaitGroup, c chan error) {
	if wg != nil {
		defer wg.Done()
	}
	fileName := files.CreateFriendlyFileName("", a.FirstName, a.LastName)
	reference, err := s.ImageRepo.GetImageReference(fileName)
	if err != nil {
		c <- err
		return
	}
	a.HeadshotObjectKey = reference
	c <- nil
}

func (s *AuthorsService) fetchHeadshotObjectKeys(authors []*AuthorOutput) error {
	var wg sync.WaitGroup
	errChannel := make(chan error, len(authors))
	for _, author := range authors {
		wg.Add(1)
		go s.getHeadshotObjectKeyAsync(author, &wg, errChannel)
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
	for _, author := range authorData {
		ao := &AuthorOutput{
			ID:        *author.ID,
			FirstName: *author.FirstName,
			LastName:  *author.LastName,
			Bio:       *author.Bio,
		}
		authors = append(authors, ao)
	}

	if includeImages && err == nil {
		err = s.fetchHeadshotObjectKeys(authors)
	}

	return authors, err
}

func (s *AuthorsService) GetAuthor(ID uint, includeImage bool) (AuthorOutput, error) {
	author, err := s.DataRepo.GetAuthor(ID)
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
		err = s.getHeadshotObjectKeySync(&ao)
		if err != nil {
			return AuthorOutput{}, err
		}
	}
	return ao, nil
}

func (s *AuthorsService) DeleteAuthor(ID uint) error {
	proceed := make(chan bool)
	authorChan := make(chan data.AuthorWithErr)
	defer close(authorChan)
	defer close(proceed)
	go s.DataRepo.DeleteAuthor(ID, proceed, authorChan)
	author := <-authorChan
	if author.Err != nil {
		return fmt.Errorf("Failed to delete author with ID %v: %v", ID, author.Err.Error())
	}

	imageName := files.CreateFriendlyFileName("", *author.FirstName, *author.LastName)

	err := s.ImageRepo.DeleteImage(imageName, proceed)
	return err
}

func (s *AuthorsService) UpdateAuthor(ID uint, a author) error {
	author := data.Author{
		ID:        &ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Bio:       a.Bio,
	}

	proceed := make(chan bool)
	authorChan := make(chan data.AuthorWithErr)
	defer close(proceed)
	defer close(authorChan)

	// If any of these fields have changed then the image needs to be updated with
	// a new title or new content
	if a.FirstName != nil || a.LastName != nil || a.Headshot != nil {
		originalAuthor, err := s.DataRepo.GetAuthor(ID)
		if err != nil {
			return err
		}
		originalFilename := files.CreateFriendlyFileName("", *originalAuthor.FirstName, *originalAuthor.LastName)
		updatedHeadshot := data.UpdateImageJSON{OriginalName: originalFilename}

		if a.Headshot != nil {
			updatedHeadshot.NewContent = a.Headshot.Content
			updatedHeadshot.NewName = files.CreateFriendlyFileName(a.Headshot.Name, originalFilename)
		}

		go s.DataRepo.UpdateAuthor(ID, author, proceed, authorChan)
		updatedAuthor := <-authorChan
		if updatedAuthor.Err != nil {
			return updatedAuthor.Err
		}
		updatedHeadshot.NewName = files.CreateFriendlyFileName(updatedHeadshot.NewName, *updatedAuthor.FirstName, *updatedAuthor.LastName)
		fmt.Println("CONTENT:", updatedHeadshot.NewContent)
		s.ImageRepo.ReplaceImage(updatedHeadshot, proceed)
	} else {
		go s.DataRepo.UpdateAuthor(ID, author, proceed, authorChan)
		updatedAuthor := <-authorChan
		if updatedAuthor.Err != nil {
			return updatedAuthor.Err
		}
		proceed <- true
	}
	return nil
}
