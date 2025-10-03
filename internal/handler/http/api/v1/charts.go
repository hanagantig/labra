package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"labra/internal/apperror"
	"labra/internal/entity"
	"labra/internal/handler/http/api/v1/models"
	"labra/pkg/middleware"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func (h *Handler) Charts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Обработка preflight-запросов (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	param := mux.Vars(r)
	pUUID, err := uuid.Parse(param["profile_id"])
	if err != nil {
		h.writeError(w, fmt.Errorf("profile id is required: %w, %w", err, apperror.ErrBadRequest))

		return
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		h.writeError(w, fmt.Errorf("get user uuid from context: %w", err))

		return
	}

	from := time.Date(time.Now().Year(), time.Now().Month()-1, time.Now().Day(), 0, 0, 0, 0, time.UTC)

	qFrom := r.URL.Query().Get("from")
	qTo := r.URL.Query().Get("to")

	if qFrom != "" {
		fromTimestamp, err := strconv.ParseInt(qFrom, 10, 64)
		if err != nil {
			h.writeError(w, fmt.Errorf("unable to parse from date: %v: %w", err.Error(), http.StatusBadRequest))

			return
		}

		from = time.Unix(fromTimestamp, 0)
	}

	to := from.AddDate(0, 1, 0)
	if qTo != "" {
		toTimestamp, err := strconv.ParseInt(qTo, 10, 64)
		if err != nil {
			h.writeError(w, fmt.Errorf("unable to parse to date: %v: %w", err.Error(), http.StatusBadRequest))

			return
		}
		to = time.Unix(toTimestamp, 0)
	}

	markers, err := h.uc.GetMarkers(ctx, userUUID, pUUID, from, to, nil)
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to get results: %w", err))

		return
	}

	if len(markers) == 0 {
		w.Write([]byte("{}"))

		return
	}

	markersRes := markers.Group()
	data := make([]*models.TimeSeriesData, 0, len(markersRes))
	minMarkerValue := float64(0)
	maxMarkerValue := float64(0)
	for mKey, m := range markersRes {
		data = append(data, &models.TimeSeriesData{
			ID:    int64(mKey.Marker.ID),
			Spots: markersToSpots(m),
			Title: mKey.Marker.Name,
			Unit:  mKey.Unit.String(),
			Style: &models.ChartStyle{
				Opacity:      1,
				PrimaryColor: int64(mKey.Marker.PrimaryColor),
			},
		})

		minV := m.MinValue()
		if minV < minMarkerValue {
			minMarkerValue = minV
		}

		maxV := m.MaxValue()
		if maxV > maxMarkerValue {
			maxMarkerValue = maxV
		}
	}

	resp := models.Chart{
		Data: data,
		//Timestamps: times,
		Title: "Chart",
		Maxx:  float64(to.Unix()),   //float64(to.Add(48 * time.Hour).Unix()),
		Minx:  float64(from.Unix()), //float64(from.Add(-240 * time.Hour).Unix()),
		Miny:  minMarkerValue - 10,
		Maxy:  maxMarkerValue + 10,
		Annotations: []*models.ChartAnnotation{
			{
				Axis:  "x",
				Start: float64(1741963200),
				End:   float64(1742963200),
				Style: &models.ChartStyle{
					Opacity:        0.3,
					PrimaryColor:   74564575345,
					SecondaryColor: 563453453454,
				},
			},
		},
	}

	res, err := resp.MarshalBinary()
	if err != nil {
		h.writeError(w, fmt.Errorf("unable to marshal response: %w", err))

		return
	}

	w.Write(res)
}

func markersToSpots(res entity.MarkerResults) []*models.ChartSpot {
	sp := make([]*models.ChartSpot, 0, len(res))
	for _, r := range res {
		sp = append(sp, &models.ChartSpot{X: float64(r.CreatedAt.Unix()), Y: r.Value})
	}

	sort.Slice(sp, func(i, j int) bool {
		return sp[i].X < sp[j].X
	})

	return sp
}

func packData(buckets []int64, data map[int64][]float64) []*models.ChartSpot {
	packedValues := make([]float64, 0, len(buckets))

	for i := 0; i < len(buckets); i++ {
		valuesToPack := []float64{}
		for k, v := range data {
			cur := buckets[i]
			next := buckets[i]

			if i < len(buckets)-1 {
				next = buckets[i+1]
			}

			if k >= cur && k < next {
				valuesToPack = append(valuesToPack, v...)
			}
		}

		if len(valuesToPack) > 0 {
			packedValues = append(packedValues, valuesToPack[len(valuesToPack)-1])
		}
	}

	res := make([]*models.ChartSpot, 0, len(packedValues))
	for i := 0; i < len(packedValues); i++ {
		res = append(res, &models.ChartSpot{X: float64(buckets[i]), Y: packedValues[i]})
	}

	return res
}
