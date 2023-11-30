package data

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type LocalImageRepository struct {
	Filepath string
}

func (r *LocalImageRepository) AddImage(image AddImageJSON, finished chan bool, errChan chan error) {
	file, err := os.Create(path.Join(r.Filepath, image.ImageName))
	if err != nil {
		finished <- false
		errChan <- err
		return
	}
	defer file.Close()
	_, err = file.Write(image.Image)
	if err != nil {
		finished <- false
	} else {
		finished <- true
	}
	errChan <- err
}

func (r *LocalImageRepository) GetImageReference(name string) (string, error) {
	files, err := os.ReadDir(r.Filepath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filename := strings.Split(file.Name(), ".")[0]
		if strings.Compare(filename, name) == 0 {
			return path.Join(r.Filepath, filename), nil
		}
	}
	return "", fmt.Errorf("File with name %v not found in directory %v", name, r.Filepath)
}

func (r *LocalImageRepository) DeleteImage(name string, finished chan bool) error {
	reference, err := r.GetImageReference(name)
	if err != nil {
		finished <- false
		return err
	}

	err = os.Remove(reference)

	if err != nil {
		finished <- false
	} else {
		finished <- true
	}

	return err
}

func (r *LocalImageRepository) ReplaceImage(updatedImage UpdateImageJSON, finished chan bool) {
	originalReference, err := r.GetImageReference(updatedImage.OriginalName)
	if err != nil {
		finished <- false
		return
	}
	if updatedImage.NewContent == nil && strings.Compare(updatedImage.OriginalName, updatedImage.NewName) != 0 {
		_, ext, _ := strings.Cut(originalReference, ".")
		newReference := path.Join(r.Filepath, updatedImage.NewName+"."+ext)
		err := os.Rename(originalReference, newReference)
		if err != nil {
			finished <- false
			return
		}
	} else if updatedImage.NewContent != nil {
		file, err := os.OpenFile(originalReference, os.O_RDWR, os.ModeAppend)
		if err != nil {
			finished <- false
			return
		}
		defer file.Close()

		err = file.Truncate(0)
		if err != nil {
			finished <- false
			return
		}

		_, err = file.Write(updatedImage.NewContent)
		if err != nil {
			finished <- false
			return
		}

		_, ext, _ := strings.Cut(originalReference, ".")
		newReference := path.Join(r.Filepath, updatedImage.NewName+"."+ext)
		err = os.Rename(originalReference, newReference)
		if err != nil {
			finished <- false
			return
		}
	}
	finished <- true
}
