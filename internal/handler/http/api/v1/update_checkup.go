package v1

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"labra/pkg/middleware"
	"net/http"
	"time"
)

func (h *Handler) UpdateCheckup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.CheckupUpdateRequest
	if err := req.UnmarshalBinary(body); err != nil {
		h.writeError(w, fmt.Errorf("unable to update checkup: %w", err))

		return
	}

	req.Checkup.Profile.UUID = pointer.ToString(req.Checkup.Profile.ID)
	err := req.Validate(strfmt.NewFormats())
	if err != nil {
		h.writeError(w, fmt.Errorf("Invalid request: %v", err.Error(), http.StatusBadRequest))

		return
	}

	reqID := pointer.GetString(req.Checkup.Profile.UUID)
	if reqID == "" {
		reqID = req.Checkup.Profile.ID
	}

	profileUUID, err := uuid.Parse(reqID)
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to parse profile id: %v: %w", err, apperror.ErrBadRequest))

		return
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		h.writeError(w, fmt.Errorf("%w, %w", err, apperror.ErrBadRequest))

		return
	}

	checkupDate := time.Unix(req.Checkup.Date, 0)
	if req.Checkup.Date <= 0 {
		h.writeError(w, fmt.Errorf("%w, %w", err, apperror.ErrBadRequest))

		return
	}

	checkupData := entity.Checkup{
		ID:    int(req.Checkup.ID),
		Title: req.Checkup.Title,
		Profile: entity.Profile{
			Uuid: profileUUID,
		},
		Lab: entity.Lab{
			ID:   int(req.Checkup.Lab.ID),
			Name: req.Checkup.Lab.Name,
		},
		Date:    checkupDate,
		Comment: "",
	}

	resultsToAdd := make(entity.MarkerResults, 0, len(req.ResultsToAdd))
	for _, res := range req.ResultsToAdd {
		if res.ID == 0 || res.Unit.ID == 0 {
			h.writeError(w, fmt.Errorf("can't add unrecognized marker or unit, %w", apperror.ErrBadRequest))

			return
		}

		marker := entity.MarkerResult{
			ID:     int(res.ID),
			Marker: entity.Marker{Name: res.Name},
			Value:  res.Value,
			Unit: entity.Unit{
				ID: int(res.Unit.ID),
			},
		}

		resultsToAdd = append(resultsToAdd, marker)
	}

	toDelete := make([]int, 0, len(req.ResultsIdsToDelete))
	for _, res := range req.ResultsIdsToDelete {
		toDelete = append(toDelete, int(res))
	}

	err = h.uc.UpdateCheckup(ctx, userUUID, checkupData, resultsToAdd, toDelete)
	if err != nil {
		h.writeError(w, fmt.Errorf("update checkup error, %w", err))

		return
	}

	res, err := h.uc.GetCheckup(ctx, checkupData.ID)
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to get updated checkup: %w", err))

		return
	}

	response := buildCheckupResponse(res)
	resp, err := response.MarshalBinary()
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
