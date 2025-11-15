package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tspl-simulator/config"
	"tspl-simulator/models"
	"tspl-simulator/parser"
	"tspl-simulator/storage"
	"tspl-simulator/validator"
)

type Client struct {
	client         mqtt.Client
	config         *config.Config
	storageService *storage.StorageService
}

var (
	mqttClient *Client
)

// SetStorageService 設定儲存服務
func SetStorageService(s *storage.StorageService) {
	if mqttClient != nil {
		mqttClient.storageService = s
	}
}

// NewClient 建立新的 MQTT 客戶端
func NewClient(cfg *config.Config) (*Client, error) {
	opts := mqtt.NewClientOptions()
	brokerURL := fmt.Sprintf("tcp://%s:%s", cfg.MQTTBroker, cfg.MQTTPort)
	opts.AddBroker(brokerURL)
	opts.SetClientID(cfg.MQTTClientID)

	if cfg.MQTTUsername != "" {
		opts.SetUsername(cfg.MQTTUsername)
		opts.SetPassword(cfg.MQTTPassword)
	}

	opts.SetDefaultPublishHandler(messageHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("MQTT 連接失敗: %v", token.Error())
	}

	mqttClient = &Client{
		client: client,
		config: cfg,
	}

	// 訂閱主題
	if token := client.Subscribe(cfg.MQTTTopic, 0, nil); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("MQTT 訂閱失敗: %v", token.Error())
	}

	log.Printf("MQTT 客戶端已連接到 %s 並訂閱主題 %s", brokerURL, cfg.MQTTTopic)

	return mqttClient, nil
}

// GetClient 取得 MQTT 客戶端實例
func GetClient() *Client {
	return mqttClient
}

// Publish 發布訊息到 MQTT
func (c *Client) Publish(topic string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化訊息失敗: %v", err)
	}

	token := c.client.Publish(topic, 0, false, payload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("發布訊息失敗: %v", token.Error())
	}

	log.Printf("已發布訊息到主題 %s", topic)
	return nil
}

// PublishRenderResult 發布渲染結果
func (c *Client) PublishRenderResult(renderData *models.RenderData) error {
	topic := c.config.MQTTTopic + "/result"

	message := models.MQTTMessage{
		Type:      "render_result",
		Timestamp: time.Now().Unix(),
	}

	return c.Publish(topic, map[string]interface{}{
		"type":      message.Type,
		"timestamp": message.Timestamp,
		"data":      renderData,
	})
}

// Close 關閉 MQTT 客戶端
func (c *Client) Close() {
	if c.client.IsConnected() {
		c.client.Disconnect(250)
		log.Println("MQTT 客戶端已斷開連接")
	}
}

// IsConnected 檢查 MQTT 是否已連接
func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}

// 訊息處理器
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("收到 MQTT 訊息 [%s]: %s", msg.Topic(), string(msg.Payload()))

	var message models.MQTTMessage
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Printf("解析 MQTT 訊息失敗: %v", err)
		return
	}

	// 處理不同類型的訊息
	switch message.Type {
	case "render_request":
		handleRenderRequest(message.TSPLCode)
	default:
		log.Printf("未知的訊息類型: %s", message.Type)
	}
}

// 處理渲染請求
func handleRenderRequest(tsplCode string) {
	log.Printf("處理 MQTT 渲染請求")

	// 驗證 TSPL 語法
	validationResult := validator.ValidateTSPL(tsplCode)
	if !validationResult.Valid {
		log.Printf("TSPL 語法驗證失敗:")
		for _, err := range validationResult.Errors {
			log.Printf("  行 %d [%s]: %s", err.Line, err.Command, err.Message)
		}
		return
	}

	// 儲存 MQTT 接收的資料
	if mqttClient != nil && mqttClient.storageService != nil {
		if filePath, err := mqttClient.storageService.SaveMQTTData(tsplCode); err != nil {
			log.Printf("儲存 MQTT 資料失敗: %v", err)
		} else {
			log.Printf("MQTT 資料已儲存至: %s", filePath)
		}
	}

	renderData, err := parser.ParseTSPL(tsplCode)
	if err != nil {
		log.Printf("解析 TSPL 失敗: %v", err)
		return
	}

	// 發布渲染結果
	if mqttClient != nil {
		if err := mqttClient.PublishRenderResult(renderData); err != nil {
			log.Printf("發布渲染結果失敗: %v", err)
		}
	}
}

// 連接處理器
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("MQTT 已連接")
}

// 斷線處理器
var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("MQTT 連接丟失: %v", err)
}
