package apperror

import (
	"fmt"
)

var (
	ErrNotFound              = newClientError("item not found")
	ErrBadRequest            = newClientError("bad request")
	ErrDuplicateEntity       = newClientError("entity already exists")
	ErrEntityAlreadyVerified = newClientError("entity already verified")
	ErrUnauthorized          = newClientError("Unauthorized")
	ErrTooManyRequests       = newClientError("Too Many Requests")
)

type ErrClient string

func (e ErrClient) Error() string {
	return string(e)
}

func newClientError(msg string) error {
	return ErrClient(fmt.Sprintf("%s: client error", msg))
}
