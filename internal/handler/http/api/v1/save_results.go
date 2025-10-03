package v1

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
	"io"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
	"time"
)

func (h *Handler) SaveResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.CheckupResult
	if err := req.UnmarshalBinary(body); err != nil {
		http.Error(w, "Unable to save: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err := req.Validate(strfmt.NewFormats())
	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusInternalServerError)

		return
	}

	layout := "2014-09-12T11:45:26.371Z"
	dateOfBirth, _ := time.Parse(layout, req.Profile.DateOfBirth)

	checkupResults := entity.CheckupResults{
		Checkup: entity.Checkup{
			Profile: entity.Profile{
				ID:          123,
				FName:       req.Profile.FName,
				DateOfBirth: dateOfBirth,
			},
			Lab: entity.Lab{
				ID:   1, //int(req.Lab.ID),
				Name: req.Lab.Name,
			},
			Date:    time.Now(), //req.Date,
			Comment: "",
		},
	}

	for _, res := range req.Results {
		if res.ID == 0 || res.Unit.ID == 0 {
			continue
		}

		marker := entity.MarkerResult{
			ID:     int(res.ID),
			Marker: entity.Marker{Name: res.Name},
			Value:  res.Value,
			Unit: entity.Unit{
				ID:   int(res.Unit.ID),
				Name: res.Unit.Name,
			},
		}

		checkupResults.Results = append(checkupResults.Results, marker)
	}

	err = h.uc.SaveResults(ctx, checkupResults)
	if err != nil {
		http.Error(w, "Save error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode("done")
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
