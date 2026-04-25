package collector

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/reski/openrouter-exporter/cache"
)

type OpenRouterCollector struct {
	cache  *cache.Cache
	logger *slog.Logger
}

func New(c *cache.Cache, logger *slog.Logger) *OpenRouterCollector {
	return &OpenRouterCollector{cache: c, logger: logger}
}

func (c *OpenRouterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- MetricInputPrice
	ch <- MetricOutputPrice
	ch <- MetricCacheReadPrice
	ch <- MetricUptime30m
	ch <- MetricUptime5m
	ch <- MetricUptime1d
	ch <- MetricThroughput
	ch <- MetricLatency
	ch <- MetricModelInfo
	ch <- MetricScrapeDuration
	ch <- MetricScrapeErrors
	ch <- MetricModelsScraped
	ch <- MetricEndpointsScraped
	ch <- MetricScrapeTimestamp
	ch <- MetricActivityRequests
	ch <- MetricActivityPromptTokens
	ch <- MetricActivityCompletionTokens
	ch <- MetricActivityToolCalls
	ch <- MetricActivityCacheHitTokens
	ch <- MetricActivityReasoningTokens
	ch <- MetricActivityReasoningRatio
	ch <- MetricActivityScrapeDuration
	ch <- MetricActivityScrapeErrors
	ch <- MetricActivityScrapeTimestamp
}

func (c *OpenRouterCollector) Collect(ch chan<- prometheus.Metric) {
	data := c.cache.Get()
	if data == nil {
		return
	}

	// Scrape metadata
	ch <- prometheus.MustNewConstMetric(MetricScrapeDuration, prometheus.GaugeValue, data.FetchDuration.Seconds())
	ch <- prometheus.MustNewConstMetric(MetricScrapeErrors, prometheus.GaugeValue, float64(data.FetchErrors))
	ch <- prometheus.MustNewConstMetric(MetricModelsScraped, prometheus.GaugeValue, float64(data.ModelsCount))
	ch <- prometheus.MustNewConstMetric(MetricEndpointsScraped, prometheus.GaugeValue, float64(data.EndpointsCount))
	ch <- prometheus.MustNewConstMetric(MetricScrapeTimestamp, prometheus.GaugeValue, float64(data.FetchedAt.Unix()))

	// Per-model info metrics
	for _, m := range data.Models {
		modality := m.Architecture.Modality
		tokenizer := m.Architecture.Tokenizer
		inputMod := strings.Join(m.Architecture.InputModalities, ",")
		outputMod := strings.Join(m.Architecture.OutputModalities, ",")

		ch <- prometheus.MustNewConstMetric(
			MetricModelInfo, prometheus.GaugeValue, 1,
			m.ID, m.Name, strconv.Itoa(m.ContextLength),
			modality, tokenizer, inputMod, outputMod,
		)
	}

	// Per-endpoint metrics — deduplicate by tag across models
	seen := make(map[string]bool)
	for modelID, epResp := range data.Endpoints {
		for _, ep := range epResp.Data.Endpoints {
			// Deduplicate: same tag can appear under different model IDs
			key := modelID + "/" + ep.Tag
			if seen[key] {
				continue
			}
			seen[key] = true

			c.emitEndpointMetrics(ch, modelID, ep)
		}
	}

	// Activity metrics
	if data.Activity != nil {
		ch <- prometheus.MustNewConstMetric(MetricActivityScrapeErrors, prometheus.GaugeValue, float64(data.ActivityErrors))

		for modelID, records := range data.Activity {
			for _, r := range records {
				date := strings.SplitN(r.Date, " ", 2)[0]
				ch <- prometheus.MustNewConstMetric(MetricActivityRequests, prometheus.GaugeValue, float64(r.Count), modelID, date)
				ch <- prometheus.MustNewConstMetric(MetricActivityPromptTokens, prometheus.GaugeValue, float64(r.TotalPromptTokens), modelID, date)
				ch <- prometheus.MustNewConstMetric(MetricActivityCompletionTokens, prometheus.GaugeValue, float64(r.TotalCompletionTokens), modelID, date)
				ch <- prometheus.MustNewConstMetric(MetricActivityToolCalls, prometheus.GaugeValue, float64(r.TotalToolCalls), modelID, date)
				ch <- prometheus.MustNewConstMetric(MetricActivityCacheHitTokens, prometheus.GaugeValue, float64(r.TotalNativeTokensCached), modelID, date)
				ch <- prometheus.MustNewConstMetric(MetricActivityReasoningTokens, prometheus.GaugeValue, float64(r.TotalNativeTokensReasoning), modelID, date)

				if r.TotalCompletionTokens > 0 {
					ratio := float64(r.TotalNativeTokensReasoning) / float64(r.TotalCompletionTokens)
					ch <- prometheus.MustNewConstMetric(MetricActivityReasoningRatio, prometheus.GaugeValue, ratio, modelID, date)
				}
			}
		}
	}
}

