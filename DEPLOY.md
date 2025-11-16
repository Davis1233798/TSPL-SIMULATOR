# CI/CD 自動部署到 Ubuntu 主機

## 📋 前置準備

### 1. GitHub 端設置

#### 設置 GitHub Secrets
進入倉庫: `Settings → Secrets and variables → Actions → New repository secret`

添加以下 4 個 secrets:

| Secret 名稱 | 值 | 說明 |
|------------|-----|------|
| `DEPLOY_HOST` | `138.2.60.98` | Ubuntu 主機 IP 或域名 |
| `DEPLOY_USER` | `ubuntu` | SSH 登入用戶名 |
| `DEPLOY_SSH_KEY` | SSH 私鑰完整內容 | 用於 SSH 認證 |
| `DEPLOY_PORT` | `22` | SSH 端口 (可選,默認 22) |

#### 生成部署 SSH 金鑰

```powershell
# 生成專用部署金鑰
ssh-keygen -t ed25519 -C "github-deploy" -f C:\Users\solidityDeveloper\.ssh\github_deploy_key

# 顯示私鑰 (添加到 GitHub Secret: DEPLOY_SSH_KEY)
cat C:\Users\solidityDeveloper\.ssh\github_deploy_key

# 顯示公鑰 (添加到 Ubuntu 主機)
cat C:\Users\solidityDeveloper\.ssh\github_deploy_key.pub
```

---

### 2. Ubuntu 主機準備

#### SSH 到主機並執行:

```bash
# 更新系統
sudo apt update && sudo apt upgrade -y

# 安裝 Nginx
sudo apt install -y nginx

# 創建部署目錄
sudo mkdir -p /opt/tspl-simulator && sudo chown -R $USER:$USER /opt/tspl-simulator

# 添加 GitHub Actions 部署公鑰
mkdir -p ~/.ssh
nano ~/.ssh/authorized_keys
# 貼上剛才生成的公鑰,保存

# 設置權限
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys

# 配置防火牆
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw allow 8080/tcp # 後端
sudo ufw enable
```

#### 創建 systemd 服務文件

```bash
sudo nano /etc/systemd/system/tspl-simulator.service
```

內容:
```ini
[Unit]
Description=TSPL Simulator Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/tspl-simulator/backend
EnvironmentFile=/opt/tspl-simulator/backend/.env
ExecStart=/opt/tspl-simulator/backend/tspl-simulator
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

保存後:
```bash
sudo systemctl daemon-reload
```

#### 配置 Nginx

**注意**: CI/CD 會自動部署 Nginx 配置,這裡僅供參考或手動配置時使用。

CI/CD 會自動從 `deploy/nginx-site.conf` 部署配置文件到:
- `/etc/nginx/sites-available/tspl-simulator`
- `/etc/nginx/sites-enabled/tspl-simulator`

如需手動配置:
```bash
sudo nano /etc/nginx/sites-available/tspl-simulator
```

內容參考 `deploy/nginx-site.conf`:
- 前端靜態文件: `/opt/tspl-simulator/frontend/build`
- API 代理: `http://127.0.0.1:8080`
- 包含 Gzip 壓縮、快取設定、安全標頭

手動啟用配置:
```bash
sudo ln -s /etc/nginx/sites-available/tspl-simulator /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 🚀 部署流程

### 自動部署 (推送到 GitHub 後)

```powershell
# 1. 修改代碼後提交
git add .
git commit -m "更新功能"

# 2. 推送到 GitHub (自動觸發 CI/CD)
git push origin master
```

GitHub Actions 會自動:
1. ✅ 運行後端測試 (Go)
2. ✅ 運行前端測試 (React)
3. ✅ 構建前後端
4. ✅ SSH 到 Ubuntu 主機
5. ✅ 部署應用
6. ✅ 重啟服務

### 查看部署進度

訪問: `https://github.com/Davis1233798/TSPL-SIMULATOR/actions`

---

## 🔍 驗證部署

### 在 Ubuntu 主機檢查

```bash
# 檢查服務狀態
sudo systemctl status tspl-simulator

# 查看日誌
sudo journalctl -u tspl-simulator -f

# 測試 API
curl http://localhost:8080/api/health

# 檢查文件
ls -la /opt/tspl-simulator/backend/
```

### 從瀏覽器訪問

```
http://your-server-ip/          # 前端
http://your-server-ip/api/health # API 健康檢查
```

---

## 🐛 故障排除

### 問題 1: GitHub Actions 部署失敗

```bash
# 檢查 SSH 連接
ssh -i ~/.ssh/github_deploy_key ubuntu@your-server-ip

# 檢查 Secrets 是否正確配置
# GitHub: Settings → Secrets and variables → Actions
```

### 問題 2: 服務啟動失敗

```bash
# 查看詳細錯誤
sudo journalctl -u tspl-simulator -n 50

# 確保環境變數存在
cat /opt/tspl-simulator/backend/.env

# 手動測試運行
cd /opt/tspl-simulator/backend
./tspl-simulator
```

### 問題 3: Nginx 502 錯誤

```bash
# 確認後端運行中
sudo systemctl status tspl-simulator

# 檢查端口
netstat -tlnp | grep 8080

# 查看 Nginx 日誌
sudo tail -f /var/log/nginx/error.log
```

---

## 📊 服務管理

```bash
# 啟動
sudo systemctl start tspl-simulator

# 停止
sudo systemctl stop tspl-simulator

# 重啟
sudo systemctl restart tspl-simulator

# 查看狀態
sudo systemctl status tspl-simulator

# 開機自啟
sudo systemctl enable tspl-simulator

# 查看日誌
sudo journalctl -u tspl-simulator -f
```

---

## ✅ 完成!

現在每次你推送代碼到 GitHub,應用會自動部署到 Ubuntu 主機! 🎉
