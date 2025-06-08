package src

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectQrCodeWithContent(t *testing.T) {
	// 测试不存在的文件
	t.Run("non-existent file", func(t *testing.T) {
		found, content := DetectQrCodeWithContent("non-existent.png")
		assert.False(t, found)
		assert.Empty(t, content)
	})

	// 测试空路径
	t.Run("empty path", func(t *testing.T) {
		found, content := DetectQrCodeWithContent("")
		assert.False(t, found)
		assert.Empty(t, content)
	})

	// 创建一个非图片文件进行测试
	t.Run("invalid image file", func(t *testing.T) {
		// 创建临时文本文件
		testFile := "test_invalid.txt"
		err := os.WriteFile(testFile, []byte("not an image"), 0644)
		assert.NoError(t, err)
		defer os.Remove(testFile)

		found, content := DetectQrCodeWithContent(testFile)
		assert.False(t, found)
		assert.Empty(t, content)
	})
}

func TestDetectQRCode(t *testing.T) {
	// 测试不存在的文件
	t.Run("non-existent file", func(t *testing.T) {
		result := DetectQRCode("non-existent.png")
		assert.NotNil(t, result)
		assert.False(t, result.Found)
		assert.Empty(t, result.Content)
	})

	// 测试空路径
	t.Run("empty path", func(t *testing.T) {
		result := DetectQRCode("")
		assert.NotNil(t, result)
		assert.False(t, result.Found)
		assert.Empty(t, result.Content)
	})
}

func TestValidateImageForQRCode(t *testing.T) {
	// 测试空路径
	t.Run("empty path", func(t *testing.T) {
		err := ValidateImageForQRCode("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "图片路径不能为空")
	})

	// 测试不存在的文件
	t.Run("non-existent file", func(t *testing.T) {
		err := ValidateImageForQRCode("non-existent.png")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "图片文件不存在")
	})

	// 测试无效的图片文件
	t.Run("invalid image file", func(t *testing.T) {
		// 创建临时文本文件
		testFile := "test_invalid_image.txt"
		err := os.WriteFile(testFile, []byte("not an image"), 0644)
		assert.NoError(t, err)
		defer os.Remove(testFile)

		err = ValidateImageForQRCode(testFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的图片格式")
	})

	// 如果有有效的图片文件，可以测试成功的情况
	// 这里由于没有真实的图片文件，我们创建一个简单的PNG文件进行测试
	t.Run("valid PNG file", func(t *testing.T) {
		// 创建一个最小的PNG文件
		// 1x1像素的PNG图片数据
		pngData := []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
			0x00, 0x00, 0x00, 0x0D, // IHDR chunk size
			0x49, 0x48, 0x44, 0x52, // IHDR
			0x00, 0x00, 0x00, 0x01, // width: 1
			0x00, 0x00, 0x00, 0x01, // height: 1
			0x08, 0x06, 0x00, 0x00, 0x00, // bit depth, color type, compression, filter, interlace
			0x1F, 0x15, 0xC4, 0x89, // CRC
			0x00, 0x00, 0x00, 0x0A, // IDAT chunk size
			0x49, 0x44, 0x41, 0x54, // IDAT
			0x78, 0x9C, 0x62, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, // compressed data
			0xE2, 0x21, 0xBC, 0x33, // CRC
			0x00, 0x00, 0x00, 0x00, // IEND chunk size
			0x49, 0x45, 0x4E, 0x44, // IEND
			0xAE, 0x42, 0x60, 0x82, // CRC
		}

		testFile := "test_valid.png"
		err := os.WriteFile(testFile, pngData, 0644)
		assert.NoError(t, err)
		defer os.Remove(testFile)

		err = ValidateImageForQRCode(testFile)
		assert.NoError(t, err)
	})
}

func TestQRCodeResult(t *testing.T) {
	// 测试QRCodeResult结构体
	result := &QRCodeResult{
		Found:   true,
		Content: "test content",
	}

	assert.True(t, result.Found)
	assert.Equal(t, "test content", result.Content)

	// 测试空结果
	emptyResult := &QRCodeResult{
		Found: false,
	}

	assert.False(t, emptyResult.Found)
	assert.Empty(t, emptyResult.Content)
}
