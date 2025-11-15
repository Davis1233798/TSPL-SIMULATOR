# 後端建置與測試指南

## 前置需求

### 1. 安裝 Go

**Windows:**
1. 下載 Go 安裝程式: https://go.dev/dl/
2. 執行安裝程式 (建議安裝 Go 1.21 或更新版本)
3. 安裝完成後,重新啟動命令提示字元
4. 驗證安裝:
   ```cmd
   go version
   ```

**macOS:**
```bash
brew install go
```

**Linux:**
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

## 建置步驟

### 1. 安裝依賴套件

```bash
cd backend
go mod download
```

這會下載所有必要的 Go 模組:
- Gin Web Framework
- Eclipse Paho MQTT Client
- CORS 中介軟體

### 2. 編譯專案

**開發模式 (直接執行):**
```bash
go run main.go
```

**編譯為執行檔:**

Windows:
```cmd
go build -o tspl-simulator.exe
```

macOS/Linux:
```bash
go build -o tspl-simulator
```

### 3. 執行程式

**開發模式:**
```bash
go run main.go
```

**執行編譯後的檔案:**

Windows:
```cmd
tspl-simulator.exe
```

macOS/Linux:
```bash
./tspl-simulator
```

## 測試

### 1. 測試 API 端點

啟動後端後,可以使用以下方式測試:

**健康檢查:**
```bash
curl http://localhost:8080/api/health
```

**測試 TSPL 渲染:**
```bash
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"Hello TSPL!\"\nPRINT 1,1"
  }'
```

**使用測試檔案:**
```bash
# Windows (PowerShell)
$content = Get-Content test_tspl_example.tspl -Raw
$json = @{tspl_code=$content} | ConvertTo-Json
Invoke-WebRequest -Uri "http://localhost:8080/api/render" -Method POST -Body $json -ContentType "application/json"

# macOS/Linux
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d "{\"tspl_code\": \"$(cat test_tspl_example.tspl | tr '\n' ' ')\"}"
```

### 2. 驗證資料儲存

執行 API 請求後,檢查資料是否正確儲存:

```bash
# Windows
dir /s data\API_print

# macOS/Linux
find data/API_print -type f
```

應該會看到類似以下結構:
```
data/
└── API_print/
    └── 2025_01_15/
        └── 10_30_45.tspl
```

### 3. 測試語法驗證

**測試無效的 TSPL (缺少 SIZE):**
```bash
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "TEXT 100,100,\"3\",0,1,1,\"Hello\"\nPRINT 1,1"
  }'
```

預期回應:
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

**測試無效的參數:**
```bash
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "SIZE 100 mm, 50 mm\nDIRECTION 99\nPRINT 1,1"
  }'
```

預期回應:
```json
{
  "success": false,
  "error": "TSPL 語法驗證失敗",
  "validation_errors": [
    {
      "line": 2,
      "command": "DIRECTION",
      "message": "方向必須在 0-3 之間"
    }
  ]
}
```

## MQTT 測試 (選用)

如果要測試 MQTT 功能,需要先安裝 MQTT Broker:

### 安裝 Mosquitto

**Windows:**
```
下載並安裝: https://mosquitto.org/download/
```

**macOS:**
```bash
brew install mosquitto
brew services start mosquitto
```

**Linux:**
```bash
sudo apt-get install mosquitto mosquitto-clients
sudo systemctl start mosquitto
```

### 測試 MQTT

1. 訂閱結果主題:
```bash
mosquitto_sub -t "tspl/commands/result"
```

2. 發布測試訊息:
```bash
mosquitto_pub -t "tspl/commands" -m '{
  "type": "render_request",
  "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"MQTT Test\"\nPRINT 1,1",
  "timestamp": 1234567890
}'
```

3. 檢查 MQTT 資料儲存:
```bash
# Windows
dir /s data\MQTT_print

# macOS/Linux
find data/MQTT_print -type f
```

## 常見問題

### Q: `go: command not found`
A: Go 尚未安裝或未加入 PATH。請按照上方的安裝步驟操作。

### Q: 編譯錯誤 `cannot find package`
A: 執行 `go mod download` 安裝所有依賴套件。

### Q: 服務器無法啟動 (port already in use)
A: 端口 8080 已被佔用。可以透過環境變數更改:
```bash
export SERVER_PORT=8081
go run main.go
```

### Q: MQTT 連接失敗
A: 確認 MQTT Broker 已啟動,或將 `.env` 中的 `MQTT_BROKER` 設為空值以停用 MQTT。

### Q: 資料儲存失敗
A: 檢查 `data` 資料夾的寫入權限,或透過環境變數指定其他路徑:
```bash
export STORAGE_PATH=/path/to/writable/directory
go run main.go
```

## 生產環境部署

### 使用 Docker

建立 Dockerfile:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o tspl-simulator

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/tspl-simulator .
EXPOSE 8080
CMD ["./tspl-simulator"]
```

建置並執行:
```bash
docker build -t tspl-simulator .
docker run -p 8080:8080 -v $(pwd)/data:/root/data tspl-simulator
```

### 使用 systemd (Linux)

建立服務檔案 `/etc/systemd/system/tspl-simulator.service`:
```ini
[Unit]
Description=TSPL Simulator Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/tspl-simulator
ExecStart=/opt/tspl-simulator/tspl-simulator
Restart=on-failure
Environment="SERVER_PORT=8080"
Environment="STORAGE_PATH=/var/lib/tspl-simulator/data"

[Install]
WantedBy=multi-user.target
```

啟用並啟動:
```bash
sudo systemctl enable tspl-simulator
sudo systemctl start tspl-simulator
sudo systemctl status tspl-simulator
```
