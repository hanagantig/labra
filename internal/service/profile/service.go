package profile

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
	"labra/internal/service"
)

type profileRepo interface {
	service.Transactor
	GetByUserID(ctx context.Context, userID int) (entity.Profiles, error)
	GetByUserUUID(ctx context.Context, userUUID uuid.UUID) (entity.Profiles, error)
	GetByAssociatedUserID(ctx context.Context, userID int) (entity.Profile, error)
	Upsert(ctx context.Context, patient entity.Profile) (entity.Profile, error)
	UpdateAssociatedUser(ctx context.Context, patientID, userID int) error
	BindToUser(ctx context.Context, patientID, userID int, role string) error
	UnBindFromUser(ctx context.Context, patientID, userID int) error
	UpdateBoundAccess(ctx context.Context, patientID int, fromAccessLevel, toAccessLevel string) error
	SoftDelete(ctx context.Context, patientID int) error
	GetByID(ctx context.Context, userID, patientID int) (entity.Profile, error)
}

type contactService interface {
	LinkUserContactsToProfile(ctx context.Context, profileID, userID int) error
	GetForProfiles(ctx context.Context, patientIDs []int) (entity.Contacts, error)
}

type Service struct {
	profileRepo profileRepo
	contactSvc  contactService
}

func NewService(pr profileRepo, cs contactService) *Service {
	return &Service{
		profileRepo: pr,
		contactSvc:  cs,
	}
}
