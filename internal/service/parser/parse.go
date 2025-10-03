package parser

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) Parse(ctx context.Context, labID int, txt string) (entity.CheckupResults, error) {
	var err error

	p, ok := s.parsers[labID]
	if !ok {
		p, err = s.loader.GetParserRules(ctx, labID)
		if err != nil {
			return entity.CheckupResults{}, nil
		}
	}

	return p.ParseResults(txt), nil
}
