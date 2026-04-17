package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const userAgent = "openrouter-exporter/1.0"

type OpenRouterClient struct {
	baseURL        string
	apiKey         string
	httpClient     *http.Client
	maxConcurrency int
}

func NewClient(baseURL, apiKey string, timeout time.Duration, maxConcurrency int) *OpenRouterClient {
	return &OpenRouterClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxConcurrency: maxConcurrency,
	}
}

func (c *OpenRouterClient) FetchModels(ctx context.Context) (*ModelsResponse, error) {
	url := c.baseURL + "/api/v1/models"

	body, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch models: %w", err)
	}

	var resp ModelsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode models response: %w", err)
	}
	return &resp, nil
}

func (c *OpenRouterClient) FetchEndpoints(ctx context.Context, modelID string) (*EndpointsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/models/%s/endpoints", c.baseURL, modelID)

	body, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch endpoints for %s: %w", modelID, err)
	}

	var resp EndpointsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode endpoints for %s: %w", modelID, err)
	}
	return &resp, nil
}

type FetchAllResult struct {
	Endpoints map[string]*EndpointsResponse
	Errors    int
}

func (c *OpenRouterClient) FetchAllEndpoints(ctx context.Context, models []Model) (*FetchAllResult, error) {
	type result struct {
		modelID   string
		endpoints *EndpointsResponse
		err       error
	}

	jobs := make(chan string, len(models))
	results := make(chan result, len(models))

	// Launch workers
	var wg sync.WaitGroup
	for i := 0; i < c.maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				ep, err := c.FetchEndpoints(ctx, id)
				results <- result{modelID: id, endpoints: ep, err: err}
			}
		}()
	}

	// Send jobs
	go func() {
		for _, m := range models {
			jobs <- m.ID
		}
		close(jobs)
	}()

	// Close results when workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	endpoints := make(map[string]*EndpointsResponse, len(models))
	errorCount := 0
	for r := range results {
		if r.err != nil {
			errorCount++
			continue
		}
		if r.endpoints != nil && len(r.endpoints.Data.Endpoints) > 0 {
			endpoints[r.modelID] = r.endpoints
		}
	}

	return &FetchAllResult{
		Endpoints: endpoints,
		Errors:    errorCount,
	}, nil
}

func (c *OpenRouterClient) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if len(body) == 0 {
		return nil, nil
	}

	return body, nil
}
