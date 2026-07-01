package repository

import (
	"database/sql"
	"strings"
)

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

// likePattern turns a free-text query into a SQL LIKE pattern (used with
// ESCAPE '\'). Each word must appear, in order, so "docker port" matches
// "docker ... port". LIKE wildcards in the input are escaped — critical for
// commands, which are full of '_'. Empty input yields "%" (match anything).
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func likePattern(query string) string {
	words := strings.Fields(query)
	for i, w := range words {
		words[i] = likeEscaper.Replace(w)
	}
	return "%" + strings.Join(words, "%") + "%"
}
