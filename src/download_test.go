package src

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectImageType(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "PNG image",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00},
			expected: "png",
		},
		{
			name:     "JPEG image",
			data:     []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46},
			expected: "jpg",
		},
		{
			name:     "unknown format",
			data:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: "",
		},
		{
			name:     "too short data",
			data:     []byte{0x89, 0x50},
			expected: "",
		},
		{
			name:     "empty data",
			data:     []byte{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectImageType(tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateMD5Hash(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "empty data",
			data:     []byte{},
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "hello world",
			data:     []byte("hello world"),
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
		{
			name:     "binary data",
			data:     []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expected: "b59121341ab26766729b7f1d7f7e0c2f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMD5Hash(tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnsureTmpDir(t *testing.T) {
	// 测试创建临时目录
	err := ensureTmpDir()
	assert.NoError(t, err)

	// 验证目录是否存在
	_, err = os.Stat(tmpDir)
	assert.NoError(t, err)

	// 清理
	defer os.RemoveAll(tmpDir)
}

func TestGenerateUniqueFilename(t *testing.T) {
	// 测试生成唯一文件名
	filename1 := generateUniqueFilename(".jpg")
	filename2 := generateUniqueFilename(".png")

	// 应该生成不同的文件名
	assert.NotEqual(t, filename1, filename2)

	// 应该包含正确的扩展名
	assert.Contains(t, filename1, ".jpg")
	assert.Contains(t, filename2, ".png")

	// 应该包含tmp目录
	assert.Contains(t, filename1, tmpDir)
	assert.Contains(t, filename2, tmpDir)
}

func TestWriteFile(t *testing.T) {
	// 确保临时目录存在
	err := ensureTmpDir()
	assert.NoError(t, err)

	testFile := tmpDir + "/test_write.txt"
	testData := []byte("test content")

	// 测试写入文件
	err = writeFile(testFile, testData)
	assert.NoError(t, err)

	// 验证文件内容
	data, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, testData, data)

	// 清理
	defer func() {
		os.Remove(testFile)
		os.RemoveAll(tmpDir)
	}()
}

func TestSaveBase64Image(t *testing.T) {
	// 小的1x1像素PNG图片的base64编码
	validPNG := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="

	tests := []struct {
		name        string
		base64Data  string
		expectError bool
	}{
		{
			name:        "valid PNG base64",
			base64Data:  validPNG,
			expectError: false,
		},
		{
			name:        "empty base64",
			base64Data:  "",
			expectError: true,
		},
		{
			name:        "invalid base64",
			base64Data:  "invalid-base64-data",
			expectError: true,
		},
		{
			name:        "valid base64 but not image",
			base64Data:  base64.StdEncoding.EncodeToString([]byte("not an image")),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := saveBase64Image(tt.base64Data)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, filename)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, filename)

				// 验证文件是否存在
				_, err = os.Stat(filename)
				assert.NoError(t, err)

				// 清理
				defer os.Remove(filename)
			}
		})
	}

	// 清理临时目录
	defer os.RemoveAll(tmpDir)
}

func TestDownloadAndSaveImage_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		imageURL    string
		expectError bool
	}{
		{
			name:        "empty URL",
			imageURL:    "",
			expectError: true,
		},
		{
			name:        "invalid URL",
			imageURL:    "not-a-valid-url",
			expectError: true,
		},
		{
			name:        "non-existent URL",
			imageURL:    "http://non-existent-domain-123456.com/image.jpg",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := downloadAndSaveImage(tt.imageURL)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, filename)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, filename)

				// 清理
				defer os.Remove(filename)
			}
		})
	}

	// 清理临时目录
	defer os.RemoveAll(tmpDir)
}

func TestCalculateMD5_FileOperations(t *testing.T) {
	// 确保临时目录存在
	err := ensureTmpDir()
	assert.NoError(t, err)

	testFile := tmpDir + "/test_md5.txt"
	testData := []byte("test content for md5")

	// 创建测试文件
	err = os.WriteFile(testFile, testData, 0644)
	assert.NoError(t, err)

	// 测试计算MD5
	hash, err := calculateMD5(testFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 32) // MD5 hash is 32 characters long

	// 验证与内存计算的MD5是否一致
	expectedHash := calculateMD5Hash(testData)
	assert.Equal(t, expectedHash, hash)

	// 测试不存在的文件
	_, err = calculateMD5("non-existent-file.txt")
	assert.Error(t, err)

	// 清理
	defer func() {
		os.Remove(testFile)
		os.RemoveAll(tmpDir)
	}()
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "jpg extension",
			filename: "image.jpg",
			expected: "jpg",
		},
		{
			name:     "png extension",
			filename: "image.png",
			expected: "png",
		},
		{
			name:     "no extension",
			filename: "image",
			expected: "",
		},
		{
			name:     "multiple dots",
			filename: "image.backup.jpg",
			expected: "jpg",
		},
		{
			name:     "path with extension",
			filename: "/path/to/image.jpeg",
			expected: "jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFileExtension(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}
