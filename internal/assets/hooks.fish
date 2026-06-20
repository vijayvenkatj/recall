# Recall Fish hook

function recall_preexec --on-event fish_preexec
    set -g RECALL_CMD $argv[1]
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
