# TSPL Simulator - 專案總覽

## 專案簡介

TSPL Simulator 是一個完整的 TSPL2 (TSC Printer Language 2) 標籤列印模擬器,提供網頁界面讓使用者可以在瀏覽器中編輯、驗證和預覽 TSPL 標籤程式碼,無需實體印表機。

## 核心功能

### 1. TSPL2 語法驗證 ✅
- 即時檢查 TSPL 命令語法
- 驗證參數格式和數值範圍
- 檢查必要命令 (SIZE, PRINT)
- 提供詳細的錯誤訊息和行號

### 2. 視覺化渲染 ✅
- 即時預覽標籤外觀
- 支援文字、條碼、QR Code
- 支援圖形元素 (矩形、線條)
- 精確的尺寸計算 (基於 203 DPI)

### 3. 自動檔案儲存 ✅
- API 接收的資料儲存至 `data/API_print/`
- MQTT 接收的資料儲存至 `data/MQTT_print/`
- 按日期組織: `年_月_日/時_分_秒.tspl`
- 自動建立資料夾結構

### 4. MQTT 整合 ✅
- 訂閱 MQTT 主題接收列印請求
- 發布渲染結果到 MQTT
- 支援認證和自動重連
- 完整的錯誤處理

### 5. RESTful API ✅
- `/api/health` - 健康檢查
- `/api/render` - TSPL 渲染
- `/api/examples` - 範例列表
- `/api/mqtt/publish` - MQTT 發布

## 技術架構

```
┌─────────────────────────────────────────────────────────────┐
│                         使用者介面                             │
│                    (React + TypeScript)                      │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  程式碼編輯器  │  │  即時預覽      │  │  範例庫        │      │
│  │  (Monaco)    │  │  (Canvas)     │  │  (10+ 範例)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌────────────────────────────────────────────────────┐     │
│  │            前端 TSPL 語法驗證器 (Fallback)           │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP / WebSocket
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      後端服務 (Go)                            │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  API Server  │  │ TSPL Validator│  │ File Storage │      │
│  │  (Gin)       │  │               │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐                         │
│  │ TSPL Parser  │  │ MQTT Client  │                         │
│  │              │  │ (Paho)       │                         │
│  └──────────────┘  └──────────────┘                         │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ MQTT
                            ▼
                   ┌─────────────────┐
                   │  MQTT Broker    │
                   │  (Mosquitto)    │
                   └─────────────────┘
```

## 專案結構

```
TSPL-simulator/
├── frontend/                  # React 前端
│   ├── src/
│   │   ├── components/        # React 元件
│   │   │   ├── Editor/        # 程式碼編輯器
│   │   │   ├── Preview/       # 預覽畫布
│   │   │   └── Examples/      # 範例選擇器
│   │   ├── services/          # API 服務
│   │   │   ├── tsplApi.ts     # API 呼叫
│   │   │   ├── mockApi.ts     # 前端 fallback
│   │   │   └── validator.ts   # 前端驗證器
│   │   ├── types/             # TypeScript 類型
│   │   └── utils/             # 工具函數
│   ├── public/
│   └── package.json
│
├── backend/                   # Go 後端
│   ├── api/                   # API 處理器
│   │   ├── handlers.go        # 請求處理
│   │   └── router.go          # 路由設定
│   ├── config/                # 配置管理
│   │   └── config.go
│   ├── models/                # 資料模型
│   │   └── models.go
│   ├── parser/                # TSPL 解析器
│   │   └── parser.go
│   ├── validator/             # TSPL 驗證器 (新增)
│   │   └── validator.go
│   ├── storage/               # 檔案儲存 (新增)
│   │   └── storage.go
│   ├── mqtt/                  # MQTT 客戶端
│   │   └── client.go
│   ├── main.go
│   ├── go.mod
│   ├── .env.example
│   ├── README.md
│   ├── BUILD.md               # 建置指南 (新增)
│   └── QUICKSTART.md          # 快速開始 (新增)
│
├── data/                      # 自動儲存的資料 (新增)
│   ├── API_print/
│   │   └── YYYY_MM_DD/
│   │       └── HH_MM_SS.tspl
│   └── MQTT_print/
│       └── YYYY_MM_DD/
│           └── HH_MM_SS.tspl
│
└── PROJECT_OVERVIEW.md        # 本檔案
```

