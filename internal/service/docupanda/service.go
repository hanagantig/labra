package recognizer

import (
	"context"
	"labra/internal/entity"
)

type dpAPI interface {
	UploadDocument(ctx context.Context, dataset string, file []byte) (string, error)
	Standardize(ctx context.Context, documentID string) (string, error)
	GetResults(ctx context.Context, docID string) (entity.CheckupResults, error)
	GetStandardizations(ctx context.Context, fileID string) ([]string, error)
}

type identifier interface {
	Identify(ctx context.Context, res entity.CheckupResults) entity.CheckupResults
}

type Service struct {
	dpRepo   dpAPI
	identSvc identifier
}

func NewService(dpRepo dpAPI, is identifier) *Service {
	return &Service{
		dpRepo:   dpRepo,
		identSvc: is,
	}
}
