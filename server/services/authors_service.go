package services

import (
	"bytes"
	"fmt"
	"libraryonthego/server/config"
	"net/http"
)

func SendAuthorImageToS3(imageBytes []byte) string {

	requestBody := []byte(fmt.Sprintf(`{
		"imageData": "%v"
	}`, imageBytes))

	bodyReader := bytes.NewReader(requestBody)

	req, err := http.NewRequest(http.MethodPost, "https://s3_service/add-author-image", bodyReader)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config.ClientTLS,
		},
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error with request: %v\n", err.Error())
	}

	fmt.Printf("Status: %v\n", res.Status)
	return ""
}

func UploadAuthorInfoToDB() {

}
