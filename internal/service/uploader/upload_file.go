package uploader

import (
	"context"
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) UploadFile(ctx context.Context, profile entity.Profile, fileType string, bytes []byte) (entity.UploadedFile, error) {
	//TODO: in transaction
	uploadedFile := entity.NewUploadedFile(bytes)
	uploadedFile.UserID = profile.UserID
	uploadedFile.ProfileID = profile.ID

	err := s.registrySvc.VerifyForUser(ctx, uploadedFile.UserID, uploadedFile.Fingerprint)
	if err != nil {
		return entity.UploadedFile{}, fmt.Errorf("unable to verify file: %w: %w", err, apperror.ErrBadRequest)
	}

	uploadedFile, err = s.storageSvc.StoreFile(ctx, uploadedFile)
	if err != nil {
		return entity.UploadedFile{}, fmt.Errorf("unable to store file: %w", err)
	}

	uploadedFile.FileType = fileType
	err = s.registrySvc.RegisterUploadedFile(ctx, uploadedFile)
	if err != nil {
		return entity.UploadedFile{}, fmt.Errorf("unable to register uploaded file: %w", err)
	}

	return uploadedFile, nil
}
