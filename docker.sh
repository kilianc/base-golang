#!/bin/bash

set -e

docker build -t base-go:local .
docker run base-go:local ./go-test.sh
