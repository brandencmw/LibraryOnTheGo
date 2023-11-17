package services

import (
	"fmt"
	"libraryonthego/server/data"
	"libraryonthego/server/utils"
	"mime/multipart"
)

type AddAuthorInfo struct {
	HeadshotFile *multipart.FileHeader
	FirstName    string
	LastName     string
	Bio          string
}

type AuthorsService struct {
	ImageRepo data.ImageRepository
	DataRepo  data.AuthorRepository
}

func NewAuthorsService(imageRepo data.ImageRepository, dataRepo data.AuthorRepository) *AuthorsService {
	return &AuthorsService{
		ImageRepo: imageRepo,
		DataRepo:  dataRepo,
	}
}

func (s *AuthorsService) AddAuthor(author AddAuthorInfo) error {
	imageContents, err := utils.GetMultipartFormContents(author.HeadshotFile)
	if err != nil {
		return fmt.Errorf("Failed to create image repository: %v\n", err.Error())
	}

	headshotFileName := author.HeadshotFile.Filename
	imageNameToStore := utils.CreateFriendlyFileName(headshotFileName, author.FirstName, author.LastName)

	if err := s.ImageRepo.AddImage(imageNameToStore, imageContents); err != nil {
		return err
	}

	createAuthorData := data.CreateAuthorData{
		FirstName: author.FirstName,
		LastName:  author.LastName,
		Bio:       author.Bio,
	}
	s.DataRepo.CreateAuthor(createAuthorData)
	return nil
}
