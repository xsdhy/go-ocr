//go:build test
// +build test

package src

import (
	"log"
	"runtime"
)

const kModelDbNet = "./models/dbnet.onnx"
const kModelAngle = "./models/angle_net.onnx"
const kModelCRNN = "./models/crnn_lite_lstm.onnx"
const kModelKeys = "./models/keys.txt"

const kDefaultBufferLen = 10 * 1024

type OCRBoxPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type OCRTextBlock struct {
	AngleIndex int           `json:"angle_index"`
	AngleScore float64       `json:"angle_score"`
	AngleTime  float64       `json:"angle_time"`
	BlockTime  float64       `json:"block_time"`
	BoxPoint   []OCRBoxPoint `json:"box_point"`
	BoxScore   float64       `json:"box_score"`
	CharScores []float64     `json:"char_scores"`
	CRNNTime   float64       `json:"crnn_time"`
	Text       string        `json:"text"`
}

type OCRResultData struct {
	DBNetTime  float64        `json:"db_net_time,omitempty"`
	DetectTime float64        `json:"detect_time,omitempty"`
	TextBlocks []OCRTextBlock `json:"text_blocks,omitempty"`
	Texts      []string       `json:"texts"`
	QRCode     bool           `json:"qr_code,omitempty"`
}

// 测试环境的存根实现
func Init() int {
	log.Println("OCR Init (stub implementation for testing)")
	threadNum := runtime.NumCPU()
	return threadNum
}

func Detect(imagePath string) (bool, *OCRResultData) {
	log.Printf("OCR Detect (stub implementation for testing): %s", imagePath)

	// 返回模拟的OCR结果
	result := &OCRResultData{
		DBNetTime:  0.1,
		DetectTime: 0.5,
		TextBlocks: []OCRTextBlock{
			{
				AngleIndex: 0,
				AngleScore: 0.99,
				AngleTime:  0.05,
				BlockTime:  0.2,
				BoxPoint: []OCRBoxPoint{
					{X: 10, Y: 10},
					{X: 100, Y: 10},
					{X: 100, Y: 30},
					{X: 10, Y: 30},
				},
				BoxScore:   0.92,
				CharScores: []float64{0.9, 0.8, 0.95, 0.85},
				CRNNTime:   0.1,
				Text:       "模拟识别结果",
			},
		},
		Texts:  []string{"模拟识别结果"},
		QRCode: false,
	}

	return true, result
}

func CleanUp() {
	log.Println("OCR CleanUp (stub implementation for testing)")
}
