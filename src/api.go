package src

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type OcrDTO struct {
	ImageUrl    string `json:"image_url"`
	ImageBase64 string `json:"image_base_64"`
	NeedBlock   bool   `json:"need_block"`
}
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OcrJson(c *gin.Context) {
	var input OcrDTO
	err := c.ShouldBindWith(&input, binding.JSON)
	if nil != err {
		SendError(c, err.Error())
		return
	}

	var imagePath string
	if input.ImageBase64 != "" {
		imagePath, err = saveBase64Image(input.ImageBase64)
	} else if input.ImageUrl != "" {
		imagePath, err = downloadAndSaveImage(input.ImageUrl)
	}
	if err != nil {
		SendError(c, err.Error())
		return
	}
	if imagePath == "" {
		SendError(c, "请输入图片信息")
		return
	}

	//识别
	err, response := ocr(imagePath, input)
	if err != nil {
		SendError(c, err.Error())
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func OcrFile(c *gin.Context) {
	var input OcrDTO

	// 单文件
	file, _ := c.FormFile("file")
	if c.DefaultPostForm("need_block", "") == "true" {
		input.NeedBlock = true
	}

	log.Println(file.Filename)
	imagePath := "./tmp/" + file.Filename
	// 上传文件至指定的完整文件路径
	err := c.SaveUploadedFile(file, imagePath)
	if err != nil {
		SendError(c, err.Error())
		return
	}
	//识别
	err, response := ocr(imagePath, input)
	if err != nil {
		SendError(c, err.Error())
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func ocr(imagePath string, input OcrDTO) (error, *Response) {
	detect, ocrResult := Detect(imagePath)
	_ = os.Remove(imagePath)
	_ = os.Remove(imagePath + "-result.jpg")

	if !detect {
		return errors.New("识别失败"), nil
	}
	if !input.NeedBlock {
		ocrResult.TextBlocks = nil
	}
	return nil, &Response{Code: 200, Msg: "ok", Data: ocrResult}
}

func SendError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{Code: 500, Msg: message, Data: nil})
}
