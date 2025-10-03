package uploader

import (
	"context"
	"labra/internal/entity"
)

type fileRegisterer interface {
	VerifyForUser(ctx context.Context, userID int, fingerprint string) error
	RegisterUploadedFile(ctx context.Context, f entity.UploadedFile) error
}

type fileStore interface {
	StoreFile(ctx context.Context, file entity.UploadedFile) (entity.UploadedFile, error)
}

type Service struct {
	registrySvc fileRegisterer
	storageSvc  fileStore
}

func NewService(rs fileRegisterer, fs fileStore) *Service {
	return &Service{
		registrySvc: rs,
		storageSvc:  fs,
	}
}
