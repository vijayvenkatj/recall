# Recall Fish hook

# Commands matching this (case-insensitive) have their secret VALUES redacted
# before being recorded, so secrets never touch the log.
set -g RECALL_SECRET_RE '(password|passwd|secret|token|api[_-]?key|access[_-]?key|bearer|credential|passphrase|private[_-]?key)'

# Drop a command entirely: leading space ("don't record this") or the user's
# RECALL_IGNORE regex.
function recall_should_drop
    set -l cmd $argv[1]
    if string match -rq '^\s' -- "$cmd"
        return 0
    end
    if set -q RECALL_IGNORE; and test -n "$RECALL_IGNORE"
        if string match -rq -- "$RECALL_IGNORE" "$cmd"
            return 0
        end
    end
    return 1
end

# Mask secret values (keys, tokens, passwords) while keeping the command shape.
function recall_redact
    printf '%s' "$argv[1]" | sed -E \
        -e 's/([A-Za-z0-9_]*(TOKEN|SECRET|PASSWORD|PASSWD|APIKEY|API_KEY|ACCESS_KEY|KEY|CREDENTIAL|PASSPHRASE|token|secret|password|passwd|apikey|api_key|access_key|key|credential|passphrase)=)[^[:space:]]*/\1<redacted>/g' \
        -e 's/(--(password|passwd|token|secret|api-key|apikey|access-key|credential|passphrase)[= ])[^[:space:]]*/\1<redacted>/g' \
        -e 's/((Bearer|bearer) )[^[:space:]]*/\1<redacted>/g'
end

function recall_preexec --on-event fish_preexec
    if recall_should_drop "$argv[1]"
        set -g RECALL_CMD ""
        return
    end
    if string match -rqi -- "$RECALL_SECRET_RE" "$argv[1]"
        set -g RECALL_CMD (recall_redact "$argv[1]" | string collect)
    else
        set -g RECALL_CMD $argv[1]
    end
    set -g RECALL_TS (date +%s)
    set -g RECALL_CWD $PWD
end

function recall_precmd --on-event fish_postexec
    set -l exit_code $status
    if not set -q RECALL_CMD; or test -z "$RECALL_CMD"
        return
    end

    set -l repo ""
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1
        set repo (basename (git rev-parse --show-toplevel 2>/dev/null) 2>/dev/null)
    end

    set -l log_path $EVENT_LOG_PATH
    if test -z "$log_path"
        set log_path "$HOME/.local/share/recall/events.log"
    end
    mkdir -p (dirname "$log_path")

    # Escape tabs and newlines
    set -l cmd_escaped (string replace -a \t "    " -- "$RECALL_CMD")
    set cmd_escaped (string replace -a \n "\\n" -- "$cmd_escaped")

    printf "%s\t%s\t%s\t%s\t%s\n" \
        "$RECALL_TS" \
        "$exit_code" \
        "$RECALL_CWD" \
        "$repo" \
        "$cmd_escaped" \
        >> "$log_path"

    set -e RECALL_CMD
end
