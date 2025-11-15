#!/bin/bash

# TSPL Simulator 部署腳本
# 用於 Ubuntu 主機自動部署

set -e

echo "========================================="
echo "TSPL Simulator 部署腳本"
echo "========================================="

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
DEPLOY_DIR="/opt/tspl-simulator"
SERVICE_NAME="tspl-simulator"
NGINX_SITE_CONFIG="/etc/nginx/sites-available/tspl-simulator"
NGINX_SITE_ENABLED="/etc/nginx/sites-enabled/tspl-simulator"

# 函數: 打印成功訊息
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# 函數: 打印錯誤訊息
print_error() {
    echo -e "${RED}✗${NC} $1"
}

# 函數: 打印警告訊息
print_warning() {
    echo -e "${YELLOW}!${NC} $1"
}

# 檢查是否為 root 或有 sudo 權限
if [[ $EUID -ne 0 ]]; then
   if ! sudo -n true 2>/dev/null; then
       print_error "此腳本需要 sudo 權限"
       exit 1
   fi
fi

# 1. 創建部署目錄
echo ""
echo "步驟 1: 創建部署目錄..."
sudo mkdir -p $DEPLOY_DIR/backend
sudo mkdir -p $DEPLOY_DIR/frontend
sudo mkdir -p $DEPLOY_DIR/backend/data/API_print
sudo mkdir -p $DEPLOY_DIR/backend/data/MQTT_print
print_success "部署目錄已創建"

# 2. 設置權限
echo ""
echo "步驟 2: 設置文件權限..."
sudo chown -R $USER:$USER $DEPLOY_DIR
chmod +x $DEPLOY_DIR/backend/tspl-simulator 2>/dev/null || true
print_success "權限已設置"

# 3. 複製環境變數文件（如果不存在）
echo ""
echo "步驟 3: 配置環境變數..."
if [ ! -f "$DEPLOY_DIR/backend/.env" ]; then
    if [ -f "$DEPLOY_DIR/backend/.env.example" ]; then
        cp $DEPLOY_DIR/backend/.env.example $DEPLOY_DIR/backend/.env
        print_success "已從 .env.example 創建 .env 文件"
        print_warning "請編輯 $DEPLOY_DIR/backend/.env 配置您的環境變數"
    else
        print_warning ".env.example 不存在,跳過環境變數配置"
    fi
else
    print_success "環境變數文件已存在"
fi

# 4. 安裝並配置 systemd 服務
echo ""
echo "步驟 4: 配置 systemd 服務..."
if [ -f "$(dirname $0)/tspl-simulator.service" ]; then
    sudo cp $(dirname $0)/tspl-simulator.service /etc/systemd/system/
    sudo systemctl daemon-reload
    print_success "systemd 服務文件已安裝"
else
    print_warning "找不到 tspl-simulator.service,請手動配置"
fi

# 5. 停止舊服務
echo ""
echo "步驟 5: 停止舊服務..."
if sudo systemctl is-active --quiet $SERVICE_NAME; then
    sudo systemctl stop $SERVICE_NAME
    print_success "舊服務已停止"
else
    print_warning "服務未運行"
fi

# 6. 啟動新服務
echo ""
echo "步驟 6: 啟動新服務..."
sudo systemctl start $SERVICE_NAME
sudo systemctl enable $SERVICE_NAME

# 等待服務啟動
sleep 2

if sudo systemctl is-active --quiet $SERVICE_NAME; then
    print_success "服務已成功啟動"
else
    print_error "服務啟動失敗"
    sudo systemctl status $SERVICE_NAME
    exit 1
fi

# 7. 配置 Nginx (可選)
echo ""
echo "步驟 7: 配置 Nginx..."
if command -v nginx &> /dev/null; then
    # 創建 Nginx 配置
    sudo tee $NGINX_SITE_CONFIG > /dev/null <<EOF
server {
    listen 80;
    server_name _;

    # 前端靜態文件
    location / {
        root $DEPLOY_DIR/frontend/build;
        index index.html;
        try_files \$uri \$uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    # Gzip 壓縮
    gzip on;
    gzip_vary on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
}
EOF

    # 啟用站點
    sudo ln -sf $NGINX_SITE_CONFIG $NGINX_SITE_ENABLED

    # 測試 Nginx 配置
    if sudo nginx -t 2>/dev/null; then
        sudo systemctl reload nginx
        print_success "Nginx 已配置並重新加載"
    else
        print_error "Nginx 配置測試失敗"
    fi
else
    print_warning "未安裝 Nginx,跳過 Nginx 配置"
fi

# 8. 顯示服務狀態
echo ""
echo "========================================="
echo "部署完成!"
echo "========================================="
echo ""
echo "服務狀態:"
sudo systemctl status $SERVICE_NAME --no-pager

echo ""
echo "查看日誌:"
echo "  sudo journalctl -u $SERVICE_NAME -f"

echo ""
echo "服務管理:"
echo "  啟動: sudo systemctl start $SERVICE_NAME"
echo "  停止: sudo systemctl stop $SERVICE_NAME"
echo "  重啟: sudo systemctl restart $SERVICE_NAME"
echo "  狀態: sudo systemctl status $SERVICE_NAME"

echo ""
print_success "部署完成! 🎉"
