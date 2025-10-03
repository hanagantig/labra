package firebase

import (
	"firebase.google.com/go/v4/messaging"
)

type Repository struct {
	client *messaging.Client
}

func NewRepository(cl *messaging.Client) *Repository {
	return &Repository{cl}
}
