package deepseek

import "net/http"

type Repository struct {
	client http.Client
	apiKey string
}

const baseURL = "https://api.deepseek.com"

func NewRepository(apiKey string) Repository {
	return Repository{
		client: http.Client{},
		apiKey: apiKey,
	}
}
