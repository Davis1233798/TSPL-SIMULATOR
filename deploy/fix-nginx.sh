#!/bin/bash

echo "================================================"
echo "修復 Nginx SSL 配置 - Cloudflare 優化版"
echo "================================================"
echo ""

# 備份當前配置
echo "1. 備份當前配置..."
sudo cp /etc/nginx/sites-available/tspl-simulator /etc/nginx/sites-available/tspl-simulator.backup.$(date +%Y%m%d_%H%M%S) 2>/dev/null || true
echo "✅ 已備份"
echo ""

# 創建新的配置文件
echo "2. 創建 Cloudflare 優化配置..."
sudo tee /etc/nginx/sites-available/tspl-simulator > /dev/null <<'EOF'
# HTTP 配置 - 處理所有 HTTP 流量
server {
    listen 80;
    server_name tsplsimulator.dpdns.org 138.2.60.98;

    # 前端靜態文件
    root /opt/tspl-simulator/frontend/build;
    index index.html;

    # Gzip 壓縮
    gzip on;
    gzip_vary on;
    gzip_min_length 10240;
    gzip_proxied expired no-cache no-store private auth;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;

        # 快取靜態資源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # API 代理到後端
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超時設定
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康檢查
    location /health {
        proxy_pass http://127.0.0.1:8080/api/health;
        access_log off;
    }

    # 安全標頭
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # 日誌
    access_log /var/log/nginx/tspl-simulator-access.log;
    error_log /var/log/nginx/tspl-simulator-error.log;
}

# HTTPS 配置 - 僅用於 Cloudflare Origin Pull
server {
    listen 443 ssl http2;
    server_name tsplsimulator.dpdns.org 138.2.60.98;

    # Cloudflare Origin Certificate
    ssl_certificate /etc/ssl/cloudflare/cert.pem;
    ssl_certificate_key /etc/ssl/cloudflare/key.pem;

    # SSL 配置 - Cloudflare 兼容設置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1h;
    ssl_session_tickets off;

    # 前端靜態文件
    root /opt/tspl-simulator/frontend/build;
    index index.html;

    # Gzip 壓縮
    gzip on;
    gzip_vary on;
    gzip_min_length 10240;
    gzip_proxied expired no-cache no-store private auth;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;

        # 快取靜態資源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # API 代理到後端
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超時設定
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康檢查
    location /health {
        proxy_pass http://127.0.0.1:8080/api/health;
        access_log off;
    }

    # 安全標頭
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # 日誌
    access_log /var/log/nginx/tspl-simulator-access.log;
    error_log /var/log/nginx/tspl-simulator-error.log;
}
EOF
echo "✅ 配置已創建"
echo ""

# 測試配置
echo "3. 測試 Nginx 配置..."
if sudo nginx -t; then
    echo "✅ 配置測試通過"
else
    echo "❌ 配置測試失敗,恢復備份"
    sudo cp /etc/nginx/sites-available/tspl-simulator.backup.* /etc/nginx/sites-available/tspl-simulator
    exit 1
fi
echo ""

# 重載 Nginx
echo "4. 重載 Nginx..."
sudo systemctl reload nginx
echo "✅ Nginx 已重載"
echo ""

# 驗證
echo "5. 驗證服務..."
echo "---"
echo "端口監聽狀態:"
sudo netstat -tlnp | grep -E ":(80|443)" | grep nginx
echo ""
echo "本地 HTTPS 測試:"
curl -k -I https://127.0.0.1/ 2>&1 | grep -E "(HTTP|Server)"
echo ""

echo "================================================"
echo "修復完成!"
echo "================================================"
echo ""
echo "請訪問 https://tsplsimulator.dpdns.org 測試"
echo ""
