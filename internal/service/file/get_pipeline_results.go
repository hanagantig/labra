package file

import (
	"context"
	"labra/internal/entity"
)

func (s *Service) GetPipelineResults(ctx context.Context, f entity.UploadedFile) (entity.CheckupResults, error) {
	res, err := s.recognizerSvc.GetResults(ctx, f.PipelineID)
	targetStatus := entity.UploadedFileStatusRecognized
	if err != nil {
		targetStatus = f.Status
		f.Details = err.Error()
	}

	f.Status = targetStatus
	updateErr := s.fileRepo.UpdateByStatus(ctx, f, entity.UploadedFileStatusRecognizing)
	if updateErr != nil {
		return entity.CheckupResults{}, err
	}

	res.Checkup.Profile.UserID = f.UserID
	res.Checkup.Profile.ID = f.ProfileID
	res.Checkup.UploadedFileID = f.ID

	return res, err
}
