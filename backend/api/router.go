package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 設定路由
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS 設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// API 路由
	api := router.Group("/api")
	{
		// 健康檢查
		api.GET("/health", HealthCheckHandler)

		// TSPL 渲染
		api.POST("/render", RenderHandler)

		// 範例管理
		api.GET("/examples", GetExamplesHandler)
		api.GET("/examples/:id", GetExampleDetailHandler)

		// MQTT 相關
		mqtt := api.Group("/mqtt")
		{
			mqtt.POST("/publish", MQTTPublishHandler)
		}
	}

	return router
}
