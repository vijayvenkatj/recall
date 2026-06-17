package repository

import (
	"context"

	sqlc "github.com/vijayvenkatj/recall/internal/db/sqlc"
)

type Metadata = sqlc.Metadata

type MetadataRepository struct {
	queries *sqlc.Queries
}

func NewMetadataRepository(queries *sqlc.Queries) *MetadataRepository {
	return &MetadataRepository{
		queries: queries,
	}
}

func (r *MetadataRepository) Get(ctx context.Context, key string) (Metadata, error) {
	return r.queries.GetMetadata(ctx, key)
}

func (r *MetadataRepository) Set(ctx context.Context, key string, value string) (Metadata, error) {
	return r.queries.UpsertMetadata(ctx,
		sqlc.UpsertMetadataParams{
			Key:   key,
			Value: value,
		},
	)
}
