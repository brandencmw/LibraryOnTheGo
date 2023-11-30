package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

type ImageJSON interface {
	AddImageJSON | UpdateImageJSON
}

type AddImageJSON struct {
	Image     []byte `json:"imageContent"`
	ImageName string `json:"imageName"`
}

type UpdateImageJSON struct {
	OriginalName string `json:"originalName"`
	NewName      string `json:"newName"`
	NewContent   []byte `json:"newContent"`
}

func (r *S3ImageRepository) AddImage(ctx context.Context, img AddImageJSON) error {
	body, err := marshalJSON(img)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, r.basePath+"/add", body)
	if err != nil {
		return err
	}
	req.WithContext(ctx)

	res, err := r.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Error with request: %v", err.Error())
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error from S3 server: %v", resBody)
	}
	return nil
}

func marshalJSON[T ImageJSON](img T) (*bytes.Reader, error) {

	jsonBody, err := json.Marshal(img)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBody), nil
}

func (r *S3ImageRepository) DeleteImage(ctx context.Context, imageName string) error {
	req, err := http.NewRequest(http.MethodDelete, r.basePath+"/delete?img-name="+imageName, nil)
	if err != nil {
		return err
	}
	req.WithContext(ctx)

	res, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Error from S3 server: %+v", string(resBody))
	}
	return nil
}

func (r *S3ImageRepository) ReplaceImage(ctx context.Context, updatedImage UpdateImageJSON) error {
	body, err := marshalJSON(updatedImage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, r.basePath+"/update", body)
	if err != nil {
		return err
	}
	req.WithContext(ctx)

	res, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Error from S3 server: %v", res.Body)
	}
	return nil
}

func (r *S3ImageRepository) GetImageReference(ctx context.Context, imageName string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, r.basePath+"/key?img="+imageName, nil)
	if err != nil {
		return "", err
	}
	req.WithContext(ctx)

	res, err := r.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error with request: %v", err.Error())
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error from S3 server: %+v", string(resBody))
	}

	var jsonData map[string]json.RawMessage

	err = json.Unmarshal([]byte(resBody), &jsonData)
	if err != nil {
		return "", err
	}

	var key string
	err = json.Unmarshal(jsonData["key"], &key)
	if err != nil {
		return "", err
	}

	return key, nil
}
