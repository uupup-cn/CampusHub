#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
go test -count=1 -timeout 180s -v ./... 2>&1 | grep -E "^(=== RUN|--- PASS|--- FAIL|PASS|FAIL|ok)" | head -80
echo "=== done ==="
