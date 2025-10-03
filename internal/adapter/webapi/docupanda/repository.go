package docupanda

import "net/http"

type Repository struct {
	client  http.Client
	apiKey  string
	baseURL string
}

func NewRepository(baseURL, apiKey string) *Repository {
	return &Repository{
		client:  http.Client{},
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}
