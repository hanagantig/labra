package docupanda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"labra/internal/entity"
	"net/http"
	"net/url"
)

func (r *Repository) GetResults(ctx context.Context, docID string) (entity.CheckupResults, error) {
	apiURL, err := url.JoinPath(r.baseURL, "standardization", docID)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", fmt.Sprintf("%s", r.apiKey))

	resp, err := r.client.Do(req)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.CheckupResults{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entity.CheckupResults{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	type ResResp struct {
		Data Report `json:"data"`
	}

	report := ResResp{}

	err = json.Unmarshal(body, &report)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	return report.Data.buildEntity(), nil
}
