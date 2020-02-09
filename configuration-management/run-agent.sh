#!/usr/bin/env bash
if [ $# -ne 1 ]
  then
    echo "a local port MUST be specified as the first argument (example: 18001)"
    exit 1
fi

AGENT_NAME="agent${1}"

echo "removing agent container: ${AGENT_NAME}"
docker rm -f ${AGENT_NAME}

set -e
echo "starting agent container: ${AGENT_NAME}"
docker run -d --rm --name "${AGENT_NAME}" \
 -p 127.0.0.1:${1}:8080 \
 gbolo/gocm

echo "logs..."
docker logs ${AGENT_NAME}

echo
echo "started ${AGENT_NAME}: url -> http://127.0.0.1:${1}"
