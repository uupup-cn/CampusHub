#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
go test -count=1 -timeout 120s -run "TestOAuth" -v . 2>&1 | grep -E "FAIL|PASS|Error|expected|missing|got|confirm" | head -25
echo "=== done ==="
