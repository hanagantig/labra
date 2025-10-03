package v1

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *Handler) CheckupDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	param := mux.Vars(r)
	checkupID, err := strconv.Atoi(param["id"])
	if err != nil {
		http.Error(w, "checkup id is required: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkup, err := h.uc.GetCheckup(ctx, checkupID)
	if err != nil {
		http.Error(w, "Unable to get checkup details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := buildCheckupResponse(checkup)

	resp, err := res.MarshalBinary()
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}

	w.Write(resp)
}
