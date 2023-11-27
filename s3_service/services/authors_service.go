package services

import (
	"bytes"
	"context"
	"fmt"
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

func (s *BucketService) GetObjectKey(imageName string) (string, error) {

	input := &s3.GetObjectAttributesInput{
		Bucket:           &s.Bucket,
		ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesObjectSize}, // arbitary attribute to determine if object was foun
	} // otherwise nil is returned by get operation

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
	objectKey, err := s.GetObjectKey(imageName)
	if err != nil {
		return &ObjectNotFoundErr{bucket: s.Bucket, image: imageName}
	}
	return deleteOperation(s.Client, s.Bucket, objectKey)
}

func deleteOperation(client *s3.Client, bucket, key string) error {
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{Bucket: &bucket, Key: &key})
	return err
}

func buildNewImageKey(originalKey, newName string) string {
	directory, originalCore, _ := strings.Cut(originalKey, "/")
	originalCore, extension, _ := strings.Cut(originalCore, ".")
	if strings.Compare(originalCore, newName) == 0 {
		return originalKey
	}
	return directory + "/" + newName + "." + extension
}

func (s *BucketService) ReplaceImage(originalImageName, newImageName string, newImage []byte) error {

	// If no changes are needed to name or content, no work is needed
	if newImage == nil && strings.Compare(originalImageName, newImageName) == 0 {
		return nil
	}

	// Need to get original key of object to copy
	originalKey, err := s.GetObjectKey(originalImageName)
	if err != nil {
		return err
	}

	// If there is no new image content then just the name needs to be changed which can be done with copy operation
	if newImage == nil {
		newKey := buildNewImageKey(originalKey, newImageName) // Since there is no content, the filetype will be the same as the original
		copySource := s.Bucket + "/" + originalKey
		_, err = s.Client.CopyObject(context.TODO(), &s3.CopyObjectInput{Bucket: &s.Bucket, CopySource: &copySource, Key: &newKey})
		if err != nil {
			return err
		}
	} else { // If there is image content then the filetype will already be known
		err := s.UploadImage(newImageName, newImage)
		if err != nil {
			return err
		}
	}

	// If names were the same then upload would just overwrite the old image
	// So, if names were not the same then a copy would be made so the original must be removed
	if strings.Compare(originalImageName, newImageName) != 0 {
		return deleteOperation(s.Client, s.Bucket, originalKey)
	}
	return nil
}
