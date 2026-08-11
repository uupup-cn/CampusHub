#!/bin/bash
export PATH="/home/pc/go-install/go/bin:$PATH"
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan/chb-backend
echo "Go version: $(go version)"
echo "Working dir: $(pwd)"
echo "=== Running tests ==="
GOPROXY=https://goproxy.cn,direct go test -v -run . -timeout 180s ./... 2>&1
echo "=== Exit code: $? ==="
