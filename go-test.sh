#!/bin/sh

# ensure xgo installed
go install github.com/xhd2015/xgo/cmd/xgo@latest

xgo version

# succeeds on linux/amd64
# go version go1.21.5 linux/arm64

# fails on darwin/arm64
# go version go1.21.5 darwin/arm64

for i in $(seq 1 100);
do
  echo "Run $i"
  xgo test -v ./
done
