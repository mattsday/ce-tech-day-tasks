#!/bin/bash
. ../common.sh

# Lock all tasks
echo Locking all tasks
tech-day --action lock

# Run asset scoring
export TASK_ID=act2-task2
export PART_ID=audition
asset-score

# Now calculate Dragon's Den places
go run ../../utils/dragons-finals/main.go

# Run Act 2 Ending
tech-day --action act2-end

