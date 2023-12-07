package data_test

import (
	"testing"
)

// func TestAddAuthorLocallyWithInvalidSourceFile(t *testing.T) {
// 	expectedAuthor := data.NewAuthor()
// 	proceed := make(chan bool, 1)
// 	errChan := make(chan error, 1)
// 	defer close(proceed)
// 	defer close(errChan)
// 	r := data.LocalAuthorRepository{Filepath: "./test_data", Filename: "invalid.txt"}
// 	go r.CreateAuthor(expectedAuthor, proceed, errChan)
// 	fmt.Printf("Before")
// 	err := <-errChan
// 	if err == nil {
// 		t.Error(err)
// 	}
// }

func TestAddAuthorWithValidSourceFile(t *testing.T) {

}
