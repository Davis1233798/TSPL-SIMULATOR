#!/bin/bash

echo "================================================"
echo "設置 Cloudflare Authenticated Origin Pulls"
echo "================================================"
echo ""

echo "此腳本將配置 nginx 以支持 Cloudflare Full (strict) 模式"
echo ""

# 步驟 1: 下載 Cloudflare Origin CA 根證書
echo "1. 下載 Cloudflare Origin CA 根證書..."
sudo mkdir -p /etc/ssl/cloudflare
sudo curl -o /etc/ssl/cloudflare/origin-pull-ca.pem https://developers.cloudflare.com/ssl/static/origin_ca_rsa_root.pem
echo "✅ 已下載"
echo ""

# 步驟 2: 驗證證書
echo "2. 驗證證書..."
if openssl x509 -in /etc/ssl/cloudflare/origin-pull-ca.pem -noout -text > /dev/null 2>&1; then
    echo "✅ 證書格式正確"
    openssl x509 -in /etc/ssl/cloudflare/origin-pull-ca.pem -noout -subject -issuer
else
    echo "❌ 證書格式錯誤"
    exit 1
fi
echo ""

# 步驟 3: 更新 nginx 配置以使用 Cloudflare Origin CA
echo "3. 更新 nginx 配置..."
sudo tee /etc/nginx/sites-available/tspl-simulator > /dev/null <<'EOF'
# HTTP 配置
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

# HTTPS 配置 - Cloudflare Origin Pull with Client Certificate Verification
server {
    listen 443 ssl http2;
    server_name tsplsimulator.dpdns.org 138.2.60.98;

    # Cloudflare Origin Certificate
    ssl_certificate /etc/ssl/cloudflare/cert.pem;
    ssl_certificate_key /etc/ssl/cloudflare/key.pem;

    # Cloudflare Authenticated Origin Pulls
    # 驗證連接來自 Cloudflare
    ssl_client_certificate /etc/ssl/cloudflare/origin-pull-ca.pem;
    ssl_verify_client optional;

    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
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
echo "✅ 配置已更新"
echo ""

# 步驟 4: 測試配置
echo "4. 測試 nginx 配置..."
if sudo nginx -t; then
    echo "✅ 配置測試通過"
else
    echo "❌ 配置測試失敗"
    exit 1
fi
echo ""

# 步驟 5: 重載 nginx
echo "5. 重載 nginx..."
sudo systemctl reload nginx
echo "✅ nginx 已重載"
echo ""

echo "================================================"
echo "設置完成!"
echo "================================================"
echo ""
echo "下一步操作:"
echo "1. 在 Cloudflare Dashboard 中:"
echo "   SSL/TLS → Origin Server → Authenticated Origin Pulls"
echo "   開啟 'Authenticated Origin Pulls' 選項"
echo ""
echo "2. 將 SSL/TLS 加密模式改為 'Full (strict)'"
echo ""
echo "3. 測試網站: https://tsplsimulator.dpdns.org"
echo ""
