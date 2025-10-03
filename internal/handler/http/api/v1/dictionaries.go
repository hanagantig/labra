package v1

import (
	"labra/internal/handler/http/api/v1/models"
	"net/http"
)

func (h *Handler) GetDictionaries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	dicts, err := h.uc.GetDictionaries(ctx)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Convert to API models
	var markerDefinitions []*models.MarkerDefinition
	for _, marker := range dicts.Markers {
		markerDef := &models.MarkerDefinition{
			ID:   int64(marker.ID),
			Name: marker.Name,
		}

		// Convert Range to ReferenceRange if available
		if marker.ReferenceRange.From != "" || marker.ReferenceRange.To != "" {
			refRange := &models.ReferenceRange{}
			//if marker.ReferenceRange.From != "" {
			//	if val, err := parseFloat(marker.ReferenceRange.From); err == nil {
			//		refRange.Min = &val
			//	}
			//}
			//if marker.ReferenceRange.To != "" {
			//	if val, err := parseFloat(marker.ReferenceRange.To); err == nil {
			//		refRange.Max = &val
			//	}
			//}
			markerDef.ReferenceRange = refRange
		}

		markerDefinitions = append(markerDefinitions, markerDef)
	}

	var apiUnits []*models.Unit
	for _, unit := range dicts.Units {
		apiUnits = append(apiUnits, &models.Unit{
			ID:          int64(unit.ID),
			Name:        unit.Name,
			FullName:    unit.FullName,
			Description: unit.Description,
		})
	}

	labs := []*models.Lab{}
	for _, lab := range dicts.Labs {
		labs = append(labs, &models.Lab{
			ID:   int64(lab.ID),
			Name: lab.Name,
		})
	}

	dictionaries := models.DictionariesResponse{
		Markers: markerDefinitions,
		Units:   apiUnits,
		Labs:    labs,
	}

	response, err := dictionaries.MarshalBinary()
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		h.logger.Error("Failed to write response")
	}
}
