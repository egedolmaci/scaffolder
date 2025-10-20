package llm

import "net/http"

type OpenAIClient struct {
	httpClient *http.Client
	apiKey string
	model string
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	if apiKey == "" {
		panic("API Key must be provided")
	}

	if model == "" {
		model = ""
	}

	return &OpenAIClient{
		httpClient: &http.Client{},
		apiKey: apiKey,
		model: model,
		
	}
}