func (c *OpenRouterCollector) emitEndpointMetrics(ch chan<- prometheus.Metric, modelID string, ep cache.Endpoint) {
	provider := ep.ProviderName
	tag := ep.Tag
	quant := ep.Quantization

	// Pricing
	if price, err := parsePrice(ep.Pricing.Input); err == nil {
		ch <- prometheus.MustNewConstMetric(MetricInputPrice, prometheus.GaugeValue, price, modelID, provider, tag, quant)
	} else {
		c.logger.Warn("invalid input price", "model", modelID, "provider", provider, "value", ep.Pricing.Input, "error", err)
	}

	if price, err := parsePrice(ep.Pricing.Output); err == nil {
		ch <- prometheus.MustNewConstMetric(MetricOutputPrice, prometheus.GaugeValue, price, modelID, provider, tag, quant)
	} else {
		c.logger.Warn("invalid output price", "model", modelID, "provider", provider, "value", ep.Pricing.Output, "error", err)
	}

	if ep.Pricing.InputCacheRead != nil {
		if price, err := parsePrice(*ep.Pricing.InputCacheRead); err == nil {
			ch <- prometheus.MustNewConstMetric(MetricCacheReadPrice, prometheus.GaugeValue, price, modelID, provider, tag, quant)
		}
	}

	// Uptime
	ch <- prometheus.MustNewConstMetric(MetricUptime5m, prometheus.GaugeValue, ep.UptimeLast5m, modelID, provider, tag, quant)
	ch <- prometheus.MustNewConstMetric(MetricUptime30m, prometheus.GaugeValue, ep.UptimeLast30m, modelID, provider, tag, quant)
	ch <- prometheus.MustNewConstMetric(MetricUptime1d, prometheus.GaugeValue, ep.UptimeLast1d, modelID, provider, tag, quant)

	// Throughput (auth-only)
	if ep.ThroughputLast30m != nil {
		ch <- prometheus.MustNewConstMetric(MetricThroughput, prometheus.GaugeValue, ep.ThroughputLast30m.P50, modelID, provider, tag, quant, "p50")
		ch <- prometheus.MustNewConstMetric(MetricThroughput, prometheus.GaugeValue, ep.ThroughputLast30m.P75, modelID, provider, tag, quant, "p75")
		ch <- prometheus.MustNewConstMetric(MetricThroughput, prometheus.GaugeValue, ep.ThroughputLast30m.P90, modelID, provider, tag, quant, "p90")
		ch <- prometheus.MustNewConstMetric(MetricThroughput, prometheus.GaugeValue, ep.ThroughputLast30m.P99, modelID, provider, tag, quant, "p99")
	}

	// Latency (auth-only)
	if ep.LatencyLast30m != nil {
		ch <- prometheus.MustNewConstMetric(MetricLatency, prometheus.GaugeValue, ep.LatencyLast30m.P50*1000, modelID, provider, tag, quant, "p50")
		ch <- prometheus.MustNewConstMetric(MetricLatency, prometheus.GaugeValue, ep.LatencyLast30m.P75*1000, modelID, provider, tag, quant, "p75")
		ch <- prometheus.MustNewConstMetric(MetricLatency, prometheus.GaugeValue, ep.LatencyLast30m.P90*1000, modelID, provider, tag, quant, "p90")
		ch <- prometheus.MustNewConstMetric(MetricLatency, prometheus.GaugeValue, ep.LatencyLast30m.P99*1000, modelID, provider, tag, quant, "p99")
	}
}

func parsePrice(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty price")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse price %q: %w", s, err)
	}
	return v * 1_000_000, nil
}
