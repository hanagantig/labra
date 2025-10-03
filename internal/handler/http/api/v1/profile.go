package v1

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"labra/internal/handler/http/api/v1/models"
	"net/http"
	"time"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	markers, err := h.uc.Profile(ctx, uuid.UUID{})
	if err != nil {
		http.Error(w, "Unable to recognize: "+err.Error(), http.StatusInternalServerError)
		return
	}

	markersMap := map[string][]*models.Datapoint{}
	unitMap := map[string]string{}
	for _, marker := range markers {
		if _, ok := markersMap[marker.Marker.Name]; !ok {
			markersMap[marker.Marker.Name] = make([]*models.Datapoint, 0)
		}

		markersMap[marker.Marker.Name] = append(markersMap[marker.Marker.Name], &models.Datapoint{
			Value:     marker.Value,
			Timestamp: float64(marker.CreatedAt.Unix()),
		})

		unitMap[marker.Marker.Name] = marker.Unit.Name
	}

	data := make([]*models.TimeSeriesData, 0, len(markersMap))
	//for field, points := range markersMap {
	//	data = append(data, &models.TimeSeriesData{
	//		Datapoints: points,
	//		Field:      field,
	//		Unit:       unitMap[field],
	//	})
	//}

	resp := models.Account{
		User: &models.User{
			Gender:   "M",
			Username: "Tigran",
		},
		Charts: []*models.Chart{
			{
				Data:  data,
				Title: "Profile",
			},
		},
	}

	res, err := resp.MarshalBinary()
	if err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}

	w.Write(res)
}

func getTimeBuckets(from, to time.Time) []int64 {
	buckets := []int64{}
	minStep := time.Hour * 24
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

	delta := to.Sub(from)
	if delta.Hours()/24 > 14 {
		minStep = (time.Hour * 24) * 7
	}

	from = from.Add(-minStep)
	to = to.Add(minStep)

	for from.Before(to) {
		buckets = append(buckets, from.Unix())
		from = from.Add(minStep)
	}

	return buckets
}
