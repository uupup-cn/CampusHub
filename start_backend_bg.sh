#!/bin/bash
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
pkill -f "cmd/server" 2>/dev/null || true
pkill -f "exe/server" 2>/dev/null || true
pkill -f "go-build.*/server" 2>/dev/null || true
fuser -k 9090/tcp 2>/dev/null || true
sleep 1
nohup go run ./cmd/server config.test.yaml > /tmp/chb_backend.log 2>&1 &
echo "backend pid: $!"
sleep 5
curl -s http://localhost:9090/api/health || tail -20 /tmp/chb_backend.log
