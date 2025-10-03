package v1

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
	"io"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.SignUpRequest
	if err := req.UnmarshalBinary(body); err != nil {
		h.writeError(w, fmt.Errorf("Unable to parse signup request: %w, %w", err, apperror.ErrBadRequest))

		return
	}

	err := req.Validate(strfmt.Default)
	if err != nil {
		h.writeError(w, fmt.Errorf("Invalid request: %w, %w", err, apperror.ErrBadRequest))
		return
	}

	if !strfmt.IsEmail(*req.Login) {
		h.writeError(w, fmt.Errorf("Invalid request: login should be email, %w", apperror.ErrBadRequest))
		return
	}

	login := entity.EmailOrPhone(pointer.GetString(req.Login))
	pass := entity.UserPassword(req.Password.String())

	_, err = h.uc.UserSignUp(ctx, login, pass)
	if err != nil {
		h.writeError(w, fmt.Errorf("Unable to sign up: %w", err))
		return
	}

	res := models.Success{Success: "OK"}

	response, err := res.MarshalBinary()
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}

	w.Write(response)
}
