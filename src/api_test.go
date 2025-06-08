package src

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置测试环境为test模式
	gin.SetMode(gin.TestMode)
}

func TestValidateOcrDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   OcrDTO
		wantErr bool
	}{
		{
			name: "valid image url",
			input: OcrDTO{
				ImageUrl: "http://example.com/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "valid base64",
			input: OcrDTO{
				ImageBase64: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   OcrDTO{},
			wantErr: true,
		},
		{
			name: "both url and base64",
			input: OcrDTO{
				ImageUrl:    "http://example.com/image.jpg",
				ImageBase64: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
			},
			wantErr: true,
		},
		{
			name: "invalid base64",
			input: OcrDTO{
				ImageBase64: "invalid base64 data!",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOcrDTO(&tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid base64",
			input:    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid length",
			input:    "abc",
			expected: false,
		},
		{
			name:     "contains whitespace",
			input:    "abcd efgh",
			expected: false,
		},
		{
			name:     "contains newline",
			input:    "abcd\nefgh",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidBase64(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidImageFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"valid jpg", "test.jpg", true},
		{"valid jpeg", "test.jpeg", true},
		{"valid png", "test.png", true},
		{"valid JPG uppercase", "test.JPG", true},
		{"invalid txt", "test.txt", false},
		{"invalid no extension", "test", false},
		{"invalid pdf", "test.pdf", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidImageFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOcrJsonAPI(t *testing.T) {
	// 创建测试路由
	router := gin.New()
	router.POST("/api/ocr", OcrJson)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "empty payload",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
			expectedCode:   500,
		},
		{
			name: "invalid json",
			payload: map[string]interface{}{
				"image_url": "invalid-url",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   500,
		},
		{
			name: "both url and base64",
			payload: map[string]interface{}{
				"image_url":     "http://example.com/image.jpg",
				"image_base_64": "somebase64",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/ocr", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

func TestOcrFileAPI(t *testing.T) {
	// 创建测试路由
	router := gin.New()
	router.POST("/api/ocr_file", OcrFile)

	t.Run("no file uploaded", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/ocr_file", nil)
		req.Header.Set("Content-Type", "multipart/form-data")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 500, response.Code)
	})

	t.Run("invalid file type", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 创建一个文本文件字段
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)

		_, err = io.WriteString(part, "some text content")
		assert.NoError(t, err)

		writer.Close()

		req, _ := http.NewRequest("POST", "/api/ocr_file", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Msg, "不支持的文件类型")
	})
}

func TestCleanupFiles(t *testing.T) {
	// 创建临时测试文件
	testFile := "test_cleanup.txt"
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	assert.NoError(t, err)

	// 确保文件存在
	_, err = os.Stat(testFile)
	assert.NoError(t, err)

	// 调用清理函数
	cleanupFiles(testFile)

	// 验证文件已被删除
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestSendError(t *testing.T) {
	// 创建测试路由
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		SendError(c, "test error message")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Equal(t, "test error message", response.Msg)
	assert.Nil(t, response.Data)
}
