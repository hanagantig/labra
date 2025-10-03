package notify

import (
	"context"
	"fmt"
	"labra/internal/entity"
)

type TemplateID string

type sender interface {
	Send(ctx context.Context, to, templateID string, args map[string]string) error
	GetType() entity.ContactType
}

type Service struct {
	registeredSenders map[entity.ContactType]sender
}

func NewService(senders ...sender) *Service {
	rs := make(map[entity.ContactType]sender)
	for _, s := range senders {
		rs[s.GetType()] = s
	}
	return &Service{
		registeredSenders: rs,
	}
}

func (s *Service) getSender(t entity.ContactType) (sender, error) {
	if sn, ok := s.registeredSenders[t]; ok {
		return sn, nil
	}

	return nil, fmt.Errorf("sender not found")
}
