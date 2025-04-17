#!/bin/bash
. ../common.sh

# Lock all tasks
echo Locking all tasks
tech-day --action lock

# Run asset scoring
export TASK_ID=act1-task4
export PART_ID=part1
asset-score

# Run Act 1 Ending
tech-day --action act1-end

