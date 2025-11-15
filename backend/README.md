# TSPL Simulator Backend

Go 語言實作的 TSPL 標籤模擬器後端服務

## 功能特色

- ✅ RESTful API 端點
- ✅ TSPL2 語法驗證與解析
- ✅ MQTT 訂閱和發布
- ✅ 自動儲存 API/MQTT 接收的 TSPL 資料
- ✅ 按日期和時間組織檔案結構
- ✅ 支援多種 TSPL 指令
- ✅ 完整的錯誤處理與驗證
- ✅ CORS 支援

## 技術棧

- **框架**: Gin Web Framework
- **MQTT**: Eclipse Paho MQTT Client
- **語言**: Go 1.21+

## 安裝與執行

### 1. 安裝依賴

```bash
cd backend
go mod download
```

### 2. 配置環境變數

複製 `.env.example` 為 `.env` 並根據需要修改配置:

```bash
cp .env.example .env
```

### 3. 啟動服務器

```bash
go run main.go
```

服務器將啟動於 `http://localhost:8080`

## API 端點

### 健康檢查

```
GET /api/health
```

回應:
```json
{
  "status": "ok",
  "mqtt": "connected"
}
```

### 渲染 TSPL

```
POST /api/render
```

請求:
```json
{
  "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"Hello TSPL!\"\nPRINT 1,1"
}
```

回應:
```json
{
  "success": true,
  "data": {
    "width": 800,
    "height": 400,
    "elements": [...],
    "labelSize": {...},
    "gap": {...}
  }
}
```

### 取得範例列表

```
GET /api/examples
```

### 取得範例詳情

```
GET /api/examples/:id
```

### 發布 MQTT 訊息

```
POST /api/mqtt/publish
```

請求:
```json
{
  "topic": "tspl/commands",
  "message": {
    "type": "render_request",
    "tspl_code": "...",
    "timestamp": 1234567890
  }
}
```

## MQTT 功能

### 訂閱主題

預設訂閱主題: `tspl/commands`

接收的訊息格式:
```json
{
  "type": "render_request",
  "tspl_code": "SIZE 100 mm, 50 mm...",
  "timestamp": 1234567890
}
```

### 發布主題

發布渲染結果到: `tspl/commands/result`

發布的訊息格式:
```json
{
  "type": "render_result",
  "timestamp": 1234567890,
  "data": {
    "width": 800,
    "height": 400,
    "elements": [...]
  }
}
```

## 支援的 TSPL 指令

- `SIZE` - 設定標籤尺寸
- `GAP` - 設定標籤間距
- `DIRECTION` - 設定列印方向
- `REFERENCE` - 設定參考點
- `CLS` - 清除緩衝區
- `TEXT` - 列印文字
- `BARCODE` - 列印條碼
- `QRCODE` - 列印 QR Code
- `BOX` - 繪製矩形
- `BAR` - 繪製實心線條
- `PRINT` - 執行列印

## 資料儲存功能

後端會自動儲存所有接收到的 TSPL 資料:

### 資料夾結構

```
data/
├── API_print/
│   ├── 2025_01_15/
│   │   ├── 10_30_45.tspl
│   │   ├── 11_20_30.tspl
│   │   └── ...
│   └── 2025_01_16/
│       └── ...
└── MQTT_print/
    ├── 2025_01_15/
    │   ├── 09_15_20.tspl
    │   └── ...
    └── 2025_01_16/
        └── ...
```

- **API_print/**: 儲存透過 API 接收的 TSPL 資料
- **MQTT_print/**: 儲存透過 MQTT 接收的 TSPL 資料
- 資料夾命名: `年_月_日` (例如: 2025_01_15)
- 檔案命名: `時_分_秒.tspl` (例如: 10_30_45.tspl)

## TSPL2 語法驗證

後端會在處理前驗證所有 TSPL 指令:

### 驗證規則

- ✅ 檢查命令是否有效
- ✅ 驗證參數數量和格式
- ✅ 檢查必要命令 (SIZE, PRINT)
- ✅ 驗證數值範圍
- ✅ 提供詳細的錯誤訊息

### 驗證錯誤回應範例

```json
{
  "success": false,
  "error": "TSPL 語法驗證失敗",
  "validation_errors": [
    {
      "line": 5,
      "command": "TEXT",
      "message": "TEXT 命令格式錯誤。正確格式: TEXT x,y,\"font\",rotation,x-scale,y-scale,\"content\""
    }
  ]
}
```

## 環境變數

| 變數 | 說明 | 預設值 |
|------|------|--------|
| `SERVER_PORT` | 服務器端口 | `8080` |
| `STORAGE_PATH` | 資料儲存路徑 | `./data` |
| `MQTT_BROKER` | MQTT Broker 地址 | `localhost` |
| `MQTT_PORT` | MQTT 端口 | `1883` |
| `MQTT_CLIENT_ID` | MQTT 客戶端 ID | `tspl-simulator` |
| `MQTT_USERNAME` | MQTT 用戶名 | - |
| `MQTT_PASSWORD` | MQTT 密碼 | - |
| `MQTT_TOPIC` | MQTT 主題 | `tspl/commands` |

## 專案結構

```
backend/
├── api/              # API 處理器和路由
│   ├── handlers.go   # 請求處理器
│   └── router.go     # 路由設定
├── config/           # 配置管理
│   └── config.go
├── models/           # 資料模型
│   └── models.go
├── mqtt/             # MQTT 客戶端
│   └── client.go
├── parser/           # TSPL 解析器
│   └── parser.go
├── storage/          # 檔案儲存服務
│   └── storage.go
├── validator/        # TSPL2 語法驗證器
│   └── validator.go
├── main.go           # 主程序入口
├── go.mod            # Go 模組定義
└── .env.example      # 環境變數範例
```

## 開發

### 執行測試

```bash
go test ./...
```

### 編譯

```bash
go build -o tspl-simulator
```

### 執行編譯後的程序

```bash
./tspl-simulator
```

## 注意事項

- MQTT 功能是可選的,如果不需要可以不配置 MQTT Broker
- 確保前端配置正確的後端 API 地址
- 支援 CORS,可以從不同域名訪問 API

## License

MIT License
