# Recall Bash hook

RECALL_CMD=""
RECALL_TS=0
RECALL_CWD=""

recall_preexec() {
    # Ignore completion commands and subshells
    [ -n "$COMP_LINE" ] && return
    [ "$BASH_COMMAND" = "$PROMPT_COMMAND" ] && return

    RECALL_CMD="$BASH_COMMAND"
    RECALL_TS=$(date +%s)
    RECALL_CWD="$PWD"
}

recall_precmd() {
    local exit_code=$?
    [ -z "$RECALL_CMD" ] && return

    local repo=""
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        repo=$(basename "$(git rev-parse --show-toplevel)" 2>/dev/null)
    fi

    local log_path="${EVENT_LOG_PATH:-$HOME/.local/share/recall/events.log}"
    mkdir -p "$(dirname "$log_path")"

    # Escape tabs and newlines
    local cmd_escaped="${RECALL_CMD//$'\t'/    }"
    cmd_escaped="${cmd_escaped//$'\n'/\\n}"

    printf '%s\t%s\t%s\t%s\t%s\n' \
        "$RECALL_TS" \
        "$exit_code" \
        "$RECALL_CWD" \
        "$repo" \
        "$cmd_escaped" \
        >> "$log_path"

    RECALL_CMD=""
}

# Register traps and prompt commands
trap 'recall_preexec' DEBUG
if [[ -z "$PROMPT_COMMAND" ]]; then
    PROMPT_COMMAND="recall_precmd"
elif [[ "$PROMPT_COMMAND" != *"recall_precmd"* ]]; then
    PROMPT_COMMAND="recall_precmd; $PROMPT_COMMAND"
fi
