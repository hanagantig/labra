package usecase

import (
	"context"
	"errors"
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (u *UseCase) ProcessRecognizedFiles(ctx context.Context) error {
	newFiles, err := u.fileSvc.GetRecognizingFiles(ctx)
	if err != nil {
		return err
	}

	var allErr error
	for _, f := range newFiles {
		err = u.processFile(ctx, f)
		allErr = errors.Join(allErr, err)
	}

	return allErr
}

func (u *UseCase) processFile(ctx context.Context, f entity.UploadedFile) error {
	ch, err := u.checkupSvc.GetCheckupByFile(ctx, f)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return fmt.Errorf("get checkup by file %d: %w", f.ID, err)
	}

	if ch.ID > 0 {
		errDuplicate := fmt.Errorf("checkup for file already exists: %w", apperror.ErrDuplicateEntity)
		f.Details = errDuplicate.Error()
		if err := u.fileSvc.MarkAsDuplicate(ctx, f); err != nil {
			// Record both signals: the duplicate condition + marking failure
			return errors.Join(errDuplicate, fmt.Errorf("mark duplicate %d: %w", f.ID, err))
		}

		return errDuplicate
	}

	// TODO: int transaction
	res, err := u.fileSvc.GetPipelineResults(ctx, f)
	if err != nil {
		return fmt.Errorf("get pipeline results %d: %w", f.ID, err)
	}

	res.Checkup.Status = entity.UnverifiedCheckupStatus
	if err := u.checkupSvc.RegisterCheckup(ctx, res); err != nil {
		return fmt.Errorf("register checkup for file #%d: %w", f.ID, err)
	}

	args := map[string]string{
		"content": "file recognized successfully!",
	}

	return u.notifySvc.Notify(ctx, entity.Contact{}, "file_recognized", args)
}
