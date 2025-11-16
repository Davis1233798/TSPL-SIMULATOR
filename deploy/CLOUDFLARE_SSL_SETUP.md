# Cloudflare SSL 憑證設置指南

## 前提條件

您需要從 Cloudflare 獲取 Origin Certificate（原始伺服器憑證）。

## 步驟 1: 下載 Cloudflare Origin Certificate

1. 登入 Cloudflare Dashboard
2. 選擇您的域名 (tsplsimulator.dpdns.org)
3. 進入 SSL/TLS → Origin Server
4. 點擊 "Create Certificate"
5. 使用默認設置（RSA, 15年有效期）
6. 複製以下內容：
   - **Origin Certificate** (cert.pem)
   - **Private Key** (key.pem)

## 步驟 2: 在 Ubuntu 伺服器上創建憑證文件

SSH 到您的伺服器並執行：

```bash
# 創建目錄
sudo mkdir -p /etc/ssl/cloudflare
sudo chmod 755 /etc/ssl/cloudflare

# 創建憑證文件
sudo nano /etc/ssl/cloudflare/cert.pem
# 貼上 Origin Certificate 內容，保存 (Ctrl+O, Enter, Ctrl+X)

# 創建私鑰文件
sudo nano /etc/ssl/cloudflare/key.pem
# 貼上 Private Key 內容，保存

# 設置權限
sudo chmod 644 /etc/ssl/cloudflare/cert.pem
sudo chmod 600 /etc/ssl/cloudflare/key.pem
sudo chown root:root /etc/ssl/cloudflare/*
```

## 步驟 3: 開放防火牆端口

```bash
# 允許 HTTPS 流量
sudo ufw allow 443/tcp

# 檢查防火牆狀態
sudo ufw status
```

## 步驟 4: 重新部署或手動重載 Nginx

### 選項 A: 通過 CI/CD 自動部署（推薦）

只需推送代碼到 GitHub，CI/CD 會自動更新 Nginx 配置。

### 選項 B: 手動更新（如果憑證已配置）

```bash
# 測試 Nginx 配置
sudo nginx -t

# 重載 Nginx
sudo systemctl reload nginx

# 檢查 Nginx 狀態
sudo systemctl status nginx
```

## 步驟 5: 配置 Cloudflare SSL/TLS 設置

在 Cloudflare Dashboard 中：

1. 進入 SSL/TLS 設置
2. 加密模式選擇：**Full (strict)**
   - 不要選 "Flexible"（會導致無限重定向）
   - 不要選 "Full"（不夠安全）

## 驗證

### 檢查憑證是否正確安裝

```bash
# 檢查憑證文件
sudo ls -la /etc/ssl/cloudflare/

# 測試 SSL 配置
sudo nginx -t

# 檢查 443 端口是否監聽
sudo ss -tlnp | grep 443
```

### 測試網站

訪問：https://tsplsimulator.dpdns.org

應該顯示：
- ✅ 綠色鎖頭圖標（安全連接）
- ✅ 無 521 錯誤
- ✅ 網站正常載入

## 故障排除

### 問題 1: 仍然出現 521 錯誤

```bash
# 檢查 Nginx 錯誤日誌
sudo tail -f /var/log/nginx/error.log

# 檢查 Nginx 是否監聽 443
sudo netstat -tlnp | grep 443

# 確認憑證文件存在
ls -la /etc/ssl/cloudflare/
```

### 問題 2: SSL 憑證錯誤

確保：
1. 憑證文件格式正確（包含 `-----BEGIN CERTIFICATE-----` 等）
2. 私鑰文件格式正確（包含 `-----BEGIN PRIVATE KEY-----` 等）
3. 文件權限正確（cert.pem: 644, key.pem: 600）

### 問題 3: Cloudflare SSL 模式錯誤

確保 Cloudflare 的 SSL/TLS 加密模式設為 **Full (strict)**，而不是 Flexible。

## 安全建議

- ✅ 私鑰文件 (key.pem) 權限必須是 600
- ✅ 只有 root 可以讀取私鑰
- ✅ 定期更新憑證（Cloudflare Origin Certificate 最長 15 年）
- ✅ 使用 TLS 1.2 或更高版本
