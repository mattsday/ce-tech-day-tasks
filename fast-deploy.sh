#!/bin/bash
# Deploy task config to Cloud.
#
# shellcheck disable=SC1091

# Print an error to stderr and stop
error() {
  echo >&2 '[Error]' "$@"
  exit 1
}

# Print information to stdout
info() {
  echo '[Info]' "$@"
}

# Checks if a command exists and halts with a fatal error if it does not
# Usage:
# check_cmd <cmd_name>
# e.g.
# check_cmd jq
check_cmd() {
  if ! command -v "$@" >/dev/null 2>&1; then
    error Command "$@" not found - please install it
  fi
}
main() {
  if [[ -z "${HOST_PROJECT}" ]]; then
    if [[ ! -f ../setup/config.sh ]]; then
      error Cannot find config. Ensure HOST_PROJECT is set or ../setup/config.sh exists
    fi
    . ../setup/config.sh
  fi

  SCORE_BUCKET_NAME="${HOST_PROJECT}-score-assets"
  
  go run ../utils/task-tool/main.go --base-folder "${PWD}" --bucket "${SCORE_BUCKET_NAME}" --host-pid="${HOST_PROJECT}" --upload-images=false

  echo Task deployment complete.
}

main
