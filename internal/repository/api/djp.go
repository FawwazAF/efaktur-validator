package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/efaktur-validator/internal/model"
)

func (repo *Repository) GetInvoicesFromDJP(ctx context.Context, url string) (model.DJPEfaktur, error) {
	if url == "" {
		url = repo.mockURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.DJPEfaktur{}, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := repo.HTTPClient.Do(req)
	if err != nil {
		return model.DJPEfaktur{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	result := model.DJPEfaktur{}
	err = xml.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return model.DJPEfaktur{}, err
	}

	return result, nil
}
