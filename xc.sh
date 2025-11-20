#!/bin/sh

KWBUILDER_IMAGE="${KWBUILDER_IMAGE:-kwbuilder}"

if ! podman image inspect "${KWBUILDER_IMAGE}" >/dev/null 2>&1; then
  echo "Generating ${KWBUILDER_IMAGE}"
  podman build --platform linux/arm64 -t "${KWBUILDER_IMAGE}" -f builder/Containerfile builder
fi

podman run -it --rm --platform linux/arm64 \
  -v "$PWD":/src:z -w /src -v gocache:/cache -e GOMAXPROCS=10 \
  "${KWBUILDER_IMAGE}" "$@"
