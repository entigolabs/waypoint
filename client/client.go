package client

import (
	"context"
	"fmt"
)

// ApiClient fetches data from the upstream data portal.
type ApiClient struct {
	http    *HttpClient
	baseURL string
}

// NewApiClient creates an ApiClient using the provided HttpClient.
func NewApiClient(baseURL string, http *HttpClient) *ApiClient {
	return &ApiClient{baseURL: baseURL, http: http}
}

// FetchCategories retrieves all categories from the upstream API.
func (c *ApiClient) FetchCategories(ctx context.Context) ([]Category, error) {
	var resp Response[[]Category]
	if _, err := c.http.GetAs(ctx, c.baseURL+"/v2/core/categories", nil, &resp, nil); err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return resp.Data, nil
}

// FetchEmsCategories retrieves all EMS categories from the upstream API.
func (c *ApiClient) FetchEmsCategories(ctx context.Context) ([]EMSCategory, error) {
	var resp Response[[]EMSCategory]
	if _, err := c.http.GetAs(ctx, c.baseURL+"/v2/core/ems-categories", nil, &resp, nil); err != nil {
		return nil, fmt.Errorf("failed to fetch ems-categories: %w", err)
	}
	return resp.Data, nil
}

// FetchEmsThemes retrieves all EMS themes from the upstream API.
func (c *ApiClient) FetchEmsThemes(ctx context.Context) ([]EMSTheme, error) {
	var resp Response[[]EMSTheme]
	if _, err := c.http.GetAs(ctx, c.baseURL+"/v2/core/ems-themes", nil, &resp, nil); err != nil {
		return nil, fmt.Errorf("failed to fetch ems-themes: %w", err)
	}
	return resp.Data, nil
}
