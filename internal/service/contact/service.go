package contact

import (
	"context"
	"labra/internal/entity"
	"labra/internal/service"
)

type contactRepository interface {
	service.Transactor
	Create(ctx context.Context, contact entity.Contact) (entity.Contact, error)
	AddLink(ctx context.Context, contactID int, entityType string, entityID string) error
	LinkFromEntity(ctx context.Context, fromEntityType, toEntityType string, fromEntityID, toEntityID int) error
	GetByValue(ctx context.Context, val string) (entity.Contact, error)
	GetByEntityIDs(ctx context.Context, entityType string, ids []int) (entity.Contacts, error)
	GetLinkedEntities(ctx context.Context, contactID int) ([]entity.LinkedEntity, error)
	SetVerifiedByEntity(ctx context.Context, linkedEntity entity.LinkedEntity) error
}

const (
	linkEntityTypeUser    = "user"
	linkEntityTypeProfile = "profile"
)

type Service struct {
	contactRepo contactRepository
}

func NewService(cr contactRepository) *Service {
	return &Service{
		contactRepo: cr,
	}
}
