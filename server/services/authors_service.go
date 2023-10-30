package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"libraryonthego/server/config"
	"net/http"
	"path"
)

type AuthorInfo struct {
	Headshot  []byte
	FirstName string
	LastName  string
	Bio       string
}

type AuthorsService interface {
	AddAuthor(AuthorInfo) error
}

type DefaultAuthorsService struct{}

func (s *DefaultAuthorsService) AddAuthor(author AuthorInfo) error {
	err := uploadAuthorImage(author.Headshot)
	if err != nil {
		return err
	}

	err = uploadAuthorInfo(author.FirstName, author.LastName, author.Bio)
	return err
}

func createS3RequestBody(image []byte) (*bytes.Reader, error) {

	reqBody := uploadToS3Request{
		Image: image,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBody), nil

}

func getTLSConfig() (*tls.Config, error) {
	const rootCertFolder = "certificates"
	const clientCertFolder = "client"
	tlsConfigProvider := config.NewTLS13ConfigProvider(
		path.Join(rootCertFolder, clientCertFolder, "backend-client.crt"),
		path.Join(rootCertFolder, clientCertFolder, "backend-client.key"),
		[]string{path.Join(rootCertFolder, "root-ca.crt")},
	)
	return tlsConfigProvider.GetTLSConfig()
}

type uploadToS3Request struct {
	Image []byte `json:"image"`
}

func uploadAuthorImage(imageBytes []byte) error {

	body, err := createS3RequestBody(imageBytes)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://s3_service/authors/add-author-image", body)
	if err != nil {
		return err
	}

	tlsConfig, err := getTLSConfig()
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error with request: %v", err.Error())
	}

	fmt.Printf("Status: %v\n", res.Status)
	return nil
}

func uploadAuthorInfo(firstName, lastName, bio string) error {
	return nil
}
