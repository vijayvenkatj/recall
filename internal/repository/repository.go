package repository

import (
	"context"
	"database/sql"

	sqlc "github.com/vijayvenkatj/recall/internal/db/sqlc"
)

type Store struct {
	DB      *sql.DB
	queries *sqlc.Queries

	Commands *CommandRepository
	Sessions *SessionRepository
	Memories *MemoryRepository
	Metadata *MetadataRepository
}

func New(db *sql.DB) *Store {
	queries := sqlc.New(db)
	return newStore(db, queries)
}

func newStore(db *sql.DB, queries *sqlc.Queries) *Store {
	return &Store{
		DB:       db,
		queries:  queries,
		Commands: NewCommandRepository(db, queries),
		Sessions: NewSessionRepository(queries),
		Memories: NewMemoryRepository(queries),
		Metadata: NewMetadataRepository(queries),
	}
}

func (s *Store) Queries() *sqlc.Queries {
	return s.queries
}

func (s *Store) InTx(ctx context.Context, fn func(*Store) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := newStore(s.DB, s.queries.WithTx(tx))
	if err := fn(txStore); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}
