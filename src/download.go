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

func saveBase64Image(base64String string) (string, error) {
	// 创建一个临时目录
	tmpDir := "tmp"
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", err
	}

	// 创建一个随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 解码base64字符串
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}

	// 生成一个唯一的文件名，包括时间戳和4位随机数
	timestamp := time.Now().UnixNano()
	randomNum := rand.Intn(10000)
	filename := fmt.Sprintf("tmp/%d_%04d.png", timestamp, randomNum)

	// 创建文件并写入解码后的图片数据
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}

	return filename, nil
}

func downloadAndSaveImage(imageURL string) (string, error) {
	// 创建一个临时目录
	tmpDir := "tmp"
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", err
	}

	// 从URL中获取文件名
	urlParts := strings.Split(imageURL, "/")
	fileName := urlParts[len(urlParts)-1]

	// 创建一个文件来保存图片
	filePath := filepath.Join(tmpDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 发起HTTP GET请求并下载图片
	response, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// 将图片数据复制到文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	// 计算图片的哈希值（这里使用MD5作为示例）
	hashValue, err := calculateMD5(filePath)
	if err != nil {
		return "", err
	}

	// 使用哈希值来重命名图片文件
	newFilePath := filepath.Join(tmpDir, hashValue+"."+getFileExtension(fileName))
	err = os.Rename(filePath, newFilePath)
	if err != nil {
		return "", err
	}

	return newFilePath, nil
}

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 将哈希值转换为16进制字符串
	hashValue := fmt.Sprintf("%x", hash.Sum(nil))
	return hashValue, nil
}

func getFileExtension(fileName string) string {
	extension := filepath.Ext(fileName)
	return strings.TrimPrefix(extension, ".")
}
