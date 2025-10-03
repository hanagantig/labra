package v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"labra/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type UseCase interface {
	UserSignIn(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.Session, error)
	UserSignUp(ctx context.Context, login entity.EmailOrPhone, pass entity.UserPassword) (entity.User, error)
	VerifyUserContact(ctx context.Context, login entity.EmailOrPhone, otpCode string) (entity.Session, error)
	UploadCheckup(ctx context.Context, userUUID, profileUUID uuid.UUID, fileType string, file []byte) (entity.CheckupResults, error)
	SaveResults(ctx context.Context, checkupResults entity.CheckupResults) error
	UpdateCheckup(ctx context.Context, userUUID uuid.UUID, checkupData entity.Checkup, resToAdd entity.MarkerResults, toDelete []int) error
	Profile(ctx context.Context, profileID uuid.UUID) (entity.MarkerResults, error)
	GetMarkers(ctx context.Context, userUUID, profileUUID uuid.UUID, from, to time.Time, names []string) (entity.MarkerResults, error)
	GetCheckups(ctx context.Context, userUUID, profileID uuid.UUID, search string, filter entity.Filter) ([]entity.CheckupResults, error)
	GetCheckup(ctx context.Context, checkupID int) (entity.CheckupResults, error)
	GetUserProfiles(ctx context.Context, userUUID uuid.UUID) (entity.Profiles, error)
	CreateProfile(ctx context.Context, userUUID uuid.UUID, patient entity.Profile) (entity.Profile, error)
	DeleteProfile(ctx context.Context, userID, profileID int) error
	SendUserOTP(ctx context.Context, login entity.EmailOrPhone) error
	RefreshUserToken(ctx context.Context, refreshToken entity.RefreshToken) (entity.Session, error)
	GetDictionaries(ctx context.Context) (entity.Dictionaries, error)
}

type Handler struct {
	uc     UseCase
	logger logger.Logger
	secret string
}

func NewHandler(uc UseCase, logs logger.Logger, secret string) *Handler {
	return &Handler{
		uc:     uc,
		logger: logs,
		secret: secret,
	}
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	errRes := models.APIError{
		Error: &models.ErrorMessage{Message: err.Error()},
	}
	switch {
	case errors.Is(err, apperror.ErrTooManyRequests):
		errRes.Code = http.StatusTooManyRequests
	case errors.Is(err, apperror.ErrNotFound):
		errRes.Code = http.StatusNotFound
	case errors.Is(err, apperror.ErrBadRequest) || errors.Is(err, apperror.ErrDuplicateEntity):
		errRes.Code = http.StatusBadRequest
	case errors.Is(err, apperror.ErrUnauthorized):
		errRes.Code = http.StatusUnauthorized
	default:
		errRes.Code = http.StatusInternalServerError
	}

	msg, _ := errRes.MarshalBinary()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(errRes.Code))
	w.Write(msg)
}

func parseFilter(params map[string]string) (entity.Filter, error) {
	fromTimestamp, err := strconv.Atoi(params["filter[from]"])
	if err != nil {
		return entity.Filter{}, err
	}

	toTimestamp, err := strconv.Atoi(params["filter[to]"])
	if err != nil {
		return entity.Filter{}, err
	}

	var categories []int

	err = json.Unmarshal([]byte(params["filter[categories]"]), &categories)
	if err != nil {
		return entity.Filter{}, err
	}

	var labIDs []int

	err = json.Unmarshal([]byte(params["filter[lab_ids]"]), &labIDs)
	if err != nil {
		return entity.Filter{}, err
	}

	var tags []string

	err = json.Unmarshal([]byte(params["filter[tags]"]), &tags)
	if err != nil {
		return entity.Filter{}, err
	}

	return entity.Filter{
		DateFrom:   time.Unix(int64(fromTimestamp), 0),
		DateTo:     time.Unix(int64(toTimestamp), 0),
		Categories: categories,
		LabIDs:     labIDs,
		Tags:       tags,
	}, nil
}
