package v1

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"labra/internal/apperror"
	"labra/pkg/middleware"
	"net/http"
)

func (h *Handler) ScanResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	param := mux.Vars(r)

	profileUUID, err := uuid.Parse(param["profile_id"])
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to parse profile id: %v: %w", err, apperror.ErrBadRequest))

		return
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		h.writeError(w, fmt.Errorf("%w, %w", err, apperror.ErrBadRequest))

		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		h.writeError(w, fmt.Errorf("Unable to parse multipart form: %v: %w", err, apperror.ErrBadRequest))

		return
	}

	file, _, err := r.FormFile("report")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		h.writeError(w, fmt.Errorf("Unable to read file: %w", http.StatusInternalServerError))

		return
	}

	fileBytes := buf.Bytes()

	contentType := http.DetectContentType(fileBytes)
	switch contentType {
	case "application/pdf", "image/jpeg", "image/png":
	default:
		h.writeError(w, fmt.Errorf("invalid file type: %s: %w", contentType, apperror.ErrBadRequest))

		return
	}

	_, err = h.uc.UploadCheckup(ctx, userUUID, profileUUID, http.DetectContentType(fileBytes), fileBytes)
	if err != nil {
		h.writeError(w, err)

		return
	}

	//resp := models.RecognizeResponse{
	//	Lab: &models.Lab{
	//		ID:   int64(checkupRes.Checkup.Lab.ID),
	//		Name: checkupRes.Checkup.Lab.Name,
	//	},
	//	Patient: &models.Patient{
	//		DateOfBirth: checkupRes.Checkup.Profile.DateOfBirth.String(),
	//		Gender:      checkupRes.Checkup.Profile.Gender,
	//		PatientName: pointer.ToString(checkupRes.Checkup.Profile.FName),
	//	},
	//	Date:    checkupRes.Checkup.Date.String(),
	//	Results: make([]*models.RecognizedMarker, 0, len(checkupRes.Results)),
	//}
	//
	//for _, m := range checkupRes.Results {
	//	rm := &models.RecognizedMarker{
	//		Marker: &models.Marker{
	//			ID:   int64(m.ID),
	//			Name: m.Name,
	//			Unit: &models.Unit{
	//				ID:   int64(m.Unit.ID),
	//				Name: m.Unit.String(),
	//			},
	//			Value: m.Value,
	//		},
	//	}
	//	if m.ID == 0 {
	//		rm.Errors = &models.ErrorMessage{
	//			Message: "Couldn't recognize marker name, please select it manually or remove",
	//		}
	//	}
	//	resp.Results = append(resp.Results, rm)
	//}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{}`))
	//err = json.NewEncoder(w).Encode("{}")
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
