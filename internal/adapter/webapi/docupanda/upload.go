package docupanda

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type CreateDocumentRequest struct {
	Document Document `json:"document"`
	Dataset  string   `json:"dataset"`
}

func (r *Repository) UploadDocument(ctx context.Context, dataset string, file []byte) (string, error) {
	apiURL, err := url.JoinPath(r.baseURL, "document")
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(file)

	payload := CreateDocumentRequest{
		Document: Document{
			File: File{
				Contents: encoded,
				Filename: "mgo_test",
			},
		},
		Dataset: dataset,
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

	var res map[string]string

	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, res["detail"])
	}

	return res["documentId"], nil
}
