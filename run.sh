#!/bin/bash
set -euo pipefail

PIPE=$(mktemp)
trap "rm \"$PIPE\"" EXIT

echo ": Building docker image..."
docker build . --iidfile "$PIPE"
echo ": Running built docker image $(cat $PIPE)"
docker run -it --rm "$(cat $PIPE)"
