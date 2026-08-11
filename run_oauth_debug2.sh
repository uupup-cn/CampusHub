#!/bin/bash
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
go test -count=1 -timeout 120s -run "TestOAuthConfirmFlow" -v . > /tmp/oauth_test.log 2>&1
echo "exit=$?"
tail -60 /tmp/oauth_test.log
