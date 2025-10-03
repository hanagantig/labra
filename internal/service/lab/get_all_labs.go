package lab

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetAllLabs() ([]entity.Lab, error) {
	return s.labRepo.GetAll(context.Background())
}
