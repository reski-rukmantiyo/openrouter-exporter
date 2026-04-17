# OpenRouter Exporter

A [Prometheus](https://prometheus.io/) exporter for [OpenRouter](https://openrouter.ai/) model pricing, uptime, throughput, and latency metrics.

## Background

[OpenRouter](https://openrouter.ai/) provides a unified API to access hundreds of LLM models across providers (OpenAI, Anthropic, Google, Meta, etc.). Pricing and availability change frequently as new models are added and providers update their rates.

This exporter periodically scrapes the OpenRouter API and exposes model-level metrics in Prometheus format, enabling you to:

- Track pricing trends across models and providers
- Monitor endpoint uptime and reliability
- Observe throughput and latency characteristics (requires an API key)
- Build Grafana dashboards for cost estimation and model comparison

## Metrics

All pricing metrics carry the labels `model_id`, `provider_name`, `tag`, and `quantization`.

| Metric | Type | Description |
|---|---|---|
| `openrouter_model_prompt_price_dollars_per_million_tokens` | gauge | Prompt price ($/M tokens) |
| `openrouter_model_completion_price_dollars_per_million_tokens` | gauge | Completion price ($/M tokens) |
| `openrouter_model_input_cache_read_price_dollars_per_million_tokens` | gauge | Cached input price ($/M tokens) |
| `openrouter_endpoint_uptime_percentage_last_5m` | gauge | Uptime % over 5 minutes |
| `openrouter_endpoint_uptime_percentage_last_30m` | gauge | Uptime % over 30 minutes |
| `openrouter_endpoint_uptime_percentage_last_1d` | gauge | Uptime % over 1 day |
| `openrouter_endpoint_throughput_tokens_per_second` | gauge | Throughput in tokens/s (API key required) |
| `openrouter_endpoint_latency_milliseconds` | gauge | Latency in ms (API key required) |
| `openrouter_model_info` | gauge | Static metadata (value=1, carries labels) |
| `openrouter_scrape_duration_seconds` | gauge | Last scrape duration |
| `openrouter_scrape_errors_total` | gauge | Errors in last scrape |
| `openrouter_models_scraped` | gauge | Models in last scrape |
| `openrouter_endpoints_scraped` | gauge | Endpoints in last scrape |
| `openrouter_scrape_timestamp_seconds` | gauge | Timestamp of last successful scrape |

Throughput and latency metrics additionally carry a `quantile` label (`p50`, `p75`, `p90`, `p99`).

## Configuration

Each flag has a corresponding environment variable that takes precedence if set.

| Flag | Environment Variable | Default | Description |
|---|---|---|---|
| `-web.listen-address` | `OPENROUTER_LISTEN_ADDR` | `:9837` | Listen address |
| `-web.metrics-path` | `OPENROUTER_METRICS_PATH` | `/metrics` | Metrics endpoint path |
| `-scrape.interval` | `OPENROUTER_SCRAPE_INTERVAL` | `5m` | API scrape interval (min 10s) |
| `-api.timeout` | `OPENROUTER_API_TIMEOUT` | `30s` | HTTP client timeout |
| `-max-concurrency` | `OPENROUTER_MAX_CONCURRENCY` | `10` | Max concurrent endpoint fetches |
| `-api.key` | `OPENROUTER_API_KEY` | | OpenRouter API key (optional) |
| `-base-url` | `OPENROUTER_BASE_URL` | `https://openrouter.ai` | OpenRouter base URL |

## Quick Start

### Docker Compose

The easiest way to get started. Starts the exporter, Prometheus, and Grafana with a pre-configured dashboard.

```bash
# Optional: enable throughput/latency metrics
export OPENROUTER_API_KEY=sk-or-...

docker compose up -d
```

| Service | URL |
|---|---|
| Exporter | http://localhost:9837/metrics |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 (admin/admin) |

### Docker

Build and run the exporter alone:

```bash
docker build -t openrouter-exporter .

docker run -p 9837:9837 \
  -e OPENROUTER_API_KEY=sk-or-... \
  openrouter-exporter
```

### Binary

Download or build from source (requires [Go 1.25+](https://go.dev/dl/)):

```bash
go build -o openrouter-exporter .

./openrouter-exporter
```

## Building

### Current platform

```bash
make build
```

Output: `./bin/openrouter-exporter`

### Cross-compile all platforms

```bash
make build-all
```

Output binaries in `./bin/`:

```
bin/
  openrouter-exporter-linux-amd64
  openrouter-exporter-linux-arm64
  openrouter-exporter-darwin-amd64
  openrouter-exporter-darwin-arm64
  openrouter-exporter-windows-amd64.exe
  openrouter-exporter-windows-arm64.exe
```

Individual targets:

```bash
make build-linux       # linux amd64 + arm64
make build-darwin      # macos amd64 + arm64
make build-windows     # windows amd64 + arm64
```

### Clean

```bash
make clean
```

## Development

### Prerequisites

- Go 1.25+
- [Docker](https://docs.docker.com/get-docker/) (for containerized runs)

### Project structure

```
.
├── main.go              # HTTP server, entry point
├── config/config.go     # CLI flags and environment variables
├── client/client.go     # OpenRouter API client
├── cache/cache.go       # Caching layer with background refresh
├── collector/
│   ├── metrics.go       # Prometheus metric descriptor definitions
│   └── collector.go     # Metric collection and emission
├── examples/
│   ├── prometheus/      # Prometheus scrape configuration
│   └── grafana/         # Grafana dashboard and provisioning
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

### Running locally

```bash
# Build and run
go run . -api.key=sk-or-...

# Or with environment variables
OPENROUTER_API_KEY=sk-or-... go run .
```

### Endpoints

| Path | Description |
|---|---|
| `/metrics` | Prometheus metrics |
| `/health` | Health check (returns 200 when data is loaded) |
| `/` | Landing page with link to metrics |

## License

MIT
