package v1

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"labra/pkg/middleware"
	"net/http"
	"strconv"
)

func (h *Handler) CheckupList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	param := mux.Vars(r)
	pUUID, err := uuid.Parse(param["profile_id"])
	if err != nil {
		h.writeError(w, fmt.Errorf("profile id is required: %w, %w", err, apperror.ErrBadRequest))

		return
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		h.writeError(w, fmt.Errorf("get user uuid from context: %w", err))

		return
	}

	search := r.URL.Query().Get("search")
	//filter, err := parseFilter(param)
	filter := entity.Filter{}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checkups, err := h.uc.GetCheckups(ctx, userUUID, pUUID, search, filter)
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to get list of checkups: %w", err))

		return
	}

	res := models.CheckupsResponse{}
	for _, checkup := range checkups {
		res = append(res, buildCheckupResponse(checkup))
	}
	resp, err := swag.WriteJSON(res)

	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}

	w.Write(resp)
}

func buildCheckupResponse(ch entity.CheckupResults) *models.CheckupWithResults {
	checkupRes := &models.Checkup{
		ID:    int64(ch.Checkup.ID),
		Title: ch.Checkup.Title,
		Date:  ch.Checkup.Date.Unix(),
		Lab: &models.Lab{
			ID:   int64(ch.Checkup.Lab.ID),
			Name: ch.Checkup.Lab.Name,
		},
		Profile: &models.Profile{
			//DateOfBirth: ch.Checkup.Profile.DateOfBirth.String(),
			Gender: ch.Checkup.Profile.Gender,
			UUID:   pointer.ToString(ch.Checkup.Profile.Uuid.String()),
			ID:     ch.Checkup.Profile.Uuid.String(),
			FName:  ch.Checkup.Profile.FName,
			LName:  ch.Checkup.Profile.LName,
		},
		Material: "кровь",
		Status:   string(ch.Checkup.Status),
		Tags: []*models.Tag{
			{
				ID:   1,
				Name: "важно",
			},
		},
	}

	if ch.Checkup.UploadedFileID > 0 {
		checkupRes.FileID = strconv.Itoa(ch.Checkup.UploadedFileID)
	}

	results := models.Markers{}
	for _, res := range ch.Results {
		results = append(results, &models.Marker{
			ID:               int64(res.Marker.ID),
			ResultID:         int64(res.ID),
			Name:             res.Marker.Name,
			UnrecognizedName: res.UnrecognizedName,
			Unit: &models.Unit{
				ID:               int64(res.Unit.ID),
				Name:             res.Unit.Name,
				UnrecognizedName: res.Unit.UnrecognizedName,
			},
			ReferenceRange: &models.ReferenceRange{
				Max: 73,
				Min: 15,
			},
			Value: res.Value,
			Date:  ch.Checkup.Date.Unix(),
		})
	}

	return &models.CheckupWithResults{
		Checkup: checkupRes,
		Results: results,
	}
}
