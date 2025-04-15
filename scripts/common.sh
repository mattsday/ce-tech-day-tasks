#!/bin/bash

export PROJECT_ID=q2-25-tech-day-host

tech-day() {
  go run ../utils/tech-day/*.go "${@}"
}

asset-score() {
  echo Asset Score Tool
}
