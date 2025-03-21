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

deploy_direct() {
  if ! command -v go >/dev/null 2>&1; then
    return 2
  fi
  if ! go run utils/task-tool/main.go --base-folder "${PWD}" --bucket "${SCORE_BUCKET_NAME}" --host-pid="${HOST_PROJECT}"; then
    return 2
  fi
}

deploy_cloud_build() {

  if [[ -z "${LOCATION}" ]]; then
    error "Expected LOCATION environment variable to be set"
  fi

  service_account="projects/${HOST_PROJECT}/serviceAccounts/task-builder@${HOST_PROJECT}.iam.gserviceaccount.com"

  check_cmd tar

  info "Submitting Tasks"

  # Create archive for build
  tar -czvf /tmp/tasks.tar.gz -C . . -C utils task-tool >/dev/null

  # Copy archive to Cloud Storage
  gcloud storage cp "/tmp/tasks.tar.gz" "gs://${ASSET_BUCKET_NAME}/tasks.tar.gz" >/dev/null

  # Run Cloud Build Job
  gcloud builds submit "gs://${ASSET_BUCKET_NAME}/tasks.tar.gz" --config cloudbuild.yaml --substitutions="_BUCKET_NAME=${SCORE_BUCKET_NAME},_HOST_PID=${HOST_PROJECT},_ROOT_DIR=../,_TASK_DIR=task-tool" --service-account="${service_account}" --project="${HOST_PROJECT}" --region "${LOCATION}" >/dev/null
}

main() {
  if [[ -z "${HOST_PROJECT}" ]]; then
    if [[ ! -f ../setup/config.sh ]]; then
      error Cannot find config. Ensure HOST_PROJECT is set or ../setup/config.sh exists
    fi
    . ../setup/config.sh
  fi

  ASSET_BUCKET_NAME="${HOST_PROJECT}-build-assets"
  SCORE_BUCKET_NAME="${HOST_PROJECT}-score-assets"

  # Attempt to run directly from here, otherwise resort to running in Cloud Build
  if ! deploy_direct; then
    deploy_cloud_build
  fi
  echo Task deployment complete.
}

main
