package services

import (
	"bytes"
	"context"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LibraryBucketService struct {
	Directory string
	Client    *s3.Client
}

type Image struct {
	Content []byte
	Name    string
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

func (s *LibraryBucketService) GetImage(objectKey string) (*Image, error) {
	bucketName := "library-pictures"
	fileExtensions := []string{"jpg", "png"}

	var imageName string
	var output *s3.GetObjectOutput
	var err error
	for _, ext := range fileExtensions {
		imageName = s.Directory + "/" + objectKey + "." + ext
		output, err = s.Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &bucketName, Key: &imageName})
		if output != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	defer output.Body.Close()

	imageContent, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return &Image{
		Content: imageContent,
		Name:    imageName,
	}, nil
}

func (s *LibraryBucketService) DeleteImage(objectKey string) error {
	return nil
}

func (s *LibraryBucketService) ReplaceImage(objectKey string, newImage []byte) error {
	return nil
}
