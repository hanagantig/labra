package docupanda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"labra/internal/apperror"
	"net/http"
	"net/url"
)

func (r *Repository) Standardize(ctx context.Context, documentID string) (string, error) {
	apiURL, err := url.JoinPath(r.baseURL, "standardize/batch")
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"documentIds": []string{documentID},
		"schemaId":    "5188c4f3",
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", fmt.Sprintf("%s", r.apiKey))

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to recognize document: %s", string(body))
	}

	var res = struct {
		JobID              string   `json:"jobId"`
		Status             string   `json:"status"`
		StandardizationIds []string `json:"standardizationIds"`
		Details            string   `json:"details"`
	}{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	if len(res.StandardizationIds) == 0 {
		return "", fmt.Errorf("no standardizations found: %w", apperror.ErrNotFound)
	}

	return res.StandardizationIds[0], nil
}
