package service

import (
	"context"
)

type Transactor interface {
	InTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}
