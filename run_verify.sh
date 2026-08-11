#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
echo "=== go vet ==="
go vet ./... 2>&1 | head -30 || true
echo "=== go build ==="
go build ./... 2>&1 | head -30
echo "=== done ==="
