# 運行 TSPL Simulator 專案指南

## 前置條件

確保已安裝:
- Go 1.21+ (執行 `go version` 確認)
- Node.js 16+ 和 npm (執行 `node --version` 確認)

## 🔧 已修復的問題

### ✅ 類型轉換錯誤已修復

**之前的錯誤**:
```
api\handlers.go:42:22: cannot use validationResult.Errors (variable of type []"tspl-simulator/validator".ValidationError) as []models.ValidationError value in struct literal
```

**修復方法**: 在 [backend/api/handlers.go](backend/api/handlers.go) 中添加了類型轉換:
```go
// 將 validator 錯誤轉換為 models 錯誤
var modelErrors []models.ValidationError
for _, err := range validationResult.Errors {
    modelErrors = append(modelErrors, models.ValidationError{
        Line:    err.Line,
        Command: err.Command,
        Message: err.Message,
    })
}
```

## 🚀 啟動後端

### 步驟 1: 設置環境變數 (可選)

在 `backend` 資料夾中創建 `.env` 檔案:
```env
# 伺服器設定
PORT=8080

# MQTT 設定 (如果需要)
MQTT_BROKER=tcp://localhost:1883
MQTT_CLIENT_ID=tspl-simulator
MQTT_USERNAME=
MQTT_PASSWORD=

# 儲存路徑
STORAGE_PATH=./data
```

或者使用預設值(不需要 .env 檔案)

### 步驟 2: 編譯並運行

打開 PowerShell 或命令提示字元:

```powershell
# 切換到後端目錄
cd backend

# 如果之前遇到網路問題,設置 GOPROXY (可選)
$env:GOPROXY = "https://goproxy.cn,direct"

# 下載依賴 (如果尚未執行)
go mod tidy

# 運行後端
go run main.go
```

### 預期輸出

```
2025/01/15 21:00:00 儲存服務已初始化: ./data
2025/01/15 21:00:00 MQTT 未配置,跳過初始化
2025/01/15 21:00:00 伺服器運行於 :8080
2025/01/15 21:00:00 API 資料儲存路徑: ./data/API_print
2025/01/15 21:00:00 MQTT 資料儲存路徑: ./data/MQTT_print
```

### 編譯成可執行檔 (可選)

```powershell
# 編譯
go build -o tspl-simulator.exe .

# 運行
./tspl-simulator.exe
```

## 🎨 啟動前端

### 步驟 1: 設置環境變數

在 `frontend` 資料夾中創建 `.env` 檔案:
```env
REACT_APP_API_URL=http://localhost:8080/api
```

### 步驟 2: 安裝依賴並運行

打開新的 PowerShell 或命令提示字元:

```powershell
# 切換到前端目錄
cd frontend

# 安裝依賴 (首次運行)
npm install

# 啟動開發伺服器
npm start
```

### 預期結果

- 瀏覽器自動打開 http://localhost:3000
- 前端介面載入完成
- 左側顯示 TSPL 編輯器
- 右側顯示標籤預覽區域

## 🧪 測試完整流程

### 測試 1: 正確的 TSPL 語法

1. 在編輯器中輸入:
```tspl
SIZE 100 mm, 50 mm
GAP 3 mm, 0 mm
CLS
TEXT 100,100,"3",0,1,1,"Hello TSPL!"
PRINT 1,1
```

2. 點擊「預覽」按鈕

**預期結果**:
- ✅ 前端語法檢查: 無錯誤
- ✅ 後端驗證: 通過
- ✅ 成功渲染標籤預覽
- ✅ 檔案儲存至: `backend/data/API_print/2025_01_15/21_05_30.tspl`

### 測試 2: 語法錯誤 (缺少 SIZE)

1. 在編輯器中輸入:
```tspl
GAP 3 mm, 0 mm
CLS
TEXT 100,100,"3",0,1,1,"Test"
PRINT 1,1
```

2. 點擊「預覽」按鈕

**預期結果**:
- ⚠️ 前端語法檢查: 警告 "建議使用 SIZE 指令"
- ❌ 後端驗證: 失敗
- ❌ 顯示 ValidationErrors 元件:
```
後端驗證錯誤
1 個錯誤

❌ 行 0
   命令: SIZE
   缺少必要的 SIZE 命令

💡 提示: 請修正上述錯誤後再次嘗試渲染
```

### 測試 3: 參數錯誤 (DIRECTION 超出範圍)

1. 在編輯器中輸入:
```tspl
SIZE 100 mm, 50 mm
DIRECTION 99
CLS
PRINT 1,1
```

2. 點擊「預覽」按鈕

