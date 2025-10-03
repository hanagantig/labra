package auth

import (
	"fmt"
	"labra/internal/apperror"
	"labra/internal/entity"
	"strconv"
)

func (s *Service) GetUserIDByToken(token entity.JWT) (int, error) {
	claim, err := token.ValidateAndGetClientClaims(s.accessTokenSecret)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(claim.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalud user id in token: %w", apperror.ErrUnauthorized)
	}

	return int(userID), nil
}
