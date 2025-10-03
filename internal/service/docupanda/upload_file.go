package recognizer

import (
	"context"
	"labra/internal/entity"
)

const sourceName = "docupanda"

func (s *Service) StoreFile(ctx context.Context, file entity.UploadedFile) (entity.UploadedFile, error) {
	docID, err := s.dpRepo.UploadDocument(ctx, "U01P01", file.Bytes())
	if err != nil {
		return entity.UploadedFile{}, err
	}

	file.FileID = docID
	file.Source = sourceName

	return file, nil
}
