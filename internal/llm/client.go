package llm

import "strings"

func NewClient(provider string, apiKey string, model string, endpoint string) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "gemini":
		return NewGeminiProvider(apiKey, model), nil
	case "ollama":
		return NewOllamaProvider(endpoint, model), nil
	default:
		return nil, nil
	}
}
