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

func (s *Service) RecognizeResults(ctx context.Context, fileType string, bytes []byte) (entity.CheckupResults, error) {
	return entity.CheckupResults{}, nil
	//text, err := s.scannerSvc.ScanByFileType(ctx, fileType, bytes)
	//if err != nil {
	//	return entity.CheckupResults{}, err
	//}
	//
	//recognizedCheckup, err := s.nerSvc.Recognize(ctx, text)
	//if err != nil {
	//	return entity.CheckupResults{}, err
	//}

	//res, err := s.parserSvc.Parse(ctx, 1, text)
	//if err != nil {
	//	return entity.CheckupResults{}, err
	//}

	//res := s.identifier.Identify(ctx, recognizedCheckup)

	//return res, nil
}
