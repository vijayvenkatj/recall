package repository

import (
	"context"
	"regexp"
	"strings"

	sqlc "github.com/vijayvenkatj/recall/internal/db/sqlc"
)

type Memory = sqlc.Memory

type MemoryRepository struct {
	queries *sqlc.Queries
}

type CreateMemoryParams struct {
	ID        string
	SessionID string
	Title     *string
	Summary   string
	CreatedAt int64
}

func NewMemoryRepository(queries *sqlc.Queries) *MemoryRepository {
	return &MemoryRepository{queries: queries}
}

func (r *MemoryRepository) Create(ctx context.Context, params CreateMemoryParams) (Memory, error) {
	return r.queries.CreateMemory(ctx, sqlc.CreateMemoryParams{
		ID:        params.ID,
		SessionID: params.SessionID,
		Title:     nullString(params.Title),
		Summary:   params.Summary,
		CreatedAt: params.CreatedAt,
	})
}

func (r *MemoryRepository) List(ctx context.Context, page Page) ([]Memory, error) {
	page = normalizePage(page)
	return r.queries.ListMemories(ctx, sqlc.ListMemoriesParams{
		Limit:  page.Limit,
		Offset: page.Offset,
	})
}

func (r *MemoryRepository) Search(ctx context.Context, query string, limit int32) ([]Memory, error) {
	formattedQuery := formatFTS5Query(query)
	if formattedQuery == "" {
		return nil, nil
	}

	rows, err := r.queries.SearchMemories(ctx, sqlc.SearchMemoriesParams{
		Query:    formattedQuery,
		LimitVal: int64(limit),
	})
	if err != nil {
		return nil, err
	}

	memories := make([]Memory, len(rows))
	for i, row := range rows {
		memories[i] = Memory{
			ID:        row.ID,
			SessionID: row.SessionID,
			Title:     row.Title,
			Summary:   row.Summary,
			CreatedAt: row.CreatedAt,
		}
	}
	return memories, nil
}

var wordRegexp = regexp.MustCompile(`[a-zA-Z0-9_]+`)

func formatFTS5Query(query string) string {
	words := wordRegexp.FindAllString(query, -1)
	if len(words) == 0 {
		return ""
	}

	var parts []string
	for _, w := range words {
		parts = append(parts, w+"*")
	}
	return strings.Join(parts, " AND ")
}
