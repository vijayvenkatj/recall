# Recall command capture

typeset -g RECALL_CMD=""
typeset -g RECALL_TS=0
typeset -g RECALL_CWD=""

# Commands matching this (case-insensitive) have their secret VALUES redacted
# before being recorded, so secrets never touch the log.
typeset -g RECALL_SECRET_RE='(password|passwd|secret|token|api[_-]?key|access[_-]?key|bearer|credential|passphrase|private[_-]?key)'

# Drop a command entirely: leading space ("don't record this", like
# HISTCONTROL=ignorespace) or the user's RECALL_IGNORE regex.
recall_should_drop() {
    local cmd="$1"
    [[ "$cmd" == ' '* ]] && return 0
    [[ -n "$RECALL_IGNORE" && "$cmd" =~ $RECALL_IGNORE ]] && return 0
    return 1
}

# Mask secret values (keys, tokens, passwords) while keeping the command shape.
recall_redact() {
    printf '%s' "$1" | sed -E \
        -e 's/([A-Za-z0-9_]*(TOKEN|SECRET|PASSWORD|PASSWD|APIKEY|API_KEY|ACCESS_KEY|KEY|CREDENTIAL|PASSPHRASE|token|secret|password|passwd|apikey|api_key|access_key|key|credential|passphrase)=)[^[:space:]]*/\1<redacted>/g' \
        -e 's/(--(password|passwd|token|secret|api-key|apikey|access-key|credential|passphrase)[= ])[^[:space:]]*/\1<redacted>/g' \
        -e 's/((Bearer|bearer) )[^[:space:]]*/\1<redacted>/g'
}

preexec() {
    if recall_should_drop "$1"; then
        RECALL_CMD=""
        return
    fi
    if [[ "${1:l}" =~ $RECALL_SECRET_RE ]]; then
        RECALL_CMD="$(recall_redact "$1")"
    else
        RECALL_CMD="$1"
    fi
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
