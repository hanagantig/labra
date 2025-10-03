package usecase

import (
	"context"
	"errors"
)

func (u *UseCase) ProcessNewUploadedFiles(ctx context.Context) error {
	newFiles, err := u.fileSvc.GetNewFiles(ctx)
	if err != nil {
		return err
	}

	var resErr error
	for _, f := range newFiles {
		err = u.fileSvc.StartRecognition(ctx, f)
		resErr = errors.Join(resErr, err)
	}

	return resErr
}
