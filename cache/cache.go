package cache

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/reski/openrouter-exporter/client"
)

type CachedData struct {
	Models         []Model
	Endpoints      map[string]*EndpointsResponse
	FetchedAt      time.Time
	FetchDuration  time.Duration
	FetchErrors    int
	ModelsCount    int
	EndpointsCount int
}

// Re-export client types for cache consumers
type Model = client.Model
type EndpointsResponse = client.EndpointsResponse
type Endpoint = client.Endpoint

type Cache struct {
	mu       sync.RWMutex
	data     *CachedData
	client   *client.OpenRouterClient
	interval time.Duration
	stopCh   chan struct{}
	logger   *slog.Logger
}

func New(c *client.OpenRouterClient, interval time.Duration, logger *slog.Logger) *Cache {
	return &Cache{
		client:   c,
		interval: interval,
		stopCh:   make(chan struct{}),
		logger:   logger,
	}
}

func (c *Cache) Start(ctx context.Context) error {
	if err := c.refresh(ctx); err != nil {
		return err
	}

	go c.run(ctx)
	return nil
}

func (c *Cache) Get() *CachedData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}

func (c *Cache) Stop() {
	close(c.stopCh)
}

func (c *Cache) run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.refresh(ctx); err != nil {
				c.logger.Error("cache refresh failed", "error", err)
			}
		case <-c.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (c *Cache) refresh(ctx context.Context) error {
	start := time.Now()

	modelsResp, err := c.client.FetchModels(ctx)
	if err != nil {
		return fmt.Errorf("fetch models: %w", err)
	}
	models := modelsResp.Data

	result, err := c.client.FetchAllEndpoints(ctx, models)
	if err != nil {
		return fmt.Errorf("fetch endpoints: %w", err)
	}

	// Count total endpoints
	totalEndpoints := 0
	for _, ep := range result.Endpoints {
		totalEndpoints += len(ep.Data.Endpoints)
	}

	cached := &CachedData{
		Models:         models,
		Endpoints:      result.Endpoints,
		FetchedAt:      time.Now(),
		FetchDuration:  time.Since(start),
		FetchErrors:    result.Errors,
		ModelsCount:    len(models),
		EndpointsCount: totalEndpoints,
	}

	c.mu.Lock()
	c.data = cached
	c.mu.Unlock()

	c.logger.Info("cache refreshed",
		"models", cached.ModelsCount,
		"endpoints", cached.EndpointsCount,
		"errors", cached.FetchErrors,
		"duration", cached.FetchDuration.Round(time.Millisecond),
	)

	return nil
}
