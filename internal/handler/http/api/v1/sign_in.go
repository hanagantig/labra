package v1

import (
	"encoding/json"
	"github.com/AlekSi/pointer"
	"go.uber.org/zap"
	"io"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.SignInRequest
	if err := req.UnmarshalBinary(body); err != nil {
		http.Error(w, "Unable to save: "+err.Error(), http.StatusInternalServerError)

		return
	}

	login := entity.EmailOrPhone(pointer.GetString(req.Login))
	pass := entity.UserPassword(req.Password.String())

	session, err := h.uc.UserSignIn(ctx, login, pass)
	if err != nil {
		h.writeError(w, err)

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
