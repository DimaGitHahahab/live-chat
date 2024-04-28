package adapter

import (
	"encoding/json"
	"io"
	"net/http"

	"chat-service/internal/domain"
)

type StorageClient struct {
	url    string
	client *http.Client
}

func NewStorageClient(url string) *StorageClient {
	return &StorageClient{
		url:    url,
		client: &http.Client{},
	}
}

func (s *StorageClient) FetchMessageHistory() ([]domain.Message, error) {
	resp, err := s.client.Get(s.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var messages []domain.Message
	err = json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		if err == io.EOF { // empty message history
			return []domain.Message{}, nil
		}
		return nil, err
	}

	return messages, err
}
