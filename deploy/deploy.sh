#!/bin/bash
set -euo pipefail

REPO_URL="https://github.com/uupup-cn/CampusHub.git"
DEPLOY_DIR="/opt/campushub"
GO_VERSION="go1.24.6"

echo "=== CampusHub Deploy ==="

# Install Go if not present
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    curl -fsSL "https://go.dev/dl/$GO_VERSION.linux-amd64.tar.gz" -o /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh
    source /etc/profile.d/go.sh
fi
export PATH=$PATH:/usr/local/go/bin
export GOPROXY=https://goproxy.cn,direct

# Install Node.js if not present
if ! command -v node &> /dev/null; then
    echo "Installing Node.js 22..."
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash -
    apt-get install -y nodejs
fi

# Clone or pull repo
if [ -d "$DEPLOY_DIR/.git" ]; then
    echo "Pulling latest code..."
    cd $DEPLOY_DIR
    git pull origin main
else
    echo "Cloning repo..."
    git clone $REPO_URL $DEPLOY_DIR
    cd $DEPLOY_DIR
fi

echo "=== Database ==="
docker exec chb-postgres psql -U chb -c "CREATE DATABASE chb_platform;" 2>/dev/null || true

echo "=== Migrations ==="
cd $DEPLOY_DIR/chb-backend
go run ./cmd/migrate /deploy/config.prod.yaml

echo "=== Build Backend ==="
go build -o $DEPLOY_DIR/chb-backend ./cmd/server

echo "=== Build Frontend ==="
cd $DEPLOY_DIR/chb-frontend
npm ci
NEXT_PUBLIC_API_BASE=http://112.213.106.104/api npm run build

echo "=== Services ==="
cp $DEPLOY_DIR/deploy/systemd/chb-backend.service /etc/systemd/system/
cp $DEPLOY_DIR/deploy/systemd/chb-frontend.service /etc/systemd/system/
systemctl daemon-reload

echo "=== Nginx ==="
if ! command -v nginx &> /dev/null; then
    echo "Installing Nginx..."
    apt-get install -y nginx
fi
cp $DEPLOY_DIR/deploy/nginx/nginx.conf /etc/nginx/sites-available/campushub
ln -sf /etc/nginx/sites-available/campushub /etc/nginx/sites-enabled/campushub
rm -f /etc/nginx/sites-enabled/default
nginx -t

echo "=== Restart ==="
systemctl restart chb-backend
systemctl restart chb-frontend
systemctl restart nginx
systemctl enable chb-backend chb-frontend nginx

echo "=== Verify ==="
sleep 3
curl -s http://localhost:9090/health && echo "" || echo "Backend health check failed"
curl -sI http://localhost:3000 | head -1 || echo "Frontend check failed"
curl -s http://localhost/api/health && echo "" || echo "Nginx proxy check failed"

echo "=== Deploy Complete ==="
echo "Frontend: http://112.213.106.104"
echo "API:      http://112.213.106.104/api/health"
echo "Forum:    http://112.213.106.104/forum/ (Discourse not deployed yet)"
