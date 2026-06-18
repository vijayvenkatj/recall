# 🧠 Recall

[![Go Version](https://img.shields.io/github/go-mod/go-version/vijayvenkatj/recall?color=00ADD8&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)](#)

Recall is a professional, privacy-first developer productivity CLI tool designed to capture your terminal history and structure them into searchable, contextual memories. Built with Go and SQLite FTS5, Recall helps you remember how you solved specific debugging issues, environment setups, and system bugs directly from your terminal.

---

## ⚡️ Key Features

- **⚡ Lightweight Shell Integration** — Leverages asynchronous zsh hooks to log commands and exit status to a local file without introducing latency into your shell execution.
- **📂 Context-Aware Session Grouping** — Automatically clusters commands into developer sessions based on working directory, active Git repository, and idle times (30-minute default threshold).
- **📟 Interactive Terminal TUI** — Smooth, terminal-based user interfaces for both saving memories (defining problem & resolution) and searching through them.
- **🔍 Full-Text Search (FTS5)** — Lightning-fast SQLite FTS5 indexing allowing you to search through your shell commands, repositories, problem descriptions, and fixes.
- **🔒 Local-First & Privacy-Focused** — 100% local database with zero external analytics or API calls. Your terminal activity and workspace histories never leave your machine.

---

## 🏗 How It Works

Recall splits the workflow into three primary stages to seamlessly capture and retrieve context:

1. **Capture**: A lightweight shell hook logs execution details (timestamp, working directory, git repository, exit status, and command) to a local raw events log (`~/.local/share/recall/events.log`) without introducing latency to your shell.
2. **Consolidate**: Running `recall sync` parses the raw events and groups them into local dev sessions (defined as commands run within a 30-minute window in the same repository). Running `recall save` opens an interactive terminal TUI to annotate these sessions with a **Problem** and a **Fix**, saving them to SQLite.
3. **Retrieve**: Running `recall <query>` queries the SQLite FTS5 index to search your memories, rendering matches inside an interactive details browser.

---

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed:
- **Go** (version 1.26.3 or higher)
- **Make**
- **sqlc** (for database code generation)

### 1. Build and Install

Clone the repository and run the Makefile install command:

```bash
git clone https://github.com/vijayvenkatj/recall.git
cd recall
make install
```

This compiles the binary, moves it to `~/.local/bin/recall`, and prompts you to complete the setup.

### 2. Configure Shell Hooks & Database

Run the installation command to initialize directories, configure database migrations, and export shell hooks:

```bash
recall install
```

### 3. Add to your Shell Profile

Add the following line to your `~/.zshrc` configuration to enable command capturing:

```bash
source ~/.config/recall/hooks.zsh
```

Then reload your shell:
```bash
source ~/.zshrc
```

---

## 📖 Usage

### Saving a Memory
When you finish a task, solve a bug, or execute a complex pipeline, run:
```bash
recall save
```
An interactive TUI will open:
1. **Select Session**: Choose from your most recent dev sessions (shows repository name, timestamp, and command count).
2. **Review Commands**: Scroll through the list of commands executed during that session to verify.
3. **Problem Details**: Input a description of the problem/error you faced.
4. **Fix Details**: Input how you resolved it.
5. **Save**: The details are saved as a searchable memory.

### Searching Memories
Search through your history with:
```bash
recall "my search term"
```
Or simply run `recall` with keywords:
```bash
recall docker container debug
```
This triggers an interactive search interface showing matching memories. Press `Enter` on any memory to view its details (Title, Problem, Fix) along with the exact command history that led to the solution.

---

## 🛠 Project Structure

- `cmd/` — CLI subcommands (`root`, `install`, `sync`, `save`) utilizing the Cobra CLI framework.
- `internal/app/` — Core application logic including interactive terminal TUI components (`memory.go`, `search.go`), database synchronization (`sync.go`), and installation helpers (`install.go`).
- `internal/repository/` — SQLite abstraction layer created using sqlc.
- `internal/db/migrations/` — Database version control schema using Goose migrations (uses modernc.org/sqlite driver).
- `internal/assets/` — Embedded templates (e.g. Zsh shell hook script).

---

## ⚙️ Configuration

Recall is configured out of the box with zero-setup defaults. If customization is required, the following environment variables can be set in your `.env` or shell configuration:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DB_DRIVER` | Database engine driver | `sqlite` |
| `DB_STRING` | Absolute path to SQLite file | `~/.local/share/recall/recall.db` |
| `EVENT_LOG_PATH` | Path to raw shell events file | `~/.local/share/recall/events.log` |
| `LOG_LEVEL` | Application logger severity level | `info` |

---

## 🤝 Contributing

We welcome contributions of any size! Here is how you can help:
1. Fork the repository.
2. Create a new branch: `git checkout -b feature/my-new-feature`.
3. If database models or queries are altered, modify files in `internal/db/queries/` or `internal/db/migrations/`, then run:
   ```bash
   make generate
   ```
4. Commit your changes and push them to your fork.
5. Open a Pull Request!

---

## 📄 License

Distributed under the MIT License. See [LICENSE](LICENSE) for more details.
