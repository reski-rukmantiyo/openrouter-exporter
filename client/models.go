package client

type ModelsResponse struct {
	Data []Model `json:"data"`
}

type Model struct {
	ID                  string      `json:"id"`
	CanonicalSlug       string      `json:"canonical_slug"`
	HuggingFaceID       *string     `json:"hugging_face_id"`
	Name                string      `json:"name"`
	Created             int64       `json:"created"`
	Description         string      `json:"description"`
	ContextLength       int         `json:"context_length"`
	Architecture        Architecture `json:"architecture"`
	Pricing             Pricing     `json:"pricing"`
	TopProvider         TopProvider `json:"top_provider"`
	PerRequestLimits    interface{} `json:"per_request_limits"`
	SupportedParameters []string    `json:"supported_parameters"`
	DefaultParameters   interface{} `json:"default_parameters"`
	KnowledgeCutoff     *string     `json:"knowledge_cutoff"`
	ExpirationDate      *string     `json:"expiration_date"`
}

type Architecture struct {
	Modality         string   `json:"modality"`
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
	Tokenizer        string   `json:"tokenizer"`
	InstructType     *string  `json:"instruct_type"`
}

type Pricing struct {
	Prompt            string  `json:"prompt"`
	Completion        string  `json:"completion"`
	InputCacheRead    *string `json:"input_cache_read,omitempty"`
	InputCacheWrite   *string `json:"input_cache_write,omitempty"`
	WebSearch         *string `json:"web_search,omitempty"`
	Image             *string `json:"image,omitempty"`
	Audio             *string `json:"audio,omitempty"`
	InternalReasoning *string `json:"internal_reasoning,omitempty"`
}

type TopProvider struct {
	ContextLength       int  `json:"context_length"`
	MaxCompletionTokens *int `json:"max_completion_tokens"`
	IsModerated         bool `json:"is_moderated"`
}
