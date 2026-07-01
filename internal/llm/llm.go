package llm

import (
	"context"
	"fmt"
	"strings"
)

type SummaryResult struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type Provider interface {
	Name() string
	GenerateSummary(ctx context.Context, repo string, commands []string, problem string, fix string) (SummaryResult, error)
}

func BuildSummaryPrompt(repo string, commands []string, problem string, fix string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Review a developer's debug/development session in the repository: '%s'.\n", repo))
	sb.WriteString("Here are the commands executed:\n")
	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf("- %s\n", cmd))
	}
	sb.WriteString(fmt.Sprintf("\nThe developer described the problem as:\n\"%s\"\n", problem))
	sb.WriteString(fmt.Sprintf("\nThe developer described the fix as:\n\"%s\"\n\n", fix))
	sb.WriteString("Based on this information, generate:\n")
	sb.WriteString("1. A short, professional Title for this memory (e.g., 'Fix Docker Port Bindings').\n")
	sb.WriteString("2. A professional, consolidated Summary explaining the problem, context from the command logs, and resolution.\n\n")
	sb.WriteString("You MUST respond ONLY with a raw JSON object matching the following structure:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"title\": \"string\",\n")
	sb.WriteString("  \"summary\": \"string\"\n")
	sb.WriteString("}\n")
	return sb.String()
}
