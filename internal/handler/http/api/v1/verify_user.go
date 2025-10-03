package v1

import (
	"encoding/json"
	"fmt"
	"github.com/AlekSi/pointer"
	"go.uber.org/zap"
	"io"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
)

func (h *Handler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.VerifyContactRequest
	if err := req.UnmarshalBinary(body); err != nil {
		h.writeError(w, fmt.Errorf("Unable to parse signup request: %w: %w", err, apperror.ErrBadRequest))
		return
	}

	login := entity.EmailOrPhone(pointer.GetString(req.Login))

	session, err := h.uc.VerifyUserContact(ctx, login, pointer.GetString(req.Code))
	if err != nil {
		h.writeError(w, fmt.Errorf("Unable to verify user: %w", err))
		return
	}

	res := models.APIAuthToken{
		RefreshToken: session.AuthTokens.RefreshToken.Plain(),
		AccessToken:  session.AuthTokens.AccessToken.String(),
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
