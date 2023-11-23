package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var FILE_EXTENSIONS []string = []string{"jpg", "png"}

type ObjectNotFoundErr struct {
	image  string
	bucket string
}

func (e *ObjectNotFoundErr) Error() string {
	return fmt.Sprintf("Could not find object %v in bucket %v", e.image, e.bucket)
}

type BucketService struct {
	Bucket    string
	Directory string
	Client    *s3.Client
}

type Image struct {
	Content []byte
	Name    string
}

func NewBucketService(bucket, directory string, client *s3.Client) *BucketService {
	return &BucketService{
		Bucket:    bucket,
		Directory: directory,
		Client:    client,
	}
}

func (s *BucketService) UploadImage(imageName string, image []byte) error {
	imageKey := path.Join(s.Directory, imageName)
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &imageKey,
		Body:   bytes.NewReader(image),
	})

	return err
}

func (s *BucketService) GetImage(imageName string) (*Image, error) {

	var objectKey string
	var output *s3.GetObjectOutput
	var err error
	for _, ext := range FILE_EXTENSIONS {
		objectKey = s.Directory + "/" + imageName + "." + ext
		output, err = s.Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &s.Bucket, Key: &objectKey})
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
	imageName = strings.Split(objectKey, "/")[1]
	return &Image{
		Content: imageContent,
		Name:    imageName,
	}, nil
}

func (s *BucketService) findObjectKey(imageName string) (string, error) {

	input := &s3.GetObjectAttributesInput{
		Bucket:           &s.Bucket,
		ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesObjectSize},
	}

	var objectKey string
	var output *s3.GetObjectAttributesOutput
	var err error
	for _, ext := range FILE_EXTENSIONS {
		objectKey = s.Directory + "/" + imageName + "." + ext
		input.Key = &objectKey
		output, err = s.Client.GetObjectAttributes(context.TODO(), input)
		if output != nil && output.ObjectSize != 0 {
			break
		}
	}
	return objectKey, err
}

func (s *BucketService) DeleteImage(imageName string) error {
	bucketName := "library-pictures"

	objectKey, err := s.findObjectKey(imageName)
	if err != nil {
		return &ObjectNotFoundErr{bucket: bucketName, image: imageName}
	}

	_, err = s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{Bucket: &bucketName, Key: &objectKey})
	return err
}

func (s *BucketService) ReplaceImage(objectKey string, newImage []byte) error {
	return nil
}
