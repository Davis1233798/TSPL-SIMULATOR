package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"tspl-simulator/models"
	"tspl-simulator/mqtt"
	"tspl-simulator/parser"
	"tspl-simulator/storage"
	"tspl-simulator/validator"
)

var storageService *storage.StorageService

// InitStorage 初始化儲存服務
func InitStorage(basePath string) {
	storageService = storage.NewStorageService(basePath)
}

// RenderHandler 處理 TSPL 渲染請求
func RenderHandler(c *gin.Context) {
	var req models.RenderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.RenderResponse{
			Success: false,
			Error:   "請求格式錯誤: " + err.Error(),
		})
		return
	}

	// 驗證 TSPL 語法
	validationResult := validator.ValidateTSPL(req.TSPLCode)
	if !validationResult.Valid {
		// 轉換驗證錯誤格式
		var modelErrors []models.ValidationError
		for _, err := range validationResult.Errors {
			modelErrors = append(modelErrors, models.ValidationError{
				Line:    err.Line,
				Command: err.Command,
				Message: err.Message,
			})
		}

		c.JSON(http.StatusBadRequest, models.RenderResponse{
			Success:          false,
			Error:            "TSPL 語法驗證失敗",
			ValidationErrors: modelErrors,
		})
		return
	}

	// 儲存 API 接收的資料
	if storageService != nil {
		if filePath, err := storageService.SaveAPIData(req.TSPLCode); err != nil {
			log.Printf("儲存 API 資料失敗: %v", err)
		} else {
			log.Printf("API 資料已儲存至: %s", filePath)
		}
	}

	// 解析 TSPL
	renderData, err := parser.ParseTSPL(req.TSPLCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.RenderResponse{
			Success: false,
			Error:   "TSPL 解析錯誤: " + err.Error(),
		})
		return
	}

	// 如果 MQTT 已連接,也發布到 MQTT
	mqttClient := mqtt.GetClient()
	if mqttClient != nil && mqttClient.IsConnected() {
		if err := mqttClient.PublishRenderResult(renderData); err != nil {
			log.Printf("發布到 MQTT 失敗: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.RenderResponse{
		Success: true,
		Data:    renderData,
	})
}

// GetExamplesHandler 取得範例列表
func GetExamplesHandler(c *gin.Context) {
	examples := []models.ExampleInfo{
		{
			ID:          "basic_text",
			Name:        "基本文字",
			Description: "簡單的文字標籤範例",
			Category:    "basic",
		},
		{
			ID:          "barcode",
			Name:        "條碼",
			Description: "Code 128 條碼範例",
			Category:    "barcode",
		},
		{
			ID:          "qrcode",
			Name:        "QR Code",
			Description: "QR Code 標籤範例",
			Category:    "qrcode",
		},
		{
			ID:          "product_label",
			Name:        "產品標籤",
			Description: "零售產品標籤範例",
			Category:    "retail",
		},
		{
			ID:          "shipping_label",
			Name:        "運輸標籤",
			Description: "物流配送標籤範例",
			Category:    "logistics",
		},
		{
			ID:          "inventory_label",
			Name:        "庫存標籤",
			Description: "倉庫管理標籤範例",
			Category:    "warehouse",
		},
		{
			ID:          "name_badge",
			Name:        "名牌",
			Description: "活動訪客證範例",
			Category:    "event",
		},
		{
			ID:          "asset_tag",
			Name:        "資產標籤",
			Description: "公司財產標籤範例",
			Category:    "asset",
		},
		{
			ID:          "price_tag",
			Name:        "價格標籤",
			Description: "商店貨架標籤範例",
			Category:    "retail",
		},
		{
			ID:          "food_label",
			Name:        "食品標籤",
			Description: "生鮮產品標籤範例",
			Category:    "food",
		},
	}

	c.JSON(http.StatusOK, models.ExamplesResponse{
		Success:  true,
		Examples: examples,
	})
}

// GetExampleDetailHandler 取得範例詳情
func GetExampleDetailHandler(c *gin.Context) {
	exampleID := c.Param("id")

	// 讀取範例檔案
	examplePath := filepath.Join("examples", exampleID+".tspl")
	content, err := ioutil.ReadFile(examplePath)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ExampleDetailResponse{
			Success: false,
			Error:   "範例不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.ExampleDetailResponse{
		Success: true,
		Code:    string(content),
	})
}

// HealthCheckHandler 健康檢查
func HealthCheckHandler(c *gin.Context) {
	mqttStatus := "disconnected"
	mqttClient := mqtt.GetClient()
	if mqttClient != nil && mqttClient.IsConnected() {
		mqttStatus = "connected"
	}

	c.JSON(http.StatusOK, models.HealthResponse{
		Status: "ok",
		MQTT:   mqttStatus,
	})
}

// MQTTPublishHandler 發布訊息到 MQTT
func MQTTPublishHandler(c *gin.Context) {
	var req struct {
		Topic   string      `json:"topic" binding:"required"`
		Message interface{} `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "請求格式錯誤: " + err.Error(),
		})
		return
	}

	mqttClient := mqtt.GetClient()
	if mqttClient == nil || !mqttClient.IsConnected() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "MQTT 未連接",
		})
		return
	}

	if err := mqttClient.Publish(req.Topic, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "發布失敗: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "訊息已發布",
	})
}