## 支援的 TSPL 指令

### 基礎設定
- `SIZE width, height` - 設定標籤尺寸
- `GAP distance, offset` - 設定標籤間距
- `DIRECTION n` - 設定列印方向 (0-3)
- `REFERENCE x, y` - 設定參考點
- `DENSITY n` - 設定列印濃度 (0-15)
- `SPEED n` - 設定列印速度 (1-14)

### 繪圖指令
- `TEXT x,y,"font",rotation,x-scale,y-scale,"content"` - 列印文字
- `BARCODE x,y,"type",height,readable,rotation,narrow,wide,"code"` - 列印條碼
- `QRCODE x,y,ecc,size,mode,rotation,"data"` - 列印 QR Code
- `BOX x,y,x_end,y_end,thickness` - 繪製矩形框
- `BAR x,y,width,height` - 繪製實心矩形

### 控制指令
- `CLS` - 清除緩衝區
- `PRINT m,n` - 執行列印

## 資料流程

### API 請求流程

```
1. 前端發送 TSPL 程式碼
   ↓
2. 後端接收請求
   ↓
3. TSPL 語法驗證
   ├─ 驗證失敗 → 回傳錯誤訊息 (含行號和錯誤詳情)
   └─ 驗證成功 ↓
4. 儲存到 data/API_print/YYYY_MM_DD/HH_MM_SS.tspl
   ↓
5. 解析 TSPL 生成渲染資料
   ↓
6. 回傳渲染資料給前端
   ↓
7. 前端在 Canvas 上渲染標籤
```

### MQTT 訊息流程

```
1. 外部系統發送 MQTT 訊息到 tspl/commands
   {
     "type": "render_request",
     "tspl_code": "...",
     "timestamp": 1234567890
   }
   ↓
2. 後端 MQTT 客戶端接收
   ↓
3. TSPL 語法驗證
   ├─ 驗證失敗 → 記錄錯誤日誌
   └─ 驗證成功 ↓
4. 儲存到 data/MQTT_print/YYYY_MM_DD/HH_MM_SS.tspl
   ↓
5. 解析 TSPL 生成渲染資料
   ↓
6. 發布結果到 tspl/commands/result
   {
     "type": "render_result",
     "timestamp": 1234567890,
     "data": { ... }
   }
```

## 檔案儲存機制

### 命名規則

- **資料夾**: `年_月_日` (例如: `2025_01_15`)
- **檔案**: `時_分_秒.tspl` (例如: `10_30_45.tspl`)
- 如果同一秒內有多個請求,會加上毫秒後綴

### 範例結構

```
data/
├── API_print/
│   ├── 2025_01_15/
│   │   ├── 09_30_15.tspl       # 上午 9:30:15 的請求
│   │   ├── 09_30_16.tspl       # 上午 9:30:16 的請求
│   │   ├── 14_25_30.tspl       # 下午 2:25:30 的請求
│   │   └── 14_25_30_123.tspl   # 同一秒的另一個請求 (含毫秒)
│   └── 2025_01_16/
│       └── 10_00_00.tspl
└── MQTT_print/
    ├── 2025_01_15/
    │   ├── 08_15_00.tspl
    │   └── 16_45_30.tspl
    └── 2025_01_16/
        └── 11_20_00.tspl
```

## 語法驗證功能

### 驗證檢查項目

