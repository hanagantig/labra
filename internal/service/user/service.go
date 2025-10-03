package user

import (
	"context"
	"github.com/google/uuid"
	"labra/internal/entity"
	"labra/internal/service"
)

type userRepo interface {
	service.Transactor
	GetByCredentials(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.User, error)
	GetByUUID(ctx context.Context, userUUID uuid.UUID) (entity.User, error)
	GetByID(ctx context.Context, userID int) (entity.User, error)
	GetByLogin(ctx context.Context, login entity.EmailOrPhone) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	SetVerifiedByID(ctx context.Context, userID int) error
	//Update(ctx context.Context, user entity.User) (entity.User, error)
}

type profileSvc interface {
	GetUserProfiles(ctx context.Context, userUUID uuid.UUID) (entity.Profiles, error)
	CreateByUser(ctx context.Context, user entity.User, linkedPatient entity.LinkedEntity) (entity.Profile, error)
	GetByID(ctx context.Context, userID, patientID int) (entity.Profile, error)
}

type otpService interface {
	GenerateOTP(ctx context.Context, userID int, otpType entity.OTPObjectType, objectID string) (entity.OTPCode, error)
	VerifyOTP(ctx context.Context, code entity.OTPCode) error
}

type notifyService interface {
	Notify(ctx context.Context, contact entity.Contact, templateID string, args map[string]string) error
}

type contactService interface {
	CreateContactByLogin(ctx context.Context, login entity.EmailOrPhone) (entity.Contact, error)
	GetByValue(ctx context.Context, login entity.EmailOrPhone) (entity.Contact, error)
	LinkUser(ctx context.Context, userID int, contactID int) error
}

type Service struct {
	userRepo   userRepo
	contactSvc contactService
	profileSvc profileSvc
	otpSvc     otpService
	notifySvc  notifyService
}

func NewService(ur userRepo, cs contactService, ps profileSvc, otp otpService, ns notifyService) *Service {
	return &Service{
		userRepo:   ur,
		contactSvc: cs,
		profileSvc: ps,
		otpSvc:     otp,
		notifySvc:  ns,
	}
}
