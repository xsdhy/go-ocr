package src

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	tmpDir      = "tmp"
)

// init 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// ensureTmpDir 确保临时目录存在
func ensureTmpDir() error {
	return os.MkdirAll(tmpDir, os.ModePerm)
}

// generateUniqueFilename 生成唯一的文件名
func generateUniqueFilename(ext string) string {
	timestamp := time.Now().UnixNano()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("%s/%d_%04d%s", tmpDir, timestamp, randomNum, ext)
}

func saveBase64Image(base64String string) (string, error) {
	if base64String == "" {
		return "", fmt.Errorf("base64字符串不能为空")
	}

	// 确保临时目录存在
	if err := ensureTmpDir(); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 解码base64字符串
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %v", err)
	}

	// 检查文件大小
	if len(decoded) > maxFileSize {
		return "", fmt.Errorf("图片文件过大，最大支持%dMB", maxFileSize/(1024*1024))
	}

	// 检查是否为有效的图片格式
	imageType := detectImageType(decoded)
	if imageType == "" {
		return "", fmt.Errorf("不支持的图片格式")
	}

	// 生成唯一的文件名
	filename := generateUniqueFilename("." + imageType)

	// 创建文件并写入解码后的图片数据
	if err := writeFile(filename, decoded); err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	return filename, nil
}

func downloadAndSaveImage(imageURL string) (string, error) {
	if imageURL == "" {
		return "", fmt.Errorf("图片URL不能为空")
	}

	// 确保临时目录存在
	if err := ensureTmpDir(); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发起HTTP GET请求
	response, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP请求失败，状态码: %d", response.StatusCode)
	}

	// 检查Content-Length头部，防止下载过大的文件
	if response.ContentLength > maxFileSize {
		return "", fmt.Errorf("图片文件过大，最大支持%dMB", maxFileSize/(1024*1024))
	}

	// 读取响应体，限制最大读取大小
	data, err := io.ReadAll(io.LimitReader(response.Body, maxFileSize))
	if err != nil {
		return "", fmt.Errorf("读取图片数据失败: %v", err)
	}

	// 检查是否为有效的图片格式
	imageType := detectImageType(data)
	if imageType == "" {
		return "", fmt.Errorf("不支持的图片格式")
	}

	// 计算图片的哈希值用于去重
	hashValue := calculateMD5Hash(data)
	filename := filepath.Join(tmpDir, hashValue+"."+imageType)

	// 检查文件是否已存在，如果存在则直接返回
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}

	// 保存图片文件
	if err := writeFile(filename, data); err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	return filename, nil
}

// writeFile 安全地写入文件
func writeFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		// 如果写入失败，删除已创建的文件
		os.Remove(filename)
		return err
	}

	return nil
}

// detectImageType 检测图片类型
func detectImageType(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	// 检查PNG格式
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "png"
	}

	// 检查JPEG格式
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "jpg"
	}

	return ""
}

// calculateMD5Hash 计算数据的MD5哈希值
func calculateMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

// calculateMD5 计算文件的MD5哈希值 (保留原有函数用于兼容)
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("计算哈希值失败: %v", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// getFileExtension 获取文件扩展名 (保留原有函数用于兼容)
func getFileExtension(fileName string) string {
	extension := filepath.Ext(fileName)
	return strings.TrimPrefix(extension, ".")
}