1. **命令有效性** - 檢查命令是否為有效的 TSPL 指令
2. **參數數量** - 驗證每個命令的參數數量是否正確
3. **參數格式** - 檢查參數格式 (數字、字串、單位等)
4. **數值範圍** - 驗證數值是否在允許範圍內
5. **必要命令** - 確保包含 SIZE 和 PRINT 命令
6. **單位一致性** - 檢查 SIZE 和 GAP 的單位是否一致

### 驗證錯誤範例

```json
{
  "success": false,
  "error": "TSPL 語法驗證失敗",
  "validation_errors": [
    {
      "line": 2,
      "command": "DIRECTION",
      "message": "方向必須在 0-3 之間"
    },
    {
      "line": 5,
      "command": "TEXT",
      "message": "TEXT 命令格式錯誤。正確格式: TEXT x,y,\"font\",rotation,x-scale,y-scale,\"content\""
    },
    {
      "line": 0,
      "command": "SIZE",
      "message": "缺少必要的 SIZE 命令"
    }
  ]
}
```

## 開發環境設定

### 前端

```bash
cd frontend
npm install
npm start
```

瀏覽器開啟 http://localhost:3000

### 後端

```bash
cd backend
go mod download
go run main.go
```

後端啟動於 http://localhost:8080

### 環境變數

**前端 (.env)**:
```env
REACT_APP_API_URL=http://localhost:8080/api
```

**後端 (.env)**:
```env
SERVER_PORT=8080
STORAGE_PATH=./data
MQTT_BROKER=localhost
MQTT_PORT=1883
MQTT_CLIENT_ID=tspl-simulator
MQTT_TOPIC=tspl/commands
```

## 生產環境部署

### Docker 部署

```bash
# 建置前端
cd frontend
npm run build

# 建置後端
cd ../backend
docker build -t tspl-simulator .

# 執行
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e SERVER_PORT=8080 \
  tspl-simulator
```

### 傳統部署

**前端**:
```bash
npm run build
# 將 build/ 資料夾部署到 Web 伺服器 (Nginx, Apache)
```

**後端**:
```bash
go build -o tspl-simulator
./tspl-simulator
```

## 測試

### API 測試

```bash
# 健康檢查
curl http://localhost:8080/api/health

# 渲染測試
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"Hello!\"\nPRINT 1,1"
  }'

# 語法錯誤測試
curl -X POST http://localhost:8080/api/render \
  -H "Content-Type: application/json" \
  -d '{
    "tspl_code": "TEXT 100,100,\"3\",0,1,1,\"No SIZE!\"\nPRINT 1,1"
  }'
```

### MQTT 測試

```bash
# 訂閱結果
mosquitto_sub -t "tspl/commands/result"

# 發布請求
mosquitto_pub -t "tspl/commands" -m '{
  "type": "render_request",
  "tspl_code": "SIZE 100 mm, 50 mm\nGAP 3 mm, 0 mm\nCLS\nTEXT 100,100,\"3\",0,1,1,\"MQTT Test\"\nPRINT 1,1",
  "timestamp": 1234567890
}'
```

## 效能考量

- **前端**: React 虛擬 DOM 優化渲染效能
- **後端**: Go 協程處理併發請求
- **儲存**: 非同步寫入,不阻塞主執行緒
- **MQTT**: 自動重連和訊息緩衝

## 安全性

- **CORS**: 限制允許的來源
- **輸入驗證**: 嚴格的 TSPL 語法檢查
- **檔案儲存**: 安全的檔案命名,防止路徑遍歷
- **MQTT**: 支援用戶名/密碼認證

## 授權

MIT License

## 貢獻

歡迎提交 Issue 和 Pull Request!

## 未來規劃

- [ ] 支援更多 TSPL 指令
- [ ] 批次處理多個標籤
- [ ] 匯出為 PDF/PNG
- [ ] 印表機實體連接
- [ ] 範本管理系統
- [ ] 使用者認證
- [ ] 資料庫整合
- [ ] WebSocket 即時同步
