#!/bin/bash

export PROJECT_ID=q2-25-tech-day-host

basepath=../utils

if [[ ! -f "${basepath}/tech-day/main.go" ]]; then
  basepath=../../utils
  if [[ ! -f "${basepath}/tech-day/main.go" ]]; then
    echo Cannot determine base path
    exit 1
  fi
fi

tech-day() {

  go run "${basepath}"/tech-day/*.go "${@}"
}

asset-score() {
  go run "${basepath}"/asset-rank/*.go
}
