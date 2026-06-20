package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiProvider struct {
	apiKey string
	model  string
}

func NewGeminiProvider(apiKey string, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
	}
}

func (p *GeminiProvider) Name() string {
	return "Gemini"
}

func (p *GeminiProvider) GenerateSuggestions(ctx context.Context, repo string, commands []string) (Suggestions, error) {
	if p.apiKey == "" {
		return Suggestions{}, fmt.Errorf("Gemini API key is not configured")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)

	prompt := BuildPrompt(repo, commands)

	payload := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"responseSchema": map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"title":   map[string]interface{}{"type": "STRING"},
					"problem": map[string]interface{}{"type": "STRING"},
					"fix":     map[string]interface{}{"type": "STRING"},
				},
				"required": []string{"title", "problem", "fix"},
			},
		},
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
		return Suggestions{}, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Suggestions{}, err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return Suggestions{}, fmt.Errorf("Gemini API returned empty response candidates")
	}

	var suggestions Suggestions
	if err := json.Unmarshal([]byte(response.Candidates[0].Content.Parts[0].Text), &suggestions); err != nil {
		return Suggestions{}, fmt.Errorf("failed to parse suggestions JSON from Gemini: %w", err)
	}

	return suggestions, nil
}

func (p *GeminiProvider) GenerateSummary(ctx context.Context, repo string, commands []string, problem string, fix string) (SummaryResult, error) {
	if p.apiKey == "" {
		return SummaryResult{}, fmt.Errorf("Gemini API key is not configured")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)

	prompt := BuildSummaryPrompt(repo, commands, problem, fix)

	payload := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"responseSchema": map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"title":   map[string]interface{}{"type": "STRING"},
					"summary": map[string]interface{}{"type": "STRING"},
				},
				"required": []string{"title", "summary"},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return SummaryResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return SummaryResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SummaryResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return SummaryResult{}, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return SummaryResult{}, err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return SummaryResult{}, fmt.Errorf("Gemini API returned empty response candidates")
	}

	var result SummaryResult
	if err := json.Unmarshal([]byte(response.Candidates[0].Content.Parts[0].Text), &result); err != nil {
		return SummaryResult{}, fmt.Errorf("failed to parse summary JSON from Gemini: %w", err)
	}

	return result, nil
}
