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

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.RefreshTokenRequest
	if err := req.UnmarshalBinary(body); err != nil {
		h.writeError(w, fmt.Errorf("unable to parse refresh token request: %w: %w", err, apperror.ErrBadRequest))
		return
	}

	if req.RefreshToken == nil || *req.RefreshToken == "" {
		h.writeError(w, fmt.Errorf("refresh token is required: %w", apperror.ErrBadRequest))
		return
	}

	rt := entity.NewRefreshTokenFromOpaque(pointer.GetString(req.RefreshToken))

	session, err := h.uc.RefreshUserToken(ctx, rt)
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to refresh user token: %w", err))
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
