package v1

import (
	"encoding/json"
	"go.uber.org/zap"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
)

func (h *Handler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	//strUserID := middleware.GetUserIDFromContext(ctx)
	//userID, err := uuid.Parse(strUserID)
	//if err != nil {
	//	http.Error(w, "Unable to get user id: "+err.Error(), http.StatusUnauthorized)
	//
	//	return
	//}
	//
	//param := mux.Vars(r)
	//patientID, err := uuid.Parse(param["id"])
	//if err != nil {
	//	http.Error(w, "patient id is required: "+err.Error(), http.StatusBadRequest)
	//	return
	//}

	err := h.uc.DeleteProfile(ctx, 123, 456)
	if err != nil {
		http.Error(w, "Unable to create a patient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.Success{}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
