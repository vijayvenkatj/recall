package repository

import (
	"context"
	"database/sql"

	sqlc "github.com/vijayvenkatj/recall/internal/db/sqlc"
)

type Command = sqlc.Command

type CommandRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

type CreateCommandParams struct {
	ID        string
	SessionID string
	Timestamp int64
	Command   string
	CWD       *string
	Repo      *string
	ExitCode  *int64
	CreatedAt int64
}

type AppendCommandResult struct {
	Command Command
	Session Session
}

func NewCommandRepository(db *sql.DB, queries *sqlc.Queries) *CommandRepository {
	return &CommandRepository{db: db, queries: queries}
}

func (r *CommandRepository) Create(ctx context.Context, params CreateCommandParams) (Command, error) {
	return r.queries.CreateCommand(ctx, createCommandParams(params))
}

func (r *CommandRepository) Append(ctx context.Context, params CreateCommandParams) (AppendCommandResult, error) {
	var result AppendCommandResult

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}

	queries := r.queries.WithTx(tx)

	command, err := queries.CreateCommand(ctx, createCommandParams(params))
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return result, rollbackErr
		}
		return result, err
	}

	session, err := queries.TouchSessionForCommand(ctx, sqlc.TouchSessionForCommandParams{
		EndTs:     params.Timestamp,
		UpdatedAt: params.CreatedAt,
		ID:        params.SessionID,
	})
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return result, rollbackErr
		}
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	result.Command = command
	result.Session = session
	return result, nil
}

func (r *CommandRepository) Get(ctx context.Context, id string) (Command, error) {
	return r.queries.GetCommand(ctx, id)
}

func (r *CommandRepository) ListBySession(ctx context.Context, sessionID string, page Page) ([]Command, error) {
	page = normalizePage(page)
	return r.queries.ListCommandsBySession(ctx, sqlc.ListCommandsBySessionParams{
		SessionID: sessionID,
		Limit:     page.Limit,
		Offset:    page.Offset,
	})
}

func (r *CommandRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteCommand(ctx, id)
}

func createCommandParams(params CreateCommandParams) sqlc.CreateCommandParams {
	return sqlc.CreateCommandParams{
		ID:        params.ID,
		SessionID: params.SessionID,
		Timestamp: params.Timestamp,
		Command:   params.Command,
		Cwd:       nullString(params.CWD),
		Repo:      nullString(params.Repo),
		ExitCode:  nullInt64(params.ExitCode),
		CreatedAt: params.CreatedAt,
	}
}
