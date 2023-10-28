package services

import (
	"bytes"
	"fmt"
	"libraryonthego/server/config"
	"net/http"
)

func SendAuthorImageToS3(imageBytes []byte) error {

	requestBody := []byte(fmt.Sprintf(`{
		"headshot": "%v"
	}`, imageBytes))

	bodyReader := bytes.NewReader(requestBody)

	req, err := http.NewRequest(http.MethodPost, "https://s3_service/add-author-image", bodyReader)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err.Error())
	}

	tlsConfigProvider := config.NewTLS13ConfigProvider(
		"./certificates",
		"client/backend-client.crt",
		"client/backend-client.key",
		[]string{"root-ca.crt"},
	)
	tlsConfig, err := tlsConfigProvider.GetTLSConfig()
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
		fmt.Printf("Error with request: %v\n", err.Error())
	}

	fmt.Printf("Status: %v\n", res.Status)
	return nil
}

func UploadAuthorInfoToDB() {

}
