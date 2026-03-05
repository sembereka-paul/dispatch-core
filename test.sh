#!/usr/bin/env bash
set -euo pipefail

setup=$1
waitfor=$2

set -m
eval "$setup" &
bg=$!

sleep 30

eval "$waitfor"
# Make sure the background process is leader of its own group: use its PID as -pgid
# (pkill with negative PID targets the process group id)
pgid=$(ps -o pgid= "$bg" | tr -d ' ')
echo "bg pid=$bg pgid=$pgid"

trap 'rc=$?; echo "Cleaning group $pgid"; kill -TERM -"$pgid" 2>/dev/null || true; sleep 2; kill -KILL -"$pgid" 2>/dev/null || true; exit "$rc"' EXIT TERM INT

# On normal completion, kill the process group
# -KILL does not cause error exit code
kill -KILL -"$pgid" 2>/dev/null || true

exit 0
