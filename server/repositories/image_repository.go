package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

type ImageRepository interface {
	AddImage(string, []byte) error
	DeleteImage(string) error
	ReplaceImage(string, []byte) error
}

type S3ImageRepository struct {
	Client   *http.Client
	basePath string
}

func NewS3ImageRepository(client *http.Client, basePath string) *S3ImageRepository {
	return &S3ImageRepository{
		Client:   client,
		basePath: basePath,
	}
}

type uploadToS3Request struct {
	Image     []byte `json:"image"`
	ImageName string `json:"imageName"`
}

func (r *S3ImageRepository) AddImage(imageName string, image []byte) error {
	body, err := createS3RequestBody(imageName, image)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, path.Join(r.basePath, "authors/add-author-image"), body)
	if err != nil {
		return err
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Error with request: %v", err.Error())
	} else if res.StatusCode != 200 {
		return fmt.Errorf("Error from S3 server: %v", res.Body)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Error from S3 server: %v", resBody)
	}

	fmt.Printf("Res body: %v\n", string(resBody))
	return nil
}

func createS3RequestBody(imageName string, image []byte) (*bytes.Reader, error) {

	reqBody := uploadToS3Request{
		Image:     image,
		ImageName: imageName,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBody), nil
}

func (r *S3ImageRepository) DeleteImage(imageName string) error {
	return nil
}

func (r *S3ImageRepository) ReplaceImage(imageName string, newImage []byte) error {
	return nil
}
