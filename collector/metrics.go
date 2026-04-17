package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	MetricPromptPrice = prometheus.NewDesc(
		"openrouter_model_prompt_price_dollars_per_million_tokens",
		"Price per million prompt tokens in USD",
		[]string{"model_id", "provider_name", "tag", "quantization"}, nil,
	)
	MetricCompletionPrice = prometheus.NewDesc(
		"openrouter_model_completion_price_dollars_per_million_tokens",
		"Price per million completion tokens in USD",
		[]string{"model_id", "provider_name", "tag", "quantization"}, nil,
	)
	MetricCacheReadPrice = prometheus.NewDesc(
		"openrouter_model_input_cache_read_price_dollars_per_million_tokens",
		"Price per million cached input tokens in USD",
		[]string{"model_id", "provider_name", "tag", "quantization"}, nil,
	)

	MetricUptime30m = prometheus.NewDesc(
		"openrouter_endpoint_uptime_percentage_last_30m",
		"Endpoint uptime percentage over the last 30 minutes",
		[]string{"model_id", "provider_name", "tag", "quantization"}, nil,
	)
	MetricUptime5m = prometheus.NewDesc(
		"openrouter_endpoint_uptime_percentage_last_5m",
		"Endpoint uptime percentage over the last 5 minutes",
		[]string{"model_id", "provider_name", "tag", "quantization"}, nil,
	)
	MetricUptime1d = prometheus.NewDesc(
		"openrouter_endpoint_uptime_percentage_last_1d",
		"Endpoint uptime percentage over the last day",
		[]string{"model_id", "provider_name", "tag", "quantization"}, nil,
	)

	MetricThroughput = prometheus.NewDesc(
		"openrouter_endpoint_throughput_tokens_per_second",
		"Throughput in tokens per second (requires API key)",
		[]string{"model_id", "provider_name", "tag", "quantization", "quantile"}, nil,
	)
	MetricLatency = prometheus.NewDesc(
		"openrouter_endpoint_latency_milliseconds",
		"Latency in milliseconds (requires API key)",
		[]string{"model_id", "provider_name", "tag", "quantization", "quantile"}, nil,
	)

	MetricModelInfo = prometheus.NewDesc(
		"openrouter_model_info",
		"Static model metadata as labels; value is always 1",
		[]string{"model_id", "model_name", "context_length", "modality", "tokenizer", "input_modalities", "output_modalities"}, nil,
	)

	MetricScrapeDuration = prometheus.NewDesc(
		"openrouter_scrape_duration_seconds",
		"Duration of the last cache refresh in seconds",
		nil, nil,
	)
	MetricScrapeErrors = prometheus.NewDesc(
		"openrouter_scrape_errors_total",
		"Total number of endpoint fetch errors in the last cache refresh",
		nil, nil,
	)
	MetricModelsScraped = prometheus.NewDesc(
		"openrouter_models_scraped",
		"Number of models in the last cache refresh",
		nil, nil,
	)
	MetricEndpointsScraped = prometheus.NewDesc(
		"openrouter_endpoints_scraped",
		"Number of endpoints in the last cache refresh",
		nil, nil,
	)
	MetricScrapeTimestamp = prometheus.NewDesc(
		"openrouter_scrape_timestamp_seconds",
		"Unix timestamp of the last successful cache refresh",
		nil, nil,
	)
)
