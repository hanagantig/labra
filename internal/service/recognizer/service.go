package recognizer

import (
	"context"
	"labra/internal/entity"
)

type identifier interface {
	Identify(ctx context.Context, res entity.CheckupResults) entity.CheckupResults
}

type scannerSvc interface {
	ScanByFileType(ctx context.Context, fType string, bytes []byte) (string, error)
}

type parserSvc interface {
	Parse(ctx context.Context, labID int, txt string) (entity.CheckupResults, error)
}

type recognizer interface {
	Recognize(ctx context.Context, text string) (entity.CheckupResults, error)
}

type Service struct {
	scannerSvc scannerSvc
	parserSvc  parserSvc
	nerSvc     recognizer
	identifier identifier
}

func NewService(pSvc parserSvc, scanSvc scannerSvc, idf identifier, ner recognizer) *Service {
	return &Service{
		parserSvc:  pSvc,
		scannerSvc: scanSvc,
		identifier: idf,
		nerSvc:     ner,
	}
}