**預期結果**:
- ❌ 後端驗證: 失敗
- ❌ 顯示錯誤: "行 2 [DIRECTION]: 方向必須在 0-3 之間"

### 測試 4: TEXT 格式錯誤

1. 在編輯器中輸入:
```tspl
SIZE 100 mm, 50 mm
CLS
TEXT 100,100,3,0,1,1,Hello
PRINT 1,1
```

2. 點擊「預覽」按鈕

**預期結果**:
- ❌ 前端語法檢查: 錯誤 "TEXT 指令格式錯誤"
- ❌ 後端驗證: 失敗
- ❌ 顯示錯誤: "TEXT 命令格式錯誤。正確格式: TEXT x,y,\"font\",rotation,x-scale,y-scale,\"content\""

## 📁 檢查儲存的檔案

### API 請求儲存路徑

```
backend/data/API_print/
└── 2025_01_15/          # 今天的日期
    ├── 10_30_45.tspl    # 10:30:45 提交的請求
    ├── 10_31_20.tspl    # 10:31:20 提交的請求
    └── 14_25_10.tspl    # 14:25:10 提交的請求
```

### 查看儲存的檔案

```powershell
# 列出今天的檔案
ls backend/data/API_print/2025_01_15/

# 查看特定檔案內容
cat backend/data/API_print/2025_01_15/10_30_45.tspl
```

## 🔍 故障排除

### 問題 1: Go 命令找不到

**錯誤**: `'go' 不是內部或外部命令`

**解決方案**:
1. 確認 Go 已安裝: 下載自 https://go.dev/dl/
2. 檢查環境變數: Go 安裝路徑應在 PATH 中
3. 重啟 PowerShell 或命令提示字元

### 問題 2: 依賴下載失敗

**錯誤**: `dial tcp: lookup proxy.golang.org: no such host`

**解決方案**:
```powershell
$env:GOPROXY = "https://goproxy.cn,direct"
go mod tidy
```

### 問題 3: 編譯錯誤

**錯誤**: 各種編譯錯誤

**解決方案**:
1. 確保所有文件都已保存
2. 檢查 `go.mod` 中的模組名稱是否為 `tspl-simulator`
3. 重新運行 `go mod tidy`
4. 清理並重建: `go clean && go build`

### 問題 4: 前端無法連接後端

**錯誤**: 前端顯示 "後端不可用"

**解決方案**:
1. 確認後端正在運行 (`go run main.go`)
2. 檢查 `frontend/.env` 中的 API URL 是否正確
3. 檢查後端端口是否為 8080
4. 檢查防火牆設置

### 問題 5: CORS 錯誤

**錯誤**: 瀏覽器控制台顯示 CORS 錯誤

**解決方案**:
後端已配置 CORS 中間件,應該不會有此問題。如果出現:
1. 確認前端 URL 是 `http://localhost:3000`
2. 檢查 [backend/main.go](backend/main.go:44-50) 中的 CORS 配置

## 📊 功能驗證清單

運行專案後,請確認以下功能:

- [ ] 後端成功啟動在端口 8080
- [ ] 前端成功啟動在端口 3000
- [ ] 前端可以載入範例 TSPL 程式碼
- [ ] 前端語法檢查器正常運作
- [ ] 提交正確的 TSPL 可以成功渲染
- [ ] 提交錯誤的 TSPL 會顯示 ValidationErrors 元件
- [ ] 檔案成功儲存到 `backend/data/API_print/年_月_日/時_分_秒.tspl`
- [ ] 後端控制台顯示儲存日誌
- [ ] 修改程式碼後錯誤自動清除

## 🎯 下一步

1. **測試 MQTT 功能** (需要 MQTT Broker):
   - 安裝 Mosquitto 或其他 MQTT Broker
   - 配置 `.env` 中的 MQTT 設定
   - 發送 MQTT 訊息測試儲存功能

2. **添加更多範例**:
   - 在前端添加更多 TSPL 範例
   - 測試各種 TSPL 命令

3. **改進 UI/UX**:
   - 添加語法高亮
   - 實現行號點擊跳轉
   - 添加自動修正建議

4. **部署到生產環境**:
   - 編譯前後端
   - 設置反向代理 (Nginx)
   - 配置 HTTPS

## 📚 相關文件

- [後端實現報告](BACKEND_IMPLEMENTATION.md)
- [前端實現報告](FRONTEND_IMPLEMENTATION.md)
- [專案概覽](PROJECT_OVERVIEW.md)
- [後端快速開始](backend/QUICKSTART.md)
- [後端構建指南](backend/BUILD.md)

---

**注意**: 所有類型轉換錯誤已修復,專案現在應該可以正常編譯和運行。如果遇到任何問題,請參考故障排除部分或查看相關文件。
