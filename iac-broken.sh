#!/bin/bash
# Enable the secret task 8 and give everyone a negative score.

TASK_NAME=2a-infra
PART_NAME=fix-prod

. ../setup/config.sh

if [[ -z "${HOST_PROJECT}" ]]; then
  echo "Error - expected HOST_PROJECT variable"
  exit 1
fi

SCORE_BUCKET_NAME="${HOST_PROJECT}-score-assets"

# Deploy the secret task 8
sed -i 's/hidden: true/hidden: false/g; s/lb_hidden: true/lb_hidden: false/g; s/enabled: false/enabled: true/g' "${TASK_NAME}/task.yaml"

go run ../utils/task-tool/main.go --base-folder "${PWD}" --bucket "${SCORE_BUCKET_NAME}" --host-pid="${HOST_PROJECT}" --upload=false --upload-images=false

# Give everyone negative points!
go run ../utils/points-tool/main.go --all-teams --host-pid "${HOST_PROJECT}" --part-id "${PART_NAME}" --task-id "${TASK_NAME}" --points "-4000"

