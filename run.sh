#!/bin/bash
set -euo pipefail

PIPE=$(mktemp)
trap "rm \"$PIPE\"" EXIT

echo ": Building docker image..."
docker build . --iidfile "$PIPE"
echo ": Running built docker image $(cat $PIPE)"
docker run -p 8080:8080 -it --rm "$(cat $PIPE)"
