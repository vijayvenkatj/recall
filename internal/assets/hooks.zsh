# Recall command capture

typeset -g RECALL_CMD=""
typeset -g RECALL_TS=0
typeset -g RECALL_CWD=""

# Commands matching this (case-insensitive) are never recorded, so secrets
# never touch the log. Set RECALL_IGNORE to a regex for your own excludes.
typeset -g RECALL_SECRET_RE='(password|passwd|secret|token|api[_-]?key|access[_-]?key|bearer|authorization|credential|private[_-]?key)'

recall_should_ignore() {
    local cmd="$1"
    # Leading space = "don't record this" (like HISTCONTROL=ignorespace).
    [[ "$cmd" == ' '* ]] && return 0
    [[ "${cmd:l}" =~ $RECALL_SECRET_RE ]] && return 0
    [[ -n "$RECALL_IGNORE" && "$cmd" =~ $RECALL_IGNORE ]] && return 0
    return 1
}

preexec() {
    if recall_should_ignore "$1"; then
        RECALL_CMD=""
        return
    fi
    RECALL_CMD="$1"
    RECALL_TS=$(date +%s)
    RECALL_CWD="$PWD"
}

precmd() {
    local exit_code=$?
    # Ignore if no command was executed
    [[ -z "$RECALL_CMD" ]] && return

    local repo=""
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        repo=$(basename "$(git rev-parse --show-toplevel)" 2>/dev/null)
    fi

    local log_path="${EVENT_LOG_PATH:-$HOME/.local/share/recall/events.log}"
    mkdir -p "$(dirname "$log_path")"

    # Escape tabs and newlines in the command
    local cmd_escaped="${RECALL_CMD//[$'\t']/    }"
    cmd_escaped="${cmd_escaped//[$'\n']/\\n}"

    printf '%s\t%s\t%s\t%s\t%s\n' \
        "$RECALL_TS" \
        "$exit_code" \
        "$RECALL_CWD" \
        "$repo" \
        "$cmd_escaped" \
        >> "$log_path"

    RECALL_CMD=""
}
