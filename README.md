# TSPL Simulator

一個功能完整的 TSPL (TSC Printer Language) 標籤列印模擬器,支援即時語法檢查、標籤預覽和自動檔案儲存。

[English](#english) | [中文](#chinese)

---

<a name="english"></a>
## English

A full-stack TSPL (TSC Printer Language) label simulator and preview tool with backend validation and file storage.

**✨ Features both frontend and backend validation with automatic file storage!**

### Features

- ✅ **Dual Validation**: Frontend real-time checking + Backend strict validation
- 📝 Online TSPL editor with syntax validation
- 🔍 Intelligent error reporting with line numbers and suggestions
- 👁️ Live label preview with Canvas rendering
- 💾 **Automatic file storage** - API and MQTT requests saved with date/time organization
- 🎨 Support for text, barcodes, QR codes, and graphics (30+ TSPL commands)
- 📱 Responsive web interface
- 🚀 **Ready for production** - Backend with Go + Frontend with React
- 📦 10+ built-in examples

### Tech Stack

- **Backend**: Go 1.21+ with Gin framework
- **Frontend**: React 18 + TypeScript
- **Rendering**: HTML5 Canvas
- **MQTT**: Eclipse Paho (optional)
- **Architecture**: Full-stack with dual validation

### Quick Start

#### Requirements

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Node.js 16+** and npm - [Download](https://nodejs.org/)

#### Running the Application

**Step 1: Start Backend** (PowerShell)
```powershell
cd backend
go run main.go
```

Expected output:
```
儲存服務已初始化: ./data
伺服器運行於 :8080
API 資料儲存路徑: ./data/API_print
```

**Step 2: Start Frontend** (New PowerShell window)
```powershell
cd frontend
npm install  # First time only
npm start
```

The application will open at http://localhost:3000

**👉 See [QUICK_START.md](QUICK_START.md) for detailed 30-second guide!**

#### Build for Production

**Backend**:
```powershell
cd backend
go build -o tspl-simulator.exe .
```

**Frontend**:
```powershell
cd frontend
npm run build
```

Build files will be in `frontend/build/` directory and backend executable `tspl-simulator.exe` is ready for deployment.

### Supported TSPL Commands (30+)

**Basic Commands**:
- **SIZE**, **GAP**, **CLS**, **PRINT** - Label setup and printing
- **DIRECTION** (0-3) - Print direction with validation

**Text Commands**:
- **TEXT** - Print text with font, rotation, and scaling

**Barcode Commands**:
- **BARCODE** - 1D barcodes (Code 128, Code 39, EAN13, etc.)
- **QRCODE** - QR codes with error correction levels

**Graphics Commands**:
- **BOX**, **BAR**, **BITMAP** - Rectangles, lines, and images

**Settings Commands**:
- **DENSITY** (0-15), **SPEED** (1-14) - Print quality and speed with validation
- **OFFSET**, **REFERENCE**, **SHIFT** - Position adjustments

**Backend validates all parameter ranges and formats!**

For detailed command reference, see [BACKEND_IMPLEMENTATION.md](BACKEND_IMPLEMENTATION.md)

### Built-in Examples

The application includes 10 practical examples:

1. **Basic Text** - Simple text label
2. **Barcode** - Code 128 barcode
3. **QR Code** - QR code label
4. **Product Label** - Retail product tag
5. **Shipping Label** - Logistics shipping label
6. **Inventory Label** - Warehouse management
7. **Name Badge** - Event visitor badge
8. **Asset Tag** - Company property tag
9. **Price Tag** - Store shelf label
10. **Food Label** - Fresh product label

All examples are available in the `examples/` directory.

### Usage Example

```tspl
SIZE 100 mm, 50 mm
GAP 3 mm, 0 mm
CLS
TEXT 100,100,"3",0,1,1,"Hello TSPL!"
BARCODE 100,200,"128",100,1,0,2,2,"123456789"
QRCODE 400,200,H,5,A,0,"https://example.com"
PRINT 1,1
```

### Automatic File Storage 💾

All validated TSPL submissions are automatically saved:

**File Structure**:
```
backend/data/
├── API_print/
│   └── 2025_11_15/              ← Year_Month_Day
│       ├── 21_30_45.tspl        ← Hour_Minute_Second
│       ├── 21_31_20.tspl
│       └── 21_35_00.tspl
└── MQTT_print/
    └── 2025_11_15/
        └── 22_15_30.tspl
```

**Only validation-passed requests are saved!** ✅

### Documentation 📚

Complete documentation is available:

- **[QUICK_START.md](QUICK_START.md)** - 30-second quick start guide
- **[TESTING_GUIDE.md](TESTING_GUIDE.md)** - 7 comprehensive test cases
- **[RUNNING_THE_PROJECT.md](RUNNING_THE_PROJECT.md)** - Detailed setup and troubleshooting
- **[BACKEND_IMPLEMENTATION.md](BACKEND_IMPLEMENTATION.md)** - Backend technical details
- **[FRONTEND_IMPLEMENTATION.md](FRONTEND_IMPLEMENTATION.md)** - Frontend technical details
- **[PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md)** - Full project architecture

### Deployment

#### Backend
Deploy the Go binary to any VPS or cloud platform:
```powershell
go build -o tspl-simulator.exe .
# Upload and run on server
```

#### Frontend

**Vercel** (Recommended):
```bash
cd frontend
vercel --prod
```

**Netlify**: Drag and drop the `frontend/build` folder

**GitHub Pages**: Add deployment script to `frontend/package.json`

### Browser Support

- Chrome (Recommended)
- Firefox
- Safari
- Edge

Requires a modern browser with HTML5 Canvas support.

### License

MIT License

### Contributing

Issues and Pull Requests are welcome!

---

<a name="chinese"></a>
## 中文

一個功能完整的 TSPL (TSC Printer Language) 標籤列印模擬器,支援即時語法檢查、標籤預覽和自動檔案儲存。

**✨ 前後端雙重驗證 + 自動檔案儲存!**

### 功能特色

- ✅ **雙重驗證**: 前端即時檢查 + 後端嚴格驗證
- 📝 線上 TSPL 編輯器,支援語法驗證
- 🔍 智能錯誤報告,包含行號和修正建議
- 👁️ 即時標籤預覽 (Canvas 渲染)
- 💾 **自動檔案儲存** - API 和 MQTT 請求按日期/時間組織
- 🎨 支援文字、條碼、QR Code 和圖形 (30+ TSPL 命令)
- 📱 響應式網頁介面
- 🚀 **生產就緒** - Go 後端 + React 前端
- 📦 10+ 內建範例

### 技術棧

- **後端**: Go 1.21+ with Gin 框架
- **前端**: React 18 + TypeScript
- **渲染**: HTML5 Canvas
- **MQTT**: Eclipse Paho (可選)
- **架構**: 全端雙重驗證

### 快速開始

#### 環境需求

- **Go 1.21+** - [下載](https://go.dev/dl/)
- **Node.js 16+** 和 npm - [下載](https://nodejs.org/)

#### 運行應用

**步驟 1: 啟動後端** (PowerShell)
```powershell
cd backend
go run main.go
```

預期輸出:
```
儲存服務已初始化: ./data
伺服器運行於 :8080
API 資料儲存路徑: ./data/API_print
```

**步驟 2: 啟動前端** (新 PowerShell 視窗)
```powershell
cd frontend
npm install  # 僅首次需要
npm start
```

應用將在 http://localhost:3000 啟動

**👉 詳見 [QUICK_START.md](QUICK_START.md) 查看 30 秒快速指南!**

#### 建置生產版本

```bash
cd frontend
npm run build
```

建置檔案將在 `frontend/build/` 目錄中,可部署到任何靜態網站託管服務。

### 支援的 TSPL 指令

- **SIZE** - 設定標籤尺寸
- **GAP** - 設定標籤間距
- **DIRECTION** - 設定列印方向
- **CLS** - 清除緩衝區
- **TEXT** - 列印文字
- **BARCODE** - 列印條碼 (Code 128, Code 39, EAN13 等)
- **QRCODE** - 列印 QR Code
- **BOX** - 繪製矩形
- **BAR** - 繪製實心線條
- **PRINT** - 執行列印

詳細指令說明請參考 [docs/TSPL_COMMANDS.md](./docs/TSPL_COMMANDS.md)

### 內建範例

應用包含 10 個實用範例:

1. **基本文字** - 簡單文字標籤
2. **條碼** - Code 128 條碼
3. **QR Code** - QR Code 標籤
4. **產品標籤** - 零售商品標籤
5. **運輸標籤** - 物流配送標籤
6. **庫存標籤** - 倉庫管理標籤
7. **名牌** - 活動訪客證
8. **資產標籤** - 公司財產標籤
9. **價格標籤** - 商店貨架標籤
10. **食品標籤** - 生鮮產品標籤

所有範例都在 `examples/` 目錄中。

### 使用範例

```tspl
SIZE 100 mm, 50 mm
GAP 3 mm, 0 mm
CLS
TEXT 100,100,"3",0,1,1,"Hello TSPL!"
BARCODE 100,200,"128",100,1,0,2,2,"123456789"
QRCODE 400,200,H,5,A,0,"https://example.com"
PRINT 1,1
```

### 部署

#### Vercel (推薦)

```bash
npm i -g vercel
cd frontend
vercel --prod
```

#### Netlify

直接拖放 `frontend/build` 資料夾到 Netlify。

#### GitHub Pages

在 `frontend/package.json` 添加部署腳本後執行 `npm run deploy`。

### 瀏覽器支援

- Chrome (推薦)
- Firefox
- Safari
- Edge

需要支援 HTML5 Canvas 的現代瀏覽器。

### 授權

MIT License

### 貢獻

歡迎提交 Issues 和 Pull Requests!

---

**Start now! Visit http://localhost:3000 after running `npm start` 🚀**

**現在就開始! 執行 `npm start` 後訪問 http://localhost:3000 🚀**
