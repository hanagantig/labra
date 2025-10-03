package identifier

import (
	"context"
	"labra/internal/entity"
	"sync"
)

func (s *Service) Identify(ctx context.Context, res entity.CheckupResults) entity.CheckupResults {
	wg := sync.WaitGroup{}

	for _, ident := range s.identifiers {
		i := ident
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := i.Identify(ctx, res)
			if err != nil {
				// log error
			}
		}()
	}

	wg.Wait()

	return res
}
