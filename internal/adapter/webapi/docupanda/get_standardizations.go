package docupanda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"labra/internal/apperror"
	"net/http"
	"net/url"
)

func (r *Repository) GetStandardizations(ctx context.Context, fileID string) ([]string, error) {
	urlString, err := url.JoinPath(r.baseURL, "standardizations")
	if err != nil {
		return nil, err
	}

	apiURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	params := apiURL.Query()
	params.Set("document_id", fileID)
	apiURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", fmt.Sprintf("%s", r.apiKey))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, apperror.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res = []struct {
		StandardizationId string `json:"standardizationId"`
	}{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no standardizations for document: %w", apperror.ErrNotFound)
	}

	stds := make([]string, 0, len(res))
	for _, st := range res {
		stds = append(stds, st.StandardizationId)
	}

	return stds, nil
}
