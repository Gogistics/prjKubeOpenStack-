#!/bin/bash
PRJ_DIR=$PWD

finish() {
  local existcode=$?
  cd $PRJ_DIR
  exit $existcode
}

trap "finish" INT TERM

function spin_up_web_app() {
  local dockerfile_path="$1"
  local port="$2"

  docker build -t web-app -f $dockerfile_path . &&
  docker run \
    --name web-app \
    --log-driver json-file \
    --log-opt mode=non-blocking \
    --log-opt max-buffer-size=4m \
    --log-opt max-size=50m \
    --log-opt max-file=5 \
    -p $port:$port \
    -d web-app
}

function spin_up_web_apis_config() {
  docker build -t web-config -f dockerfiles/Dockerfile.apis_config . &&
  docker run \
    --name web-config \
    --log-driver json-file \
    --log-opt mode=non-blocking \
    --log-opt max-buffer-size=4m \
    --log-opt max-size=50m \
    --log-opt max-file=5 \
    -p 3001:3001 \
    -d web-config
}

function spin_up_web_apis_state() {
  docker build -t web-state -f dockerfiles/Dockerfile.apis_state . &&
  docker run \
    --name web-state \
    --log-driver json-file \
    --log-opt mode=non-blocking \
    --log-opt max-buffer-size=4m \
    --log-opt max-size=50m \
    --log-opt max-file=5 \
    -p 3002:3002 \
    -d web-state
}

spin_up_web_app "dockerfiles/Dockerfile.web_app" "3000" || exit 1
spin_up_web_apis_config || exit 1
spin_up_web_apis_state || exit 1
