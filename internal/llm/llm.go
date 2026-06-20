package llm

import (
	"context"
	"fmt"
	"strings"
)

type Suggestions struct {
	Title   string `json:"title"`
	Problem string `json:"problem"`
	Fix     string `json:"fix"`
}

type Provider interface {
	Name() string
	GenerateSuggestions(ctx context.Context, repo string, commands []string) (Suggestions, error)
}

func BuildPrompt(repo string, commands []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Analyze the following shell commands executed in the repository '%s' during a development session.\n", repo))
	sb.WriteString("Based on these commands, suggest:\n")
	sb.WriteString("1. A short, descriptive Title for this session (e.g., 'Configure SQLite FTS5 index').\n")
	sb.WriteString("2. A concise description of the Problem faced.\n")
	sb.WriteString("3. A concise description of the Fix implemented.\n\n")
	sb.WriteString("You MUST respond ONLY with a raw JSON object matching the following structure:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"title\": \"string\",\n")
	sb.WriteString("  \"problem\": \"string\",\n")
	sb.WriteString("  \"fix\": \"string\"\n")
	sb.WriteString("}\n\n")
	sb.WriteString("Commands list:\n")
	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf("- %s\n", cmd))
	}
	return sb.String()
}
