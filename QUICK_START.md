# 🚀 TSPL Simulator 快速啟動

## 30 秒快速開始

### 第 1 步: 啟動前端 (已完成 ✅)

你的前端已經在運行! 瀏覽器應該會自動打開 http://localhost:3000

### 第 2 步: 啟動後端 ⚡

**打開新的 PowerShell 視窗**,執行:

```powershell
cd C:\Users\solidityDeveloper\TSPL-simlator\backend
go run main.go
```

**預期看到**:
```
2025/11/15 21:20:00 儲存服務已初始化: ./data
2025/11/15 21:20:00 MQTT 未配置,跳過初始化
2025/11/15 21:20:00 伺服器運行於 :8080
```

✅ 看到這些訊息 = 後端成功啟動!

### 第 3 步: 測試功能 🧪

1. 打開瀏覽器 http://localhost:3000
2. 在編輯器中輸入:
```tspl
SIZE 100 mm, 50 mm
GAP 3 mm, 0 mm
CLS
TEXT 100,100,"3",0,1,1,"Hello World!"
PRINT 1,1
```
3. 點擊「預覽」

**應該看到**:
- ✅ 右側顯示渲染的標籤
- ✅ 後端控制台顯示: "API 資料已儲存至: ..."

---

## 測試驗證錯誤

試試這個**錯誤的** TSPL:

```tspl
DIRECTION 99
TEXT 100,100,3,0,1,1,Error
```

**應該看到**:
- ❌ 紅色的 ValidationErrors 元件
- ❌ 列出所有錯誤 (缺少 SIZE, DIRECTION 超出範圍, TEXT 格式錯誤)

---

## 📁 查看儲存的檔案

```powershell
# 列出今天的檔案
ls backend\data\API_print\2025_11_15\

# 讀取檔案內容
cat backend\data\API_print\2025_11_15\21_30_45.tspl
```

---

## 🔧 如果遇到問題

### 前端顯示 "後端不可用"
→ 確認後端正在運行 (`go run main.go`)

### 點擊預覽無反應
→ 按 F12 打開瀏覽器開發者工具,查看 Console 錯誤

### Go 命令找不到
→ 確認已安裝 Go: https://go.dev/dl/

---

## 📚 詳細文件

- [完整測試指南](TESTING_GUIDE.md) - 7 個測試案例
- [運行指南](RUNNING_THE_PROJECT.md) - 詳細啟動步驟
- [後端實現報告](BACKEND_IMPLEMENTATION.md) - 後端技術細節
- [前端實現報告](FRONTEND_IMPLEMENTATION.md) - 前端技術細節

---

## ✅ 功能清單

- [x] 前端即時語法檢查
- [x] 後端嚴格驗證
- [x] 美化的錯誤顯示
- [x] 自動檔案儲存 (日期/時間組織)
- [x] 雙重驗證機制
- [x] 範例選擇器
- [x] 標籤預覽
- [x] CORS 支援

---

**就這麼簡單!** 🎉

現在你有一個完整的 TSPL 模擬器,具備:
- 📝 即時語法檢查
- 🎨 標籤預覽
- ❌ 詳細錯誤報告
- 💾 自動檔案儲存

開始測試吧! 🚀
