package repository

import "database/sql"

const (
	defaultLimit = int64(50)
	maxLimit     = int64(500)
)

type Page struct {
	Limit  int64
	Offset int64
}

func normalizePage(page Page) Page {
	if page.Limit <= 0 {
		page.Limit = defaultLimit
	}
	if page.Limit > maxLimit {
		page.Limit = maxLimit
	}
	if page.Offset < 0 {
		page.Offset = 0
	}
	return page
}

func nullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func nullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}
