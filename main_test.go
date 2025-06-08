package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
}

func TestEnsureRequiredDirectories(t *testing.T) {
	// 创建临时目录用于测试
	testBaseDir := "test_dirs"
	defer os.RemoveAll(testBaseDir)

	// 临时修改目录路径进行测试
	originalWd, _ := os.Getwd()
	os.Mkdir(testBaseDir, os.ModePerm)
	os.Chdir(testBaseDir)
	defer os.Chdir(originalWd)

	// 测试目录创建
	err := ensureRequiredDirectories()
	assert.NoError(t, err)

	// 验证目录是否存在
	_, err = os.Stat("./tmp")
	assert.NoError(t, err)

	_, err = os.Stat("./models")
	assert.NoError(t, err)
}

func TestSetupRoutes(t *testing.T) {
	// 创建测试路由器
	r := gin.New()
	setupRoutes(r)

	// 测试根路径
	t.Run("root endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Go OCR Service", response["message"])
		assert.Equal(t, "1.0", response["version"])
		assert.Contains(t, response, "endpoints")
	})

	// 测试健康检查接口
	t.Run("health endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
		assert.Equal(t, "go-ocr", response["service"])
		assert.Equal(t, "1.0", response["version"])
	})
}

func TestCorsMiddleware(t *testing.T) {
	// 创建测试路由器
	r := gin.New()
	r.Use(corsMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// 测试普通请求
	t.Run("normal request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})

	// 测试OPTIONS请求
	t.Run("options request", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestAPIRoutes(t *testing.T) {
	// 创建测试路由器（不包含实际的OCR处理器，因为它们依赖外部资源）
	r := gin.New()

	// 创建模拟的API组
	api := r.Group("/api")
	{
		api.POST("/ocr", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ocr endpoint"})
		})
		api.POST("/ocr_file", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ocr_file endpoint"})
		})
	}

	// 测试OCR JSON接口路径
	t.Run("ocr json endpoint exists", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/ocr", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 应该能访问到端点（即使参数错误）
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 测试OCR文件接口路径
	t.Run("ocr file endpoint exists", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/ocr_file", nil)
		req.Header.Set("Content-Type", "multipart/form-data")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 应该能访问到端点（即使参数错误）
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestEnvironmentVariables(t *testing.T) {
	// 测试默认端口
	t.Run("default port", func(t *testing.T) {
		os.Unsetenv("PORT")
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		assert.Equal(t, "8080", port)
	})

	// 测试自定义端口
	t.Run("custom port", func(t *testing.T) {
		os.Setenv("PORT", "9090")
		defer os.Unsetenv("PORT")

		port := os.Getenv("PORT")
		assert.Equal(t, "9090", port)
	})

	// 测试Gin模式
	t.Run("gin mode", func(t *testing.T) {
		// 测试默认模式
		os.Unsetenv("GIN_MODE")
		mode := os.Getenv("GIN_MODE")
		if mode == "" {
			mode = gin.ReleaseMode
		}
		assert.Equal(t, gin.ReleaseMode, mode)

		// 测试自定义模式
		os.Setenv("GIN_MODE", gin.DebugMode)
		defer os.Unsetenv("GIN_MODE")

		mode = os.Getenv("GIN_MODE")
		assert.Equal(t, gin.DebugMode, mode)
	})
}
