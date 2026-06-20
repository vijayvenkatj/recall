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

type UpdateMemoryParams struct {
	ID      string
	Title   *string
	Summary string
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

func (r *MemoryRepository) Get(ctx context.Context, id string) (Memory, error) {
	return r.queries.GetMemory(ctx, id)
}

func (r *MemoryRepository) GetBySession(ctx context.Context, sessionID string) (Memory, error) {
	return r.queries.GetMemoryBySession(ctx, sessionID)
}

func (r *MemoryRepository) List(ctx context.Context, page Page) ([]Memory, error) {
	page = normalizePage(page)
	return r.queries.ListMemories(ctx, sqlc.ListMemoriesParams{
		Limit:  page.Limit,
		Offset: page.Offset,
	})
}

func (r *MemoryRepository) Update(ctx context.Context, params UpdateMemoryParams) (Memory, error) {
	return r.queries.UpdateMemory(ctx, sqlc.UpdateMemoryParams{
		Title:   nullString(params.Title),
		Summary: params.Summary,
		ID:      params.ID,
	})
}

func (r *MemoryRepository) Search(ctx context.Context, query string, limit int32) ([]Memory, error) {
	formattedQuery := formatFTS5Query(query)
	if formattedQuery == "" {
		return nil, nil
	}

	return r.queries.SearchMemories(ctx, sqlc.SearchMemoriesParams{
		Query:    formattedQuery,
		LimitVal: int64(limit),
	})
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteMemory(ctx, id)
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
