#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
echo "=== port check ==="
ss -tlnp 2>/dev/null | grep 5433 || true
echo "=== integration tests ==="
go test -count=1 -timeout 180s ./... 2>&1 | tail -40
echo "=== done ==="
