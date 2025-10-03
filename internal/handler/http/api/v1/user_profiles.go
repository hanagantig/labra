package v1

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"labra/pkg/middleware"
	"net/http"
)

func (h *Handler) UserProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	userUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		h.writeError(w, fmt.Errorf("%w, %w", err, apperror.ErrBadRequest))

		return
	}

	patients, err := h.uc.GetUserProfiles(ctx, userUID)
	if err != nil {
		h.writeError(w, fmt.Errorf("Unable to get user patients: %w", err))

		return
	}

	res := models.ProfilesResponse{}
	for _, p := range patients {
		res = append(res, &models.Profile{
			DateOfBirth: p.DateOfBirth.String(),
			Gender:      p.Gender,
			UUID:        pointer.ToString(p.Uuid.String()),
			ID:          p.Uuid.String(),
			FName:       p.FName,
			LName:       p.LName,
			//AssociatedUser: strfmt.UUID(p.UserID.String()),
			Access:   p.LinkedUserAccess,
			Contacts: buildContacts(p.Contacts),
		})
	}

	resp, err := swag.WriteJSON(res)

	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}

	w.Write(resp)
}

func buildContacts(contacts entity.Contacts) []*models.Contact {
	res := []*models.Contact{}
	for _, c := range contacts {
		res = append(res, &models.Contact{
			IsVerified: c.VerifiedByPatient(),
			Type:       string(c.Type),
			Value:      c.Value,
		})
	}

	return res
}
