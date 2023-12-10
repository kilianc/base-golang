#!/bin/sh

go version

# succeeds on linux/amd64
# go version go1.21.5 linux/arm64

# fails on darwin/arm64
# go version go1.21.5 darwin/arm64

for i in $(seq 1 10);
do
  echo "Run $i"
  go test -gcflags=all=-l
done
