package repository

import (
	"context"

	sqlc "github.com/vijayvenkatj/recall/internal/db/sqlc"
)

type Session = sqlc.Session

type SessionRepository struct {
	queries *sqlc.Queries
}

type CreateSessionParams struct {
	ID           string
	Repo         string
	StartTs      int64
	EndTs        int64
	CommandCount int64
	CreatedAt    int64
	UpdatedAt    int64
}

func NewSessionRepository(queries *sqlc.Queries) *SessionRepository {
	return &SessionRepository{queries: queries}
}

func (r *SessionRepository) Create(ctx context.Context, params CreateSessionParams) (Session, error) {
	return r.queries.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:           params.ID,
		Repo:         params.Repo,
		StartTs:      params.StartTs,
		EndTs:        params.EndTs,
		CommandCount: params.CommandCount,
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.UpdatedAt,
	})
}

func (r *SessionRepository) Get(ctx context.Context, id string) (Session, error) {
	return r.queries.GetSession(ctx, id)
}

func (r *SessionRepository) ListRecent(ctx context.Context, page Page) ([]Session, error) {
	page = normalizePage(page)
	return r.queries.ListRecentSessions(ctx, sqlc.ListRecentSessionsParams{
		Limit:  page.Limit,
		Offset: page.Offset,
	})
}

// SearchByCommand returns sessions containing a command matching the query
// (substring, newest first).
func (r *SessionRepository) SearchByCommand(ctx context.Context, query string, limit int64) ([]Session, error) {
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	return r.queries.SearchSessionsByCommand(ctx, sqlc.SearchSessionsByCommandParams{
		Pattern:  likePattern(query),
		LimitVal: limit,
	})
}

func (r *SessionRepository) TouchForCommand(ctx context.Context, id string, endTs int64, updatedAt int64) (Session, error) {
	return r.queries.TouchSessionForCommand(ctx, sqlc.TouchSessionForCommandParams{
		EndTs:     endTs,
		UpdatedAt: updatedAt,
		ID:        id,
	})
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteSession(ctx, id)
}
