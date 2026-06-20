# Recall

Recall is a privacy-first, local-only CLI tool that captures your terminal command history and groups them into structured, searchable memories. Built with Go and SQLite FTS5, Recall helps you remember exactly how you solved specific debugging issues, environment setups, or system bugs directly from your terminal.

---

## Features

- **Asynchronous Capture**: Lightweight shell integration hook records commands, working directories, Git repositories, and exit status to a local event log asynchronously without introducing latency to your shell.
- **Context-Aware Sessions**: Automatically groups commands into local developer sessions based on repository name, working directory, and inactive thresholds (default 30-minute idle window).
- **BM25 Search**: Matches terms dynamically against problems, resolutions, and even exact shell commands, sorted by SQLite FTS5 BM25 relevance scores.
- **Pluggable LLM Suggestions**: Leverages local Ollama instances or Google Gemini to analyze your command logs and pre-fill memory titles, problem statements, and fixes inside the wizard.
- **100% Privacy-First**: 100% local database and assets with no external analytics or telemetry.

---

## How It Works

Recall operates in three stages:

1. **Capture**: A lightweight shell hook logs execution details (timestamp, directory, exit code, Git repository, and command) to `~/.local/share/recall/events.log`.
2. **Consolidate**: Running `recall save` opens an interactive TUI to review recent sessions. If an LLM provider is configured, it auto-drafts suggestions for the problem and fix in the background.
3. **Retrieve**: Running `recall <query>` queries the SQLite FTS5 virtual table to find relevant memories and command histories inside an interactive viewer.

---

## Getting Started

### Prerequisites

Ensure you have the following installed:
- Go (1.26.3 or higher)
- Make
- sqlc

### 1. Build and Install

Clone the repository and run the install target:

```bash
git clone https://github.com/vijayvenkatj/recall.git
cd recall
make install
```

This compiles the binary and installs it to `~/.local/bin/recall`.

### 2. Configure Shell Hooks

Initialize configurations, databases, and hooks:

```bash
recall install
```

### 3. Add to Shell Profile

Source the generated hook file in your shell configuration.

- **Zsh (`~/.zshrc`)**:
  ```bash
  source ~/.config/recall/hooks.zsh
  ```
- **Bash (`~/.bashrc`)**:
  ```bash
  source ~/.config/recall/hooks.bash
  ```
- **Fish (`~/.config/fish/config.fish`)**:
  ```fish
  source ~/.config/recall/hooks.fish
  ```

Restart your terminal or reload your configuration (e.g., `source ~/.zshrc`).

---

## Usage

### Saving a Memory
Whenever you solve a bug or complete a task, run:
```bash
recall save
```
An interactive wizard will launch:
1. Select a recent session to review.
2. Review command logs.
3. Annotate or edit the generated suggestions for the problem and resolution.
4. Save the entry to SQLite.

### Searching Memories
Search through your database by typing keyword fragments:
```bash
recall docker container mapping
```
Running `recall` with no arguments will open the search viewer listing your 20 most recent memories:
```bash
recall
```
Inside the search viewer, press `/` to filter matching items. Press `Enter` to expand details and view the command logs that solved the problem.

---

## Configuration

Recall is configured via a YAML file located at `~/.config/recall/config.yaml`. 

Settings can be overridden by environment variables or a `.env` file:

| YAML Configuration Key | Environment Variable | Description | Default |
| :--- | :--- | :--- | :--- |
| `db_driver` | `DB_DRIVER` | Database engine driver | `sqlite` |
| `db_string` | `DB_STRING` | Absolute path to SQLite file | `~/.local/share/recall/recall.db` |
| `event_log_path` | `EVENT_LOG_PATH` | Path to raw shell events file | `~/.local/share/recall/events.log` |
| `log_level` | `LOG_LEVEL` | Application logging level | `info` |
| `llm_provider` | `LLM_PROVIDER` | LLM suggestions provider (`gemini` or `ollama`) | `""` (Disabled) |
| `llm_api_key` | `LLM_API_KEY` | API key (required for Gemini) | `""` |
| `llm_model` | `LLM_MODEL` | LLM model name (e.g. `gemini-2.5-flash` or `llama3`) | Provider default |
| `llm_endpoint` | `LLM_ENDPOINT` | Local service endpoint (for Ollama) | `http://localhost:11434` |

### LLM Configurations

#### Google Gemini API Setup
In `~/.config/recall/config.yaml`:
```yaml
llm_provider: "gemini"
llm_api_key: "YOUR_GEMINI_API_KEY"
llm_model: "gemini-2.5-flash"
```

#### Local Ollama Setup
Start your local Ollama server (`ollama run llama3`), then configure `~/.config/recall/config.yaml`:
```yaml
llm_provider: "ollama"
llm_endpoint: "http://localhost:11434"
llm_model: "llama3"
```

---

## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/my-feature`.
3. If database models or queries are altered, run:
   ```bash
   make generate
   ```
4. Push changes and submit a Pull Request.

---

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.
