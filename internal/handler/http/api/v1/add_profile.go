package v1

import (
	"encoding/json"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
	"io"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"labra/pkg/middleware"
	"net/http"
)

func (h *Handler) AddProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	userID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		h.writeError(w, fmt.Errorf("%w, %w", err, apperror.ErrBadRequest))

		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.Person
	if err := req.UnmarshalBinary(body); err != nil {
		h.writeError(w, fmt.Errorf("Unable to parse add profile request: %w, %w", err, apperror.ErrBadRequest))

		return
	}

	err = req.Validate(strfmt.NewFormats())
	if err != nil {
		h.writeError(w, fmt.Errorf("Invalid request: %w, %w", err, apperror.ErrBadRequest))

		return
	}

	p := entity.Profile{
		FName:  pointer.Get(req.FirstName),
		Gender: pointer.Get(req.Gender),
	}

	_, err = h.uc.CreateProfile(ctx, userID, p)
	if err != nil {
		h.writeError(w, fmt.Errorf("Unable to create a profile: %w", err))

		return
	}

	res := models.Success{}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
