package file

import (
	"context"
	"errors"
	"fmt"
	"labra/internal/apperror"
	"time"
)

const week = 7 * 24 * time.Hour
const filesPerWeek = 3

func (s *Service) VerifyForUser(ctx context.Context, userID int, fingerprint string) error {
	uploadedFiles, err := s.fileRepo.GetForUserByDuration(ctx, userID, week)
	if err != nil {
		return err
	}

	if len(uploadedFiles) >= filesPerWeek {
		return fmt.Errorf("uploaded files weekly limit exceeded: %w", apperror.ErrTooManyRequests)
	}

	existing, err := s.fileRepo.GetByFingerprint(ctx, fingerprint)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return err
	}

	if existing.Fingerprint != "" {
		return fmt.Errorf("file with fingerprint already exists")
	}

	return nil
}
