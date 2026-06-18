-- +goose Up
CREATE VIRTUAL TABLE memories_fts USING fts5(
    title,
    summary,
    content='memories',
    content_rowid='rowid'
);

-- +goose StatementBegin
CREATE TRIGGER memories_ai AFTER INSERT ON memories BEGIN
  INSERT INTO memories_fts(rowid, title, summary) VALUES (new.rowid, new.title, new.summary);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER memories_ad AFTER DELETE ON memories BEGIN
  INSERT INTO memories_fts(memories_fts, rowid, title, summary) VALUES('delete', old.rowid, old.title, old.summary);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER memories_au AFTER UPDATE ON memories BEGIN
  INSERT INTO memories_fts(memories_fts, rowid, title, summary) VALUES('delete', old.rowid, old.title, old.summary);
  INSERT INTO memories_fts(rowid, title, summary) VALUES (new.rowid, new.title, new.summary);
END;
-- +goose StatementEnd

-- Populate FTS index with existing data
INSERT INTO memories_fts(rowid, title, summary) SELECT rowid, title, summary FROM memories;

-- +goose Down
DROP TRIGGER IF EXISTS memories_ai;
DROP TRIGGER IF EXISTS memories_ad;
DROP TRIGGER IF EXISTS memories_au;
DROP TABLE IF EXISTS memories_fts;
