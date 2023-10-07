package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	Init()
	r := gin.Default()
	r.POST("/api/ocr", OcrJson)
	r.POST("/api/ocr_file", OcrFile)
	r.Run()
}
