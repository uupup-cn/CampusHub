#!/bin/bash
echo "=== Checking Go ==="
ls -la /usr/local/go/bin/go 2>/dev/null && echo "Found /usr/local/go"
ls -la $HOME/go-install/go/bin/go 2>/dev/null && echo "Found go-install"
ls -la $HOME/go/bin/go 2>/dev/null && echo "Found go/bin"
which go 2>/dev/null && echo "go in PATH"
echo "---"
echo "PATH=$PATH"
