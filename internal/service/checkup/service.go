package checkup

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
	"labra/internal/service"
)

type checkupRepo interface {
	service.Transactor
	Save(ctx context.Context, checkup entity.Checkup) (entity.Checkup, error)
	UpdateCheckup(ctx context.Context, e entity.Checkup) (entity.Checkup, error)
	GetCheckupsByUUID(ctx context.Context, profileID uuid.UUID) (entity.Checkups, error)
	GetByID(ctx context.Context, checkupID int) (entity.Checkup, error)
	GetByFile(ctx context.Context, f entity.UploadedFile) (entity.Checkup, error)
}

type resultsRepo interface {
	Save(ctx context.Context, checkupID int, markers entity.MarkerResults) (entity.MarkerResults, error)
	DeleteByID(ctx context.Context, checkupID int, resultIDs []int) error
	GetByProfile(ctx context.Context, profileID int, filter entity.MarkerFilter, limit int) (entity.MarkerResults, error)
	GetByCheckup(ctx context.Context, checkupID int) (entity.MarkerResults, error)
	GetByCheckups(ctx context.Context, checkupIDs []int, markerIDs []int) (map[int]entity.MarkerResults, error)
}

type markersRepo interface {
	SearchByName(ctx context.Context, search string) (entity.Markers, error)
}

type Service struct {
	checkupRepo checkupRepo
	resultsRepo resultsRepo
	mr          markersRepo
}

func NewService(checkupRepo checkupRepo, resultsRepo resultsRepo, mr markersRepo) *Service {
	return &Service{
		checkupRepo: checkupRepo,
		resultsRepo: resultsRepo,
		mr:          mr,
	}
}
