package services

import (
	"fmt"
	"libraryonthego/server/repositories"
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
	ImageRepo repositories.ImageRepository
}

func NewAuthorsService(imageRepo repositories.ImageRepository) *AuthorsService {
	return &AuthorsService{
		ImageRepo: imageRepo,
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

	err = uploadAuthorInfo(author.FirstName, author.LastName, author.Bio)
	return err
}

func uploadAuthorInfo(firstName, lastName, bio string) error {
	return nil
}
