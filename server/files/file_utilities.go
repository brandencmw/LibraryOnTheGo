package files

import (
	"mime/multipart"
	"strings"
)

func CreateFriendlyFileName(baseFileName string, friendlyNamePieces ...string) string {
	fileName := strings.ReplaceAll(strings.ToLower(strings.Join(friendlyNamePieces, "-")), " ", "-")
	fileExtension := GetFileExtension(baseFileName)
	return fileName + fileExtension
}

func GetMultipartFormContents(fileHeader *multipart.FileHeader) (fileContents []byte, err error) {
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return nil, err
	}

	fileContents = make([]byte, 0)
	for {
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break
		}
		fileContents = append(fileContents, buffer[:n]...)
	}
	return fileContents, err
}

func GetFileExtension(fileName string) string {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) < 2 {
		return ""
	}
	return "." + fileNameParts[1]
}
