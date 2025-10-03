package verifier

import (
	"context"
	"labra/internal/entity"
)

type otpService interface {
	GenerateOTP(ctx context.Context, userID int, otpType entity.OTPObjectType, objectID string) (entity.OTPCode, error)
	GetOTP(ctx context.Context, objType entity.OTPObjectType, objID int, code string) (entity.OTPCode, error)
	VerifyOTP(ctx context.Context, code entity.OTPCode) error
}

type notifyService interface {
	Notify(ctx context.Context, contact entity.Contact, templateID string, args map[string]string) error
}

type profileService interface {
	GetAssociatedProfile(ctx context.Context, userID int) (entity.Profile, error)
	MigrateToUser(ctx context.Context, profileID, userID int) error
	CreateForUser(ctx context.Context, userID int, profile entity.Profile) (entity.Profile, error)
}

type contactService interface {
	VerifyForUser(ctx context.Context, userID int, contact entity.Contact, associatedPatient entity.Profile) (entity.Contact, error)
}

type Service struct {
	otpSvc     otpService
	notifySvc  notifyService
	profileSvc profileService
	contactSvc contactService
}

func NewService(os otpService, ns notifyService, ps profileService, cs contactService) *Service {
	return &Service{
		otpSvc:     os,
		notifySvc:  ns,
		profileSvc: ps,
		contactSvc: cs,
	}
}
