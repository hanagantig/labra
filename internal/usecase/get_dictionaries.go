package usecase

import (
	"context"
	"labra/internal/entity"
	"sync"
)

func (u *UseCase) GetDictionaries(ctx context.Context) (entity.Dictionaries, error) {
	// Get all available markers

	wg := sync.WaitGroup{}
	errChan := make(chan error, 3)

	dicts := entity.Dictionaries{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		markers, err := u.markerSvc.GetAllMarkers(ctx)
		if err != nil {
			errChan <- err
		}

		dicts.Markers = markers
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Get all available units
		units, err := u.unitSvc.GetAllUnits(ctx)
		if err != nil {
			errChan <- err
		}

		dicts.Units = units
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Get all available units
		labs, err := u.labSvc.GetAllLabs()
		if err != nil {
			errChan <- err
		}

		dicts.Labs = labs
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return entity.Dictionaries{}, err
		}
	}

	return dicts, nil
}
