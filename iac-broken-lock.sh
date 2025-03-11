#!/bin/bash
# Enable the secret task 8 and give everyone a negative score.

TASK_NAME=2a-infra

. ../setup/config.sh

if [[ -z "${HOST_PROJECT}" ]]; then
  echo "Error - expected HOST_PROJECT variable"
  exit 1
fi

SCORE_BUCKET_NAME="${HOST_PROJECT}-score-assets"

# Deploy the secret task 8
sed -i 's/enabled: true/enabled: false/g' "${TASK_NAME}"/task.yaml

go run ../utils/task-tool/main.go --base-folder "${PWD}" --bucket "${SCORE_BUCKET_NAME}" --host-pid="${HOST_PROJECT}" --upload=false
