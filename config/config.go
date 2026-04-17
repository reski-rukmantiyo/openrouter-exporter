package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ListenAddress  string
	ScrapeInterval time.Duration
	APITimeout     time.Duration
	MaxConcurrency int
	APIKey         string
	BaseURL        string
	MetricsPath    string
}

func Load() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenAddress, "web.listen-address", ":9837", "Address to listen on")
	flag.DurationVar(&cfg.ScrapeInterval, "scrape.interval", 5*time.Minute, "How often to scrape the OpenRouter API")
	flag.DurationVar(&cfg.APITimeout, "api.timeout", 30*time.Second, "HTTP client timeout for API requests")
	flag.IntVar(&cfg.MaxConcurrency, "max-concurrency", 10, "Maximum concurrent endpoint fetches")
	flag.StringVar(&cfg.APIKey, "api.key", "", "OpenRouter API key (optional, enables throughput/latency metrics)")
	flag.StringVar(&cfg.BaseURL, "base-url", "https://openrouter.ai", "OpenRouter base URL")
	flag.StringVar(&cfg.MetricsPath, "web.metrics-path", "/metrics", "Path under which to expose metrics")

	flag.Parse()

	// Override from environment variables if set
	if v := os.Getenv("OPENROUTER_LISTEN_ADDR"); v != "" {
		cfg.ListenAddress = v
	}
	if v := os.Getenv("OPENROUTER_SCRAPE_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPENROUTER_SCRAPE_INTERVAL: %w", err)
		}
		cfg.ScrapeInterval = d
	}
	if v := os.Getenv("OPENROUTER_API_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPENROUTER_API_TIMEOUT: %w", err)
		}
		cfg.APITimeout = d
	}
	if v := os.Getenv("OPENROUTER_MAX_CONCURRENCY"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPENROUTER_MAX_CONCURRENCY: %w", err)
		}
		cfg.MaxConcurrency = n
	}
	if v := os.Getenv("OPENROUTER_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("OPENROUTER_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("OPENROUTER_METRICS_PATH"); v != "" {
		cfg.MetricsPath = v
	}

	// Validate
	if cfg.MaxConcurrency < 1 {
		return nil, fmt.Errorf("max-concurrency must be >= 1, got %d", cfg.MaxConcurrency)
	}
	if cfg.ScrapeInterval < 10*time.Second {
		return nil, fmt.Errorf("scrape.interval must be >= 10s, got %s", cfg.ScrapeInterval)
	}

	return cfg, nil
}
