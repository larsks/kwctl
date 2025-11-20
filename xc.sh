#!/bin/sh

KWBUILDER_IMAGE="${KWBUILDER_IMAGE:-kwbuilder}"
KWBUILD_PLATFORM="${KWBUILDER_PLATFORM:-linux/arm64}"

if ! podman image inspect "${KWBUILDER_IMAGE}" >/dev/null 2>&1; then
  echo "Generating ${KWBUILDER_IMAGE}"
  podman build --platform "${KWBUILD_PLATFORM}" -t "${KWBUILDER_IMAGE}" -f builder/Containerfile builder
fi

podman run -it --rm --platform "${KWBUILD_PLATFORM}" \
  -v "$PWD":/src:z -w /src -v gocache:/cache -e GOMAXPROCS=10 \
  "${KWBUILDER_IMAGE}" "$@"
