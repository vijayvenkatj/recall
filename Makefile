migrate:
	goose sqlite3 replay.db up
generate:
	sqlc generate
build:
	go build ./...
