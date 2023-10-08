package main

import (
	"github.com/gin-gonic/gin"
	"ocr/src"
)

func main() {
	src.Init()
	r := gin.Default()
	r.POST("/api/ocr", src.OcrJson)
	r.POST("/api/ocr_file", src.OcrFile)
	r.Run()
}
