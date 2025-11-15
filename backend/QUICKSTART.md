# 快速入門指南

## 5 分鐘快速開始

### 1. 安裝 Go (如果尚未安裝)

訪問 https://go.dev/dl/ 下載並安裝 Go 1.21+

### 2. 啟動後端

```bash
cd backend
go mod download
go run main.go
```

你應該會看到:
```
2025/01/15 10:30:00 儲存服務已初始化,資料路徑: ./data
2025/01/15 10:30:00 MQTT 未配置,僅啟用 API 功能
2025/01/15 10:30:00 TSPL Simulator 服務器啟動於 :8080
2025/01/15 10:30:00 API 端點: http://localhost:8080/api
2025/01/15 10:30:00 API 資料儲存: data\API_print
2025/01/15 10:30:00 MQTT 資料儲存: data\MQTT_print
```

### 3. 測試 API

開啟新的終端視窗,測試健康檢查:

```bash
curl http://localhost:8080/api/health
```

預期回應:
```json
{
  "status": "ok",
  "mqtt": "disconnected"
}
```

### 4. 測試 TSPL 渲染

**使用 curl (macOS/Linux):**
```bash
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"Hello TSPL!\"\nPRINT 1,1"
  }'
```

**使用 PowerShell (Windows):**
```powershell
$body = @{
    tspl_code = "SIZE 100 mm, 50 mm`nGAP 3 mm, 0 mm`nCLS`nTEXT 100,100,`"3`",0,1,1,`"Hello TSPL!`"`nPRINT 1,1"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/render" -Method POST -Body $body -ContentType "application/json"
```

**使用瀏覽器工具 (推薦):**

安裝瀏覽器擴充如 [Postman](https://www.postman.com/) 或使用內建的開發者工具:

1. 開啟瀏覽器
2. 按 F12 開啟開發者工具
3. 切換到 Console 標籤
4. 貼上以下程式碼:

```javascript
fetch('http://localhost:8080/api/render', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    tspl_code: 'SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,"3",0,1,1,"Hello TSPL!"\nPRINT 1,1'
  })
})
.then(r => r.json())
.then(console.log)
```

### 5. 檢查儲存的檔案

```bash
# Windows
dir /s data\API_print

# macOS/Linux
ls -R data/API_print
```

你應該會看到:
```
data/API_print/2025_01_15/10_30_45.tspl
```

檔案內容就是剛才發送的 TSPL 程式碼。

## 核心功能演示

### 功能 1: 語法驗證

發送一個**無效**的 TSPL (缺少 SIZE):

```bash
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "TEXT 100,100,\"3\",0,1,1,\"Hello\"\nPRINT 1,1"
  }'
```

回應:
```json
{
  "success": false,
  "error": "TSPL 語法驗證失敗",
  "validation_errors": [
    {
      "line": 0,
      "command": "SIZE",
      "message": "缺少必要的 SIZE 命令"
    }
  ]
}
```

### 功能 2: 自動檔案儲存

每次 API 請求都會自動儲存:
- **路徑**: `data/API_print/年_月_日/時_分_秒.tspl`
- **範例**: `data/API_print/2025_01_15/10_30_45.tspl`

### 功能 3: 完整的 TSPL 範例

使用提供的測試檔案:

```bash
# Windows PowerShell
$tsplCode = Get-Content test_tspl_example.tspl -Raw
$body = @{tspl_code=$tsplCode} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/render" -Method POST -Body $body -ContentType "application/json"

# macOS/Linux
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d "{\"tspl_code\": \"$(cat test_tspl_example.tspl)\"}"
```

## 前後端整合測試

### 1. 啟動前端 (另一個終端)

```bash
cd ../frontend
npm install
npm start
```

### 2. 開啟瀏覽器

訪問 http://localhost:3000

### 3. 測試整合

1. 在前端輸入 TSPL 程式碼
2. 點擊「渲染」按鈕
3. 查看渲染結果
4. 檢查後端資料夾確認檔案已儲存

## MQTT 功能 (進階)

如果需要 MQTT 功能:

### 1. 安裝 Mosquitto

**Windows**: https://mosquitto.org/download/
**macOS**: `brew install mosquitto && brew services start mosquitto`
**Linux**: `sudo apt install mosquitto && sudo systemctl start mosquitto`

### 2. 配置後端

建立 `.env` 檔案:
```bash
SERVER_PORT=8080
STORAGE_PATH=./data
MQTT_BROKER=localhost
MQTT_PORT=1883
MQTT_CLIENT_ID=tspl-simulator
MQTT_TOPIC=tspl/commands
```

### 3. 重啟後端

```bash
go run main.go
```

你應該會看到:
```
2025/01/15 10:30:00 MQTT 客戶端已連接到 tcp://localhost:1883 並訂閱主題 tspl/commands
2025/01/15 10:30:00 MQTT 客戶端已成功初始化
```

### 4. 測試 MQTT

**訂閱結果:**
```bash
mosquitto_sub -t "tspl/commands/result"
```

**發布訊息:**
```bash
mosquitto_pub -t "tspl/commands" -m '{
  "type": "render_request",
  "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"MQTT Test\"\nPRINT 1,1",
  "timestamp": 1234567890
}'
```

**檢查 MQTT 資料儲存:**
```bash
ls -R data/MQTT_print
```

## 下一步

1. 閱讀 [README.md](README.md) 了解完整功能
2. 查看 [BUILD.md](BUILD.md) 了解建置和部署
3. 探索 [API 文件](#api-端點) 了解所有端點

## 疑難排解

### 問題: `go: command not found`
**解決**: 安裝 Go 或將 Go 加入 PATH

### 問題: `port 8080 already in use`
**解決**:
```bash
export SERVER_PORT=8081
go run main.go
```

### 問題: 前端無法連接後端
**解決**:
1. 確認後端正在運行
2. 檢查前端的 `.env` 檔案中的 `REACT_APP_API_URL`
3. 檢查 CORS 設定

### 問題: 檔案儲存失敗
**解決**: 檢查 `data` 資料夾權限,或更改 `STORAGE_PATH` 環境變數

## 聯絡與支援

遇到問題?
1. 查看 [README.md](README.md) 的常見問題章節
2. 檢查終端的錯誤訊息
3. 確認所有依賴都已正確安裝
