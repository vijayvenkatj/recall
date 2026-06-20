-- +goose Up
CREATE VIRTUAL TABLE memories_fts USING fts5(
    title,
    summary,
    commands
);

-- +goose StatementBegin
CREATE TRIGGER memories_ai AFTER INSERT ON memories BEGIN
  INSERT INTO memories_fts(rowid, title, summary, commands)
  VALUES (
    new.rowid,
    new.title,
    new.summary,
    (SELECT group_concat(command, ' ') FROM commands WHERE session_id = new.session_id)
  );
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER memories_ad AFTER DELETE ON memories BEGIN
  DELETE FROM memories_fts WHERE rowid = old.rowid;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER memories_au AFTER UPDATE ON memories BEGIN
  UPDATE memories_fts
  SET 
    title = new.title,
    summary = new.summary,
    commands = (SELECT group_concat(command, ' ') FROM commands WHERE session_id = new.session_id)
  WHERE rowid = new.rowid;
END;
-- +goose StatementEnd

-- Populate FTS index with existing data
INSERT INTO memories_fts(rowid, title, summary, commands)
SELECT 
  rowid, 
  title, 
  summary,
  (SELECT group_concat(command, ' ') FROM commands WHERE session_id = memories.session_id)
FROM memories;

-- +goose Down
DROP TRIGGER IF EXISTS memories_ai;
DROP TRIGGER IF EXISTS memories_ad;
DROP TRIGGER IF EXISTS memories_au;
DROP TABLE IF EXISTS memories_fts;
