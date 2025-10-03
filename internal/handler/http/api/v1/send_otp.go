package v1

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"io"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
)

func (h *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var req models.OTPRequest
	if err := req.UnmarshalBinary(body); err != nil {
		h.writeError(w, fmt.Errorf("Unable to parse otp request: %w: %w", err, apperror.ErrBadRequest))
		return
	}

	login := entity.EmailOrPhone(pointer.GetString(req.Login))

	err := h.uc.SendUserOTP(ctx, login)
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to send otp: %w", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
