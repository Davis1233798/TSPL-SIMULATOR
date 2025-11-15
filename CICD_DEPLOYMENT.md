# CI/CD 部署指南

## 📋 目錄

1. [GitHub Actions 設置](#github-actions-設置)
2. [Ubuntu 主機準備](#ubuntu-主機準備)
3. [GitHub Secrets 配置](#github-secrets-配置)
4. [部署流程](#部署流程)
5. [Docker 部署](#docker-部署)
6. [故障排除](#故障排除)

---

## 🚀 GitHub Actions 設置

### 工作流程概覽

CI/CD pipeline 包含三個主要階段:

1. **後端測試** - Go 單元測試和構建
2. **前端測試和構建** - React 測試、Lint 和生產構建
3. **自動部署** - 部署到 Ubuntu 主機

### Workflow 文件

已創建 `.github/workflows/ci-cd.yml`

**觸發條件**:
- Push 到 `master` 或 `main` 分支
- Pull Request 到 `master` 或 `main` 分支

**執行步驟**:
```
後端測試 ────┐
             ├──> 部署到 Ubuntu
前端測試 ────┘
```

---

## 🖥️ Ubuntu 主機準備

### 1. 安裝必要軟體

```bash
# 更新套件列表
sudo apt update && sudo apt upgrade -y

# 安裝 Nginx (可選,用於反向代理)
sudo apt install -y nginx

# 安裝 Git (如需手動部署)
sudo apt install -y git

# 安裝 systemd (通常已安裝)
systemctl --version
```

### 2. 創建部署用戶 (可選)

```bash
# 創建專用部署用戶
sudo adduser tspl-deployer

# 賦予 sudo 權限
sudo usermod -aG sudo tspl-deployer

# 切換到部署用戶
su - tspl-deployer
```

### 3. 設置 SSH 金鑰認證

**在本地機器上生成 SSH 金鑰**:

```bash
# 生成新的 SSH 金鑰 (在你的本地電腦)
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/tspl_deploy_key

# 顯示公鑰
cat ~/.ssh/tspl_deploy_key.pub
```

**在 Ubuntu 主機上**:

```bash
# 創建 .ssh 目錄
mkdir -p ~/.ssh
chmod 700 ~/.ssh

# 添加公鑰到 authorized_keys
nano ~/.ssh/authorized_keys
# 貼上剛才生成的公鑰內容

# 設置權限
chmod 600 ~/.ssh/authorized_keys
```

**測試 SSH 連接**:

```bash
# 在本地測試 (替換為你的主機 IP)
ssh -i ~/.ssh/tspl_deploy_key user@your-server-ip
```

### 4. 創建部署目錄

```bash
sudo mkdir -p /opt/tspl-simulator
sudo chown -R $USER:$USER /opt/tspl-simulator
```

### 5. 配置防火牆

```bash
# 允許 SSH
sudo ufw allow 22/tcp

# 允許 HTTP
sudo ufw allow 80/tcp

# 允許 HTTPS (如果需要)
sudo ufw allow 443/tcp

# 允許後端端口 (如果直接訪問)
sudo ufw allow 8080/tcp

# 啟用防火牆
sudo ufw enable
```

---

## 🔐 GitHub Secrets 配置

在 GitHub 倉庫設置以下 Secrets:

### 1. 進入 GitHub 倉庫設置

```
你的倉庫 → Settings → Secrets and variables → Actions → New repository secret
```

### 2. 添加以下 Secrets

| Secret 名稱 | 說明 | 範例值 |
|------------|------|--------|
| `DEPLOY_HOST` | Ubuntu 主機 IP 或域名 | `192.168.1.100` 或 `tspl.example.com` |
| `DEPLOY_USER` | SSH 登入用戶名 | `ubuntu` 或 `tspl-deployer` |
| `DEPLOY_SSH_KEY` | SSH 私鑰內容 | 從 `~/.ssh/tspl_deploy_key` 複製完整內容 |
| `DEPLOY_PORT` | SSH 端口 (可選) | `22` (默認) |

### 3. 複製 SSH 私鑰

**Windows (PowerShell)**:
```powershell
Get-Content ~/.ssh/tspl_deploy_key | clip
```

**Linux/Mac**:
```bash
cat ~/.ssh/tspl_deploy_key | pbcopy  # Mac
cat ~/.ssh/tspl_deploy_key | xclip   # Linux
```

**重要**: 複製私鑰時,包含以下格式:
```
-----BEGIN OPENSSH PRIVATE KEY-----
... 私鑰內容 ...
-----END OPENSSH PRIVATE KEY-----
```

---

## 📦 部署流程

### 自動部署 (推送到 GitHub 後)

1. **推送代碼到 GitHub**:
```bash
git add .
git commit -m "部署更新"
git push origin master
```

2. **GitHub Actions 自動執行**:
- ✅ 運行後端測試
- ✅ 運行前端測試和構建
- ✅ 部署到 Ubuntu 主機
- ✅ 啟動服務

3. **查看部署狀態**:
- GitHub 倉庫 → Actions 標籤
- 查看最新的 workflow 運行

### 手動部署 (在 Ubuntu 主機上)

```bash
# 1. SSH 連接到主機
ssh user@your-server-ip

# 2. 進入部署目錄
cd /opt/tspl-simulator

# 3. 執行部署腳本
chmod +x deploy/deploy.sh
./deploy/deploy.sh
```

---

## 🐳 Docker 部署

### 方法 1: 使用 Docker Compose (推薦)

```bash
# 1. 克隆倉庫
git clone https://github.com/Davis1233798/TSPL-SIMULATOR.git
cd TSPL-SIMULATOR

# 2. 創建環境變數文件
cp backend/.env.example backend/.env
nano backend/.env  # 編輯配置

# 3. 構建並啟動
docker-compose up -d --build

# 4. 查看日誌
docker-compose logs -f tspl-simulator

# 5. 停止服務
docker-compose down
```

### 方法 2: 僅使用 Docker

```bash
# 1. 構建映像
docker build -t tspl-simulator .

# 2. 運行容器
docker run -d \
  --name tspl-simulator \
  -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/backend/.env:/root/.env \
  tspl-simulator

# 3. 查看日誌
docker logs -f tspl-simulator

# 4. 停止容器
docker stop tspl-simulator
docker rm tspl-simulator
```

### Docker 命令

```bash
# 查看運行中的容器
docker ps

# 進入容器
docker exec -it tspl-simulator sh

# 重啟容器
docker restart tspl-simulator

# 查看資源使用
docker stats tspl-simulator
```

---

## 🔍 服務管理

### Systemd 服務命令

```bash
# 啟動服務
sudo systemctl start tspl-simulator

# 停止服務
sudo systemctl stop tspl-simulator

# 重啟服務
sudo systemctl restart tspl-simulator

# 查看狀態
sudo systemctl status tspl-simulator

# 啟用開機自啟
sudo systemctl enable tspl-simulator

# 禁用開機自啟
sudo systemctl disable tspl-simulator

# 查看日誌
sudo journalctl -u tspl-simulator -f

# 查看最近的日誌
sudo journalctl -u tspl-simulator -n 100
```

### Nginx 命令

```bash
# 測試配置
sudo nginx -t

# 重新加載配置
sudo systemctl reload nginx

# 重啟 Nginx
sudo systemctl restart nginx

# 查看狀態
sudo systemctl status nginx

# 查看錯誤日誌
sudo tail -f /var/log/nginx/error.log

# 查看訪問日誌
sudo tail -f /var/log/nginx/access.log
```

---

## 🐛 故障排除

### 問題 1: GitHub Actions 部署失敗

**症狀**: SSH 連接失敗

**解決方案**:
1. 檢查 GitHub Secrets 是否正確設置
2. 確認 SSH 私鑰格式正確 (包含 BEGIN 和 END 標記)
3. 測試本地 SSH 連接:
```bash
ssh -i ~/.ssh/tspl_deploy_key user@your-server-ip
```

### 問題 2: 服務啟動失敗

**症狀**: `systemctl status tspl-simulator` 顯示 failed

**解決方案**:
```bash
# 查看詳細錯誤
sudo journalctl -u tspl-simulator -n 50

# 檢查文件權限
ls -la /opt/tspl-simulator/backend/tspl-simulator

# 確保可執行
chmod +x /opt/tspl-simulator/backend/tspl-simulator

# 檢查環境變數
cat /opt/tspl-simulator/backend/.env
```

### 問題 3: 端口被占用

**症狀**: `bind: address already in use`

**解決方案**:
```bash
# 查看占用端口的進程
sudo lsof -i :8080

# 殺死進程
sudo kill -9 <PID>

# 或更改端口
nano /opt/tspl-simulator/backend/.env
# 修改 SERVER_PORT=8081
```

### 問題 4: Nginx 502 Bad Gateway

**症狀**: 前端可以訪問,API 返回 502

**解決方案**:
```bash
# 確認後端服務運行中
sudo systemctl status tspl-simulator

# 檢查後端是否監聽正確端口
netstat -tlnp | grep 8080

# 檢查 Nginx 錯誤日誌
sudo tail -f /var/log/nginx/error.log

# 測試後端健康檢查
curl http://localhost:8080/api/health
```

### 問題 5: 前端連接後端失敗

**症狀**: 前端顯示 "後端不可用"

**解決方案**:
```bash
# 1. 檢查前端環境變數
cat /opt/tspl-simulator/frontend/build/.env

# 2. 確認 API URL 正確
# 應該是: REACT_APP_API_URL=http://your-domain/api

# 3. 測試 CORS
curl -H "Origin: http://your-domain" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8080/api/render
```

---

## 📊 監控和日誌

### 實時監控

```bash
# 系統資源
htop

# 服務日誌
sudo journalctl -u tspl-simulator -f

# Nginx 訪問日誌
sudo tail -f /var/log/nginx/access.log

# 後端儲存的文件
watch -n 1 'ls -lh /opt/tspl-simulator/backend/data/API_print/*/*.tspl'
```

### 日誌位置

| 服務 | 日誌位置 |
|------|---------|
| systemd | `sudo journalctl -u tspl-simulator` |
| Nginx 訪問 | `/var/log/nginx/access.log` |
| Nginx 錯誤 | `/var/log/nginx/error.log` |
| TSPL 文件 | `/opt/tspl-simulator/backend/data/` |

---

## 🔄 更新部署

### 自動更新 (推送到 GitHub)

```bash
# 本地修改代碼
git add .
git commit -m "更新功能"
git push origin master

# GitHub Actions 自動部署
```

### 手動更新

```bash
# SSH 到主機
ssh user@your-server-ip

# 拉取最新代碼
cd /opt/tspl-simulator
git pull origin master

# 運行部署腳本
./deploy/deploy.sh
```

---

## ✅ 部署檢查清單

- [ ] Ubuntu 主機已安裝必要軟體
- [ ] SSH 金鑰認證已設置
- [ ] GitHub Secrets 已配置
- [ ] 防火牆規則已設置
- [ ] 部署目錄已創建
- [ ] systemd 服務文件已安裝
- [ ] Nginx 配置已設置 (如使用)
- [ ] 環境變數 `.env` 已配置
- [ ] 推送到 GitHub 觸發 CI/CD
- [ ] 服務成功啟動
- [ ] 前端可以訪問
- [ ] API 端點可以訪問
- [ ] 文件儲存功能正常

---

## 📞 快速指令參考

```bash
# 查看服務狀態
sudo systemctl status tspl-simulator

# 查看實時日誌
sudo journalctl -u tspl-simulator -f

# 重啟服務
sudo systemctl restart tspl-simulator

# 測試 API
curl http://localhost:8080/api/health

# 查看儲存的文件
ls -lh /opt/tspl-simulator/backend/data/API_print/

# 檢查端口
netstat -tlnp | grep 8080

# 查看 Docker 容器 (如使用 Docker)
docker ps
docker logs -f tspl-simulator
```

---

## 🎉 完成!

你的 TSPL Simulator 現在已經配置了完整的 CI/CD pipeline!

- ✅ 自動測試
- ✅ 自動構建
- ✅ 自動部署
- ✅ 生產就緒

每次推送到 GitHub 後,應用會自動部署到 Ubuntu 主機! 🚀
