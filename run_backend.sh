#!/bin/bash
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
echo "Starting backend on :9090 with config.test.yaml"
exec go run ./cmd/server config.test.yaml
