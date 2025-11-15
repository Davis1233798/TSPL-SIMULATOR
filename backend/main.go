package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"tspl-simulator/api"
	"tspl-simulator/config"
	"tspl-simulator/mqtt"
	"tspl-simulator/storage"
)

func main() {
	// 載入配置
	cfg := config.LoadConfig()

	// 初始化儲存服務
	storagePath := getEnv("STORAGE_PATH", "./data")
	storageService := storage.NewStorageService(storagePath)
	api.InitStorage(storagePath)
	log.Printf("儲存服務已初始化,資料路徑: %s", storagePath)

	// 初始化 MQTT 客戶端 (可選)
	var mqttClient *mqtt.Client
	var err error

	if cfg.MQTTBroker != "" && cfg.MQTTBroker != "localhost" {
		mqttClient, err = mqtt.NewClient(cfg)
		if err != nil {
			log.Printf("警告: MQTT 初始化失敗: %v", err)
			log.Println("繼續運行但不使用 MQTT 功能")
		} else {
			defer mqttClient.Close()
			mqtt.SetStorageService(storageService)
			log.Println("MQTT 客戶端已成功初始化")
		}
	} else {
		log.Println("MQTT 未配置,僅啟用 API 功能")
	}

	// 設定路由
	router := api.SetupRouter()

	// 啟動服務器
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("TSPL Simulator 服務器啟動於 %s", serverAddr)
	log.Printf("API 端點: http://localhost%s/api", serverAddr)
	log.Printf("API 資料儲存: %s", filepath.Join(storagePath, "API_print"))
	log.Printf("MQTT 資料儲存: %s", filepath.Join(storagePath, "MQTT_print"))

	// 優雅關閉
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("服務器啟動失敗: %v", err)
		}
	}()

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在關閉服務器...")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
