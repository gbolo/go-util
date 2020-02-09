#!/usr/bin/env bash
set -e

echo "starting hq container..."
docker run -it --rm --name "gocm-hq" \
 --net host \
 -v $(pwd)/testdata:/testdata:ro \
 gbolo/gocm /hq
