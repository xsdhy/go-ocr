package src

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// QRCodeResult 二维码识别结果
type QRCodeResult struct {
	Found   bool   `json:"found"`
	Content string `json:"content,omitempty"`
}

// DetectQrCodeWithContent 检测二维码并返回内容
func DetectQrCodeWithContent(imagePath string) (bool, string) {
	result := detectQRCode(imagePath)
	return result.Found, result.Content
}

// DetectQRCode 检测二维码并返回完整结果
func DetectQRCode(imagePath string) *QRCodeResult {
	return detectQRCode(imagePath)
}

func detectQRCode(imagePath string) *QRCodeResult {
	if imagePath == "" {
		log.Println("二维码识别: 图片路径为空")
		return &QRCodeResult{Found: false}
	}

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("二维码识别: 图片文件不存在: %s", imagePath)
		return &QRCodeResult{Found: false}
	}

	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("二维码识别: 打开图片文件失败: %v", err)
		return &QRCodeResult{Found: false}
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("二维码识别: 关闭文件失败: %v", err)
		}
	}()

	// 解码图片
	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("二维码识别: 解码图片失败: %v", err)
		return &QRCodeResult{Found: false}
	}

	log.Printf("二维码识别: 成功解码图片格式: %s", format)

	// 创建bitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		log.Printf("二维码识别: 创建bitmap失败: %v", err)
		return &QRCodeResult{Found: false}
	}

	// 创建二维码读取器
	qrReader := qrcode.NewQRCodeReader()

	// 尝试读取二维码
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		log.Printf("二维码识别: 未检测到二维码: %v", err)
		return &QRCodeResult{Found: false}
	}

	content := result.GetText()
	log.Printf("二维码识别: 成功检测到二维码，内容长度: %d", len(content))

	return &QRCodeResult{
		Found:   true,
		Content: content,
	}
}

// ValidateImageForQRCode 验证图片是否适合进行二维码识别
func ValidateImageForQRCode(imagePath string) error {
	if imagePath == "" {
		return fmt.Errorf("图片路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("图片文件不存在: %s", imagePath)
	}

	// 尝试打开并解码图片
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("无法打开图片文件: %v", err)
	}
	defer file.Close()

	_, _, err = image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("无效的图片格式: %v", err)
	}

	return nil
}
