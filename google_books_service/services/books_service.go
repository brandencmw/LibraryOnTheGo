package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type BooksService struct {
	client *http.Client
}

type BookSearchOutput struct {
	ID          string
	Title       string
	Authors     []string
	PublishDate string
	Description string
	PageCount   uint
	Categories  []string
}

func NewBooksService(client *http.Client) *BooksService {
	return &BooksService{
		client: client,
	}
}

func buildQuery(parts ...string) string {
	queryWords := strings.Join(parts, " ")
	return url.QueryEscape(queryWords)
}

func extractSearchResults(responseData map[string]any) *BookSearchOutput {
	var output *BookSearchOutput
	if id, ok := responseData["id"].(string); ok {
		output.ID = id
	}

	if volumeInfo, ok := responseData["volumeInfo"].(map[string]any); ok {
		if title, ok := volumeInfo["title"].(string); ok {
			output.Title = title
		}
		if authors, ok := volumeInfo["authors"].([]string); ok {
			output.Authors = authors
		}
		if publishedDate, ok := volumeInfo["publishedData"].(string); ok {
			output.PublishDate = publishedDate
		}
		if description, ok := volumeInfo["description"].(string); ok {
			output.Description = description
		}
		if pageCount, ok := volumeInfo["pageCount"].(int); ok {
			output.PageCount = uint(pageCount)
		}
		if mainCategory, ok := volumeInfo["mainCategory"].(string); ok {
			output.Categories = append(output.Categories, mainCategory)
		}
		if categories, ok := volumeInfo["categories"].([]string); ok {
			output.Categories = append(output.Categories, categories...)
		}
	}
	return output
}

func (s *BooksService) SearchBooks(ctx context.Context, title, author string) ([]*BookSearchOutput, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	query := buildQuery(title, author)
	apiURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?key=%v&q=%v", apiKey, query)

	fmt.Println("URL:", apiURL)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Problem with request: %v", err.Error())
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Problem with response: %v", err.Error())
	}

	defer res.Body.Close()

	var responseData map[string]any
	content, err := io.ReadAll(res.Body)
	fmt.Printf(string(content))
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		return nil, fmt.Errorf("Problem decoding: %v", err.Error())
	}
	fmt.Println("Response:", responseData)
	res.Body.Close()

	var books []*BookSearchOutput
	if responseBooks, ok := responseData["items"].([]map[string]any); ok {
		books = make([]*BookSearchOutput, 0, len(responseBooks))
		for _, book := range responseBooks {
			books = append(books, extractSearchResults(book))
		}
	}

	return books, nil
}
