package main

import (
	"log"
	"ocr/src"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 检查必要的目录
	if err := ensureRequiredDirectories(); err != nil {
		log.Fatalf("初始化目录失败: %v", err)
	}

	// 初始化OCR
	log.Println("正在初始化OCR引擎...")
	if result := src.Init(); result != 0 {
		log.Printf("OCR初始化完成，返回值: %d", result)
	} else {
		log.Println("OCR初始化可能失败，但服务将继续运行")
	}

	// 确保在程序退出时清理资源
	defer func() {
		log.Println("正在清理OCR资源...")
		src.CleanUp()
		log.Println("服务已停止")
	}()

	// 设置Gin模式
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	// 创建Gin路由器
	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 添加CORS中间件
	r.Use(corsMiddleware())

	// 注册路由
	setupRoutes(r)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务正在启动，监听端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// ensureRequiredDirectories 确保必要的目录存在
func ensureRequiredDirectories() error {
	directories := []string{"./tmp", "./models"}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
		log.Printf("目录检查完成: %s", dir)
	}

	return nil
}

// setupRoutes 设置API路由
func setupRoutes(r *gin.Engine) {
	// API组
	api := r.Group("/api")
	{
		api.POST("/ocr", src.OcrJson)
		api.POST("/ocr_file", src.OcrFile)
	}

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "go-ocr",
			"version": "1.0",
		})
	})

}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
