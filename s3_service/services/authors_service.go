package services

import (
	"bytes"
	"context"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LibraryBucketService struct {
	Directory string
	Client    *s3.Client
}

func NewLibraryBucketService(directory string, client *s3.Client) *LibraryBucketService {
	return &LibraryBucketService{
		Directory: directory,
		Client:    client,
	}
}

func (s *LibraryBucketService) UploadImage(imageName string, image []byte) error {

	bucketName := "library-pictures"
	imageKey := path.Join(s.Directory, imageName)
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &imageKey,
		Body:   bytes.NewReader(image),
	})

	return err
}

func (s *LibraryBucketService) GetImage(objectKey string) []byte {
	return nil
}

func (s *LibraryBucketService) DeleteImage(objectKey string) error {
	return nil
}

func (s *LibraryBucketService) ReplaceImage(objectKey string, newImage []byte) error {
	return nil
}
