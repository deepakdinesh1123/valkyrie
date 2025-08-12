#!/usr/bin/env sh

# Get the parent PID of exec.sh
parent_id=$(pgrep -f exec.sh)

echo "Parent PID(s): $parent_id"

if [ -z "$parent_id" ]; then
    echo "No process named exec.sh found."
else
    for pid in $parent_id; do
        # Get child PIDs of the parent PID
        child_pids=$(pgrep -P "$pid")
        
        for child_pid in $child_pids; do
            if kill -2 "$child_pid" 2>/dev/null; then
                echo "Sent SIGINT to process $child_pid"
            else
                echo "Failed to send SIGINT to process $child_pid (may not exist)"
            fi
        done
    done
fi
