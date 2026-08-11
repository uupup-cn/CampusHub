#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
go test -count=1 -timeout 120s -run TestMarketplaceMyItems -v . 2>&1 | grep -E "FAIL|Error|expected|got|My items|items" | head -20
echo "=== done ==="
