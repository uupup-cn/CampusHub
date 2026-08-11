#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
go test -count=1 -timeout 120s -run TestResponseFormatConsistency -v . 2>&1 | grep -E "FAIL|PASS|Error|expected|missing|snake|got" | head -20
echo "=== done ==="
