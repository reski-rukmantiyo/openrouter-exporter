package client

type EndpointsResponse struct {
	Data EndpointsData `json:"data"`
}

type EndpointsData struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Created     int64      `json:"created"`
	Description string     `json:"description"`
	Endpoints   []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name                    string               `json:"name"`
	ModelID                 string               `json:"model_id"`
	ModelName               string               `json:"model_name"`
	ContextLength           int                  `json:"context_length"`
	Pricing                 EndpointPricing      `json:"pricing"`
	ProviderName            string               `json:"provider_name"`
	Tag                     string               `json:"tag"`
	Quantization            string               `json:"quantization"`
	MaxCompletionTokens     *int                 `json:"max_completion_tokens"`
	MaxPromptTokens         *int                 `json:"max_prompt_tokens"`
	SupportedParameters     []string             `json:"supported_parameters"`
	Status                  int                  `json:"status"`
	UptimeLast30m           float64              `json:"uptime_last_30m"`
	UptimeLast5m            float64              `json:"uptime_last_5m"`
	UptimeLast1d            float64              `json:"uptime_last_1d"`
	SupportsImplicitCaching bool                 `json:"supports_implicit_caching"`
	LatencyLast30m          *LatencyPercentiles  `json:"latency_last_30m"`
	ThroughputLast30m       *ThroughputPercentiles `json:"throughput_last_30m"`
}

type EndpointPricing struct {
	Input          string  `json:"prompt"`
	Output         string  `json:"completion"`
	InputCacheRead *string `json:"input_cache_read,omitempty"`
	Discount       float64 `json:"discount"`
}

type LatencyPercentiles struct {
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P99 float64 `json:"p99"`
}

type ThroughputPercentiles struct {
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P99 float64 `json:"p99"`
}
