BINARY_NAME=recall
DATA_DIR=$(HOME)/.local/share/recall
DB_PATH=$(DATA_DIR)/recall.db

.PHONY: all build install clean migrate generate

all: build

build: generate
	go build -o build/$(BINARY_NAME) main.go

install: build
	mkdir -p $(HOME)/.local/bin
	cp build/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Built and moved binary to $(HOME)/.local/bin/$(BINARY_NAME)"
	@echo "\nRun 'recall install' to complete the setup (migrations & shell hooks)."

migrate:
	@echo "Running migrations via recall CLI..."
	go run main.go install

generate:
	sqlc generate

clean:
	rm -rf build/
