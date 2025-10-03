package file

import (
	"context"
	"labra/internal/entity"
	"time"
)

type recognizer interface {
	RunRecognizePipeline(ctx context.Context, f entity.UploadedFile) (string, error)
	GetResults(ctx context.Context, stdID string) (entity.CheckupResults, error)
	GetRecognizePipelineID(ctx context.Context, f entity.UploadedFile) (string, error)
}

type fileRepository interface {
	SaveUploadedFile(ctx context.Context, file entity.UploadedFile) error
	GetForUserByDuration(ctx context.Context, userID int, duration time.Duration) ([]entity.UploadedFile, error)
	GetByFingerprint(ctx context.Context, fingerprint string) (entity.UploadedFile, error)
	GetByStatus(ctx context.Context, status entity.UploadedFileStatus) ([]entity.UploadedFile, error)
	UpdateByStatus(ctx context.Context, file entity.UploadedFile, status entity.UploadedFileStatus) error
}
type Service struct {
	fileRepo      fileRepository
	recognizerSvc recognizer
}

func NewService(fileRepo fileRepository, rs recognizer) *Service {
	return &Service{
		fileRepo:      fileRepo,
		recognizerSvc: rs,
	}
}
