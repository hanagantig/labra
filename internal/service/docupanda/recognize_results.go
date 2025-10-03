package recognizer

import (
	"context"
	"labra/internal/entity"
)

var markers = []string{
	"Гемоглобин",
	"Эритроциты",
	"Средний объём эритроцитов",
	"Среднее содержание Hb в эритроците",
	"Средняя концентрация Нь в эритроците",
	"Гетерогенность эритроцитов по объёму",
	"Гематокрит",
	"Тромбоциты",
	"Средний объём тромбоцитов",
	"Гетерогенность тромбоцитов по объёму",
	"Тромбокрит",
	"Лейкоциты",
	"Нейтрофилы",
	"Эозинофилы",
	"Базофилы",
	"Моноциты",
	"Лимфоциты",
	"Нейтрофилы",
	"Эозинофилы",
	"Базофилы",
	"Моноциты",
	"Лимфоциты",
}

func (s *Service) GetResults(ctx context.Context, stdID string) (entity.CheckupResults, error) {
	results, err := s.dpRepo.GetResults(ctx, stdID)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	return s.identSvc.Identify(ctx, results), nil
}
