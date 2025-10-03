package file

import (
	"context"
	"errors"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (s *Service) StartRecognition(ctx context.Context, f entity.UploadedFile) error {
	//TODO: in transaction
	pipelineID, err := s.recognizerSvc.GetRecognizePipelineID(ctx, f)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return err
	}

	if pipelineID == "" {
		pipelineID, err = s.recognizerSvc.RunRecognizePipeline(ctx, f)
	}

	targetStatus := entity.UploadedFileStatusRecognizing
	if err != nil {
		targetStatus = f.Status
		f.Details = err.Error()
	}

	f.PipelineID = pipelineID
	f.Status = targetStatus
	updateErr := s.fileRepo.UpdateByStatus(ctx, f, entity.UploadedFileStatusNew)
	if updateErr != nil {
		return updateErr
	}

	return err
}
