# Recall command capture

typeset -g RECALL_CMD=""
typeset -g RECALL_TS=0
typeset -g RECALL_CWD=""

preexec() {
    RECALL_CMD="$1"
    RECALL_TS=$(date +%s)
    RECALL_CWD="$PWD"
}

precmd() {
    # Ignore if no command was executed
    [[ -z "$RECALL_CMD" ]] && return

    local exit_code=$?
    local repo=""

    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        repo=$(basename "$(git rev-parse --show-toplevel)" 2>/dev/null)
    fi

    mkdir -p "$HOME/.local/share/recall"

    printf '%s\t%s\t%s\t%s\t%s\n' \
        "$RECALL_TS" \
        "$exit_code" \
        "$RECALL_CWD" \
        "$repo" \
        "$RECALL_CMD" \
        >> "$HOME/.local/share/recall/events.log"

    RECALL_CMD=""
}
