#!/bin/bash
. ../common.sh

# Lock all tasks
echo Locking all tasks
tech-day --action lock

# Run Act 3 Ending
tech-day --action act3-end

tech-day --action lock
