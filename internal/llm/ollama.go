package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaProvider struct {
	endpoint string
	model    string
}

func NewOllamaProvider(endpoint string, model string) *OllamaProvider {
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	if model == "" {
		model = "llama3"
	}
	return &OllamaProvider{
		endpoint: endpoint,
		model:    model,
	}
}

func (p *OllamaProvider) Name() string {
	return "Ollama (" + p.model + ")"
}

func (p *OllamaProvider) GenerateSuggestions(ctx context.Context, repo string, commands []string) (Suggestions, error) {
	url := fmt.Sprintf("%s/api/generate", p.endpoint)

	prompt := BuildPrompt(repo, commands)

	payload := map[string]interface{}{
		"model":  p.model,
		"prompt": prompt,
		"format": "json",
		"stream": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Suggestions{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Suggestions{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Suggestions{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Suggestions{}, fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Suggestions{}, err
	}

	var suggestions Suggestions
	if err := json.Unmarshal([]byte(response.Response), &suggestions); err != nil {
		return Suggestions{}, fmt.Errorf("failed to parse suggestions JSON from Ollama: %w", err)
	}

	return suggestions, nil
}
