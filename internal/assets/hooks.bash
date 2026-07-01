# Recall Bash hook

RECALL_CMD=""
RECALL_TS=0
RECALL_CWD=""

# Commands matching this (case-insensitive) have their secret VALUES redacted
# before being recorded, so secrets never touch the log.
RECALL_SECRET_RE='(password|passwd|secret|token|api[_-]?key|access[_-]?key|bearer|credential|passphrase|private[_-]?key)'

# Drop a command entirely: leading space ("don't record this"; best effort in
# Bash since BASH_COMMAND usually strips it) or the user's RECALL_IGNORE regex.
recall_should_drop() {
    local cmd="$1"
    case "$cmd" in ' '*) return 0 ;; esac
    [ -n "$RECALL_IGNORE" ] && [[ "$cmd" =~ $RECALL_IGNORE ]] && return 0
    return 1
}

recall_is_secret() {
    local rc=1 restore
    shopt -q nocasematch && restore="-s" || restore="-u"
    shopt -s nocasematch
    [[ "$1" =~ $RECALL_SECRET_RE ]] && rc=0
    shopt $restore nocasematch
    return $rc
}

# Mask secret values (keys, tokens, passwords) while keeping the command shape.
recall_redact() {
    printf '%s' "$1" | sed -E \
        -e 's/([A-Za-z0-9_]*(TOKEN|SECRET|PASSWORD|PASSWD|APIKEY|API_KEY|ACCESS_KEY|KEY|CREDENTIAL|PASSPHRASE|token|secret|password|passwd|apikey|api_key|access_key|key|credential|passphrase)=)[^[:space:]]*/\1<redacted>/g' \
        -e 's/(--(password|passwd|token|secret|api-key|apikey|access-key|credential|passphrase)[= ])[^[:space:]]*/\1<redacted>/g' \
        -e 's/((Bearer|bearer) )[^[:space:]]*/\1<redacted>/g'
}

recall_preexec() {
    # Ignore completion commands and subshells
    [ -n "$COMP_LINE" ] && return
    [ "$BASH_COMMAND" = "$PROMPT_COMMAND" ] && return

    if recall_should_drop "$BASH_COMMAND"; then
        RECALL_CMD=""
        return
    fi

    if recall_is_secret "$BASH_COMMAND"; then
        RECALL_CMD="$(recall_redact "$BASH_COMMAND")"
    else
        RECALL_CMD="$BASH_COMMAND"
    fi
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
