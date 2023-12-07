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
	Client *s3.Client
}

type ObjectPath struct {
	Bucket     string
	Directory  string
	ObjectName string
}

type Image struct {
	Content []byte
	Name    string
}

func NewBucketService(client *s3.Client) *BucketService {
	return &BucketService{
		Client: client,
	}
}

func (s *BucketService) UploadImage(ctx context.Context, p ObjectPath, image []byte) error {
	imageKey := path.Join(p.Directory, p.ObjectName)
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &p.Bucket,
		Key:    &imageKey,
		Body:   bytes.NewReader(image),
	})

	return err
}

func (s *BucketService) GetObjectKey(ctx context.Context, p ObjectPath) (string, error) {

	input := &s3.GetObjectAttributesInput{
		Bucket:           &p.Bucket,
		ObjectAttributes: []types.ObjectAttributes{types.ObjectAttributesObjectSize}, // arbitary attribute to determine if object was foun
	} // otherwise nil is returned by get operation

	var objectKey string
	var output *s3.GetObjectAttributesOutput
	var err error
	for _, ext := range FILE_EXTENSIONS {
		objectKey = path.Join(p.Directory, p.ObjectName) + "." + ext
		input.Key = &objectKey
		output, err = s.Client.GetObjectAttributes(ctx, input)
		if output != nil && output.ObjectSize != 0 {
			break
		}
	}
	return objectKey, err
}

func (s *BucketService) DeleteImage(ctx context.Context, p ObjectPath) error {
	objectKey, err := s.GetObjectKey(ctx, p)
	if err != nil {
		return &ObjectNotFoundErr{bucket: p.Bucket, image: p.ObjectName}
	}
	return deleteOperation(ctx, s.Client, p.Bucket, objectKey)
}

func deleteOperation(ctx context.Context, client *s3.Client, bucket, key string) error {
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &bucket, Key: &key})
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

func (s *BucketService) ReplaceImage(ctx context.Context, originalLocation ObjectPath, newImageName string, newImage []byte) error {

	// If no changes are needed to name or content, no work is needed
	if newImage == nil && strings.Compare(originalLocation.ObjectName, newImageName) == 0 {
		return nil
	}

	// Need to get original key of object to copy
	originalKey, err := s.GetObjectKey(ctx, originalLocation)
	if err != nil {
		return err
	}

	// If there is no new image content then just the name needs to be changed which can be done with copy operation
	if newImage == nil {
		newKey := buildNewImageKey(originalKey, newImageName) // Since there is no content, the filetype will be the same as the original
		copySource := originalLocation.Bucket + "/" + originalKey
		_, err = s.Client.CopyObject(ctx, &s3.CopyObjectInput{Bucket: &originalLocation.Bucket, CopySource: &copySource, Key: &newKey})
		if err != nil {
			return err
		}
	} else { // If there is image content then the filetype will already be known
		newLocation := ObjectPath{
			Bucket:     originalLocation.Bucket,
			Directory:  originalLocation.Directory,
			ObjectName: newImageName,
		}
		err := s.UploadImage(ctx, newLocation, newImage)
		if err != nil {
			return err
		}
	}

	// If names were the same then upload would just overwrite the old image
	// So, if names were not the same then a copy would be made so the original must be removed
	if strings.Compare(originalLocation.ObjectName, newImageName) != 0 {
		return deleteOperation(ctx, s.Client, originalLocation.Bucket, originalKey)
	}
	return nil
}
