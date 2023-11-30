package data

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

type LocalAuthorRepository struct {
	Filepath string
	Filename string
}

func makeAuthorLine(author Author) []byte {
	builder := strings.Builder{}
	builder.WriteString(*author.ID + ",")
	builder.WriteString(*author.FirstName + ",")
	builder.WriteString(*author.LastName + ",")
	builder.WriteString(*author.Bio + "\n")

	return []byte(builder.String())
}

func (r *LocalAuthorRepository) CreateAuthor(author *Author, proceed chan bool, errChan chan error) {
	id := uuid.New().String()
	author.ID = &id

	record := makeAuthorLine(*author)

	file, err := os.OpenFile(path.Join(r.Filepath, r.Filename), os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		errChan <- err
		return
	}
	errChan <- nil

	defer file.Close()

	if p := <-proceed; p {
		_, err = file.Write(record)
	}
	errChan <- err
}

func (r *LocalAuthorRepository) GetAuthor(ID string) (Author, error) {
	file, err := os.Open(r.Filepath)
	if err != nil {
		return Author{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ",")
		if strings.Compare(ID, lineParts[0]) == 0 {
			return Author{
				ID:        &lineParts[0],
				FirstName: &lineParts[1],
				LastName:  &lineParts[2],
				Bio:       &lineParts[3],
			}, nil
		}
	}

	if scanner.Err() != nil {
		return Author{}, scanner.Err()
	}
	return Author{}, fmt.Errorf("Author with ID %v could not be found", ID)
}

func (r *LocalAuthorRepository) GetAllAuthors() ([]Author, error) {
	authors := make([]Author, 0)

	file, err := os.Open(r.Filepath)
	if err != nil {
		return authors, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ",")
		lineAuthor := Author{
			ID:        &lineParts[0],
			FirstName: &lineParts[1],
			LastName:  &lineParts[2],
			Bio:       &lineParts[3],
		}
		authors = append(authors, lineAuthor)
	}

	return authors, scanner.Err()
}

func (r *LocalAuthorRepository) UpdateAuthor(author Author, proceed chan bool, result chan AuthorWithErr) {
	a := AuthorWithErr{}

	file, err := os.Open(path.Join(r.Filepath, r.Filename))
	if err != nil {
		a.Err = err
		result <- a
		return
	}
	defer file.Close()

	tmp, err := os.CreateTemp(r.Filepath, "tempfile_*.txt")
	if err != nil {
		a.Err = err
		result <- a
		return
	}
	defer tmp.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ",")
		if strings.Compare(*a.ID, lineParts[0]) != 0 {
			if _, err = tmp.WriteString(line + "\n"); err != nil {
				a.Err = err
				result <- a
				return
			}
		} else {
			updatedAuthor := Author{ID: &lineParts[0]}
			if author.FirstName != nil {
				updatedAuthor.FirstName = author.FirstName
			} else {
				updatedAuthor.FirstName = &lineParts[1]
			}
			if author.LastName != nil {
				updatedAuthor.LastName = author.LastName
			} else {
				updatedAuthor.LastName = &lineParts[2]
			}
			if author.Bio != nil {
				updatedAuthor.Bio = author.Bio
			} else {
				updatedAuthor.Bio = &lineParts[3]
			}
			authorLine := makeAuthorLine(updatedAuthor)
			if _, err := tmp.WriteString(string(authorLine)); err != nil {
				a.Err = err
				result <- a
				return
			}
			a.Author = updatedAuthor
			result <- a
		}
	}

	if err = scanner.Err(); err != nil {
		a.Err = err
		result <- a
		return
	}

	if !<-proceed {
		tmp.Close()
		os.Remove(tmp.Name())
		return
	}

	file.Close()
	os.Remove(path.Join(r.Filepath, r.Filename))
	os.Rename(tmp.Name(), path.Join(r.Filepath, r.Filename))
}

func (r *LocalAuthorRepository) DeleteAuthor(ID string, proceed chan bool, result chan AuthorWithErr) {
	a := AuthorWithErr{}

	file, err := os.Open(path.Join(r.Filepath, r.Filename))
	if err != nil {
		a.Err = err
		result <- a
		return
	}
	defer file.Close()

	tmp, err := os.CreateTemp(r.Filepath, "tempfile_*.txt")
	if err != nil {
		a.Err = err
		result <- a
		return
	}
	defer tmp.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ",")
		if strings.Compare(ID, lineParts[0]) != 0 {
			if _, err = tmp.WriteString(line + "\n"); err != nil {
				a.Err = err
				result <- a
				return
			}
		} else {
			a.FirstName = &lineParts[1]
			a.LastName = &lineParts[2]
			result <- a
		}
	}

	if err = scanner.Err(); err != nil {
		a.Err = err
		result <- a
		return
	}

	if !<-proceed {
		tmp.Close()
		os.Remove(tmp.Name())
		return
	}

	file.Close()
	os.Remove(path.Join(r.Filepath, r.Filename))
	os.Rename(tmp.Name(), path.Join(r.Filepath, r.Filename))
}
