package data

import (
	"bytes"
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

type ImageJSON struct {
	Image     []byte `json:"imageContent"`
	ImageName string `json:"imageName"`
}

func (r *S3ImageRepository) AddImage(img ImageJSON) error {
	body, err := marshalImageJSON(img)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, r.basePath+"/add", body)
	if err != nil {
		return err
	}

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

func marshalImageJSON(img ImageJSON) (*bytes.Reader, error) {

	jsonBody, err := json.Marshal(img)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBody), nil
}

func (r *S3ImageRepository) DeleteImage(imageName string) error {
	return nil
}

func (r *S3ImageRepository) ReplaceImage(ImageJSON) error {
	return nil
}

func (r *S3ImageRepository) GetImage(imageName string) (*ImageJSON, error) {

	req, err := http.NewRequest(http.MethodGet, r.basePath+"?img-name="+imageName, nil)
	if err != nil {
		return nil, err
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error with request: %v", err.Error())
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error from S3 server: %+v", string(resBody))
	}

	var img ImageJSON
	fmt.Printf("JSON: %+v\n", string(resBody))
	err = json.Unmarshal(resBody, &img)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal: %v", err.Error())
	}

	return &img, nil
}
