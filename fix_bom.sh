#!/bin/bash
MIGRATIONS="/mnt/c/Users/PC/Desktop/claudecode/xiaoyuan/chb-backend/migrations"
echo "=== Checking BOM in SQL files ==="
for f in "$MIGRATIONS"/*.up.sql; do
    fname=$(basename "$f")
    header=$(xxd -l 4 "$f" | head -1)
    echo "$fname: $header"
done
