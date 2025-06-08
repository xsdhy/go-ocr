package src

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type OcrDTO struct {
	ImageUrl    string `json:"image_url" binding:"omitempty,url"`
	ImageBase64 string `json:"image_base_64"`
	NeedBlock   bool   `json:"need_block"`
	QrCode      bool   `json:"qr_code"` // 是否识别二维码
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// validateOcrDTO 验证输入参数
func validateOcrDTO(input *OcrDTO) error {
	if input.ImageUrl == "" && input.ImageBase64 == "" {
		return fmt.Errorf("image_url和image_base64至少需要提供一个")
	}
	if input.ImageUrl != "" && input.ImageBase64 != "" {
		return fmt.Errorf("image_url和image_base64只能提供一个")
	}
	if input.ImageBase64 != "" && !isValidBase64(input.ImageBase64) {
		return fmt.Errorf("无效的base64图片数据")
	}
	return nil
}

// isValidBase64 简单验证base64格式
func isValidBase64(s string) bool {
	return len(s) > 0 && len(s)%4 == 0 && !strings.ContainsAny(s, " \t\n\r")
}

func OcrJson(c *gin.Context) {
	var input OcrDTO
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		log.Printf("参数绑定失败: %v", err)
		SendError(c, "参数格式错误: "+err.Error())
		return
	}

	// 验证输入参数
	if err := validateOcrDTO(&input); err != nil {
		log.Printf("参数验证失败: %v", err)
		SendError(c, err.Error())
		return
	}

	var imagePath string
	var err error

	if input.ImageBase64 != "" {
		log.Println("处理base64图片")
		imagePath, err = saveBase64Image(input.ImageBase64)
	} else if input.ImageUrl != "" {
		log.Printf("下载图片: %s", input.ImageUrl)
		imagePath, err = downloadAndSaveImage(input.ImageUrl)
	}

	if err != nil {
		log.Printf("图片处理失败: %v", err)
		SendError(c, "图片处理失败: "+err.Error())
		return
	}

	// 确保文件存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("图片文件不存在: %s", imagePath)
		SendError(c, "图片文件处理失败")
		return
	}

	// 识别
	response, err := performOCR(imagePath, input)
	if err != nil {
		log.Printf("OCR识别失败: %v", err)
		SendError(c, "OCR识别失败: "+err.Error())
	} else {
		log.Printf("OCR识别成功: %s", imagePath)
		c.JSON(http.StatusOK, response)
	}
}

func OcrFile(c *gin.Context) {
	var input OcrDTO

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("获取上传文件失败: %v", err)
		SendError(c, "获取上传文件失败: "+err.Error())
		return
	}

	// 验证文件类型
	if !isValidImageFile(file.Filename) {
		log.Printf("不支持的文件类型: %s", file.Filename)
		SendError(c, "不支持的文件类型，请上传jpg、jpeg、png格式的图片")
		return
	}

	// 获取表单参数
	if c.DefaultPostForm("need_block", "") == "true" {
		input.NeedBlock = true
	}
	if c.DefaultPostForm("qr_code", "") == "true" {
		input.QrCode = true
	}

	log.Printf("处理上传文件: %s", file.Filename)

	// 确保tmp目录存在
	if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
		log.Printf("创建临时目录失败: %v", err)
		SendError(c, "服务器内部错误")
		return
	}

	imagePath := "./tmp/" + file.Filename

	// 保存上传的文件
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		log.Printf("保存上传文件失败: %v", err)
		SendError(c, "保存文件失败: "+err.Error())
		return
	}

	// 识别
	response, err := performOCR(imagePath, input)
	if err != nil {
		log.Printf("OCR识别失败: %v", err)
		SendError(c, "OCR识别失败: "+err.Error())
	} else {
		log.Printf("OCR识别成功: %s", imagePath)
		c.JSON(http.StatusOK, response)
	}
}

// performOCR 执行OCR识别的核心逻辑
func performOCR(imagePath string, input OcrDTO) (*Response, error) {
	// 确保在函数结束时清理临时文件
	defer func() {
		cleanupFiles(imagePath)
	}()

	detect, ocrResult := Detect(imagePath)
	if !detect {
		return nil, fmt.Errorf("OCR识别失败")
	}

	// 根据需要返回文本块信息
	if !input.NeedBlock {
		ocrResult.TextBlocks = nil
	}

	// 如果需要识别二维码
	if input.QrCode {
		hasQR, qrContent := DetectQrCodeWithContent(imagePath)
		ocrResult.QRCode = hasQR
		if hasQR {
			// 可以在这里添加二维码内容到结果中
			log.Printf("检测到二维码: %s", qrContent)
		}
	}

	return &Response{Code: 200, Msg: "ok", Data: ocrResult}, nil
}

// cleanupFiles 清理临时文件
func cleanupFiles(imagePath string) {
	files := []string{imagePath, imagePath + "-result.jpg"}
	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			log.Printf("清理文件失败 %s: %v", file, err)
		}
	}
}

// isValidImageFile 验证文件类型
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filename)
	return strings.HasSuffix(ext, ".jpg") ||
		strings.HasSuffix(ext, ".jpeg") ||
		strings.HasSuffix(ext, ".png")
}

func SendError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{Code: 500, Msg: message, Data: nil})
}
