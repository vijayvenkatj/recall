-- +goose Up
CREATE TABLE metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- +goose down
DROP TABLE metadata;