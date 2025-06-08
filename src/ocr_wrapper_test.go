package src

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOCRBoxPoint(t *testing.T) {
	// 测试OCRBoxPoint结构体
	point := OCRBoxPoint{
		X: 100,
		Y: 200,
	}

	assert.Equal(t, 100, point.X)
	assert.Equal(t, 200, point.Y)

	// 测试JSON序列化
	jsonData, err := json.Marshal(point)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"x":100`)
	assert.Contains(t, string(jsonData), `"y":200`)
}

func TestOCRTextBlock(t *testing.T) {
	// 创建测试用的OCRTextBlock
	textBlock := OCRTextBlock{
		AngleIndex: 0,
		AngleScore: 0.99,
		AngleTime:  0.1,
		BlockTime:  0.5,
		BoxPoint: []OCRBoxPoint{
			{X: 10, Y: 20},
			{X: 100, Y: 20},
			{X: 100, Y: 50},
			{X: 10, Y: 50},
		},
		BoxScore:   0.95,
		CharScores: []float64{0.9, 0.8, 0.9, 0.85},
		CRNNTime:   0.3,
		Text:       "测试文本",
	}

	// 验证字段
	assert.Equal(t, 0, textBlock.AngleIndex)
	assert.Equal(t, 0.99, textBlock.AngleScore)
	assert.Equal(t, "测试文本", textBlock.Text)
	assert.Len(t, textBlock.BoxPoint, 4)
	assert.Len(t, textBlock.CharScores, 4)

	// 测试JSON序列化
	jsonData, err := json.Marshal(textBlock)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"text":"测试文本"`)
	assert.Contains(t, string(jsonData), `"box_score":0.95`)

	// 测试JSON反序列化
	var deserializedBlock OCRTextBlock
	err = json.Unmarshal(jsonData, &deserializedBlock)
	assert.NoError(t, err)
	assert.Equal(t, textBlock.Text, deserializedBlock.Text)
	assert.Equal(t, textBlock.BoxScore, deserializedBlock.BoxScore)
}

func TestOCRResultData(t *testing.T) {
	// 创建测试用的OCRResultData
	resultData := OCRResultData{
		DBNetTime:  0.1,
		DetectTime: 0.5,
		TextBlocks: []OCRTextBlock{
			{
				Text:     "第一行文本",
				BoxScore: 0.95,
			},
			{
				Text:     "第二行文本",
				BoxScore: 0.88,
			},
		},
		Texts:  []string{"第一行文本", "第二行文本"},
		QRCode: true,
	}

	// 验证字段
	assert.Equal(t, 0.1, resultData.DBNetTime)
	assert.Equal(t, 0.5, resultData.DetectTime)
	assert.Len(t, resultData.TextBlocks, 2)
	assert.Len(t, resultData.Texts, 2)
	assert.True(t, resultData.QRCode)

	// 验证文本内容
	assert.Contains(t, resultData.Texts, "第一行文本")
	assert.Contains(t, resultData.Texts, "第二行文本")

	// 测试JSON序列化
	jsonData, err := json.Marshal(resultData)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"texts":["第一行文本","第二行文本"]`)
	assert.Contains(t, string(jsonData), `"qr_code":true`)

	// 测试JSON反序列化
	var deserializedData OCRResultData
	err = json.Unmarshal(jsonData, &deserializedData)
	assert.NoError(t, err)
	assert.Equal(t, resultData.QRCode, deserializedData.QRCode)
	assert.Equal(t, len(resultData.Texts), len(deserializedData.Texts))
}

func TestOCRConstants(t *testing.T) {
	// 验证OCR模型路径常量
	assert.Equal(t, "./models/dbnet.onnx", kModelDbNet)
	assert.Equal(t, "./models/angle_net.onnx", kModelAngle)
	assert.Equal(t, "./models/crnn_lite_lstm.onnx", kModelCRNN)
	assert.Equal(t, "./models/keys.txt", kModelKeys)
	assert.Equal(t, 10*1024, kDefaultBufferLen)
}

func TestJSONMarshalUnmarshal(t *testing.T) {
	// 测试复杂的JSON序列化和反序列化
	originalData := OCRResultData{
		DBNetTime:  0.15,
		DetectTime: 0.8,
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
				CharScores: []float64{0.9, 0.8, 0.95, 0.85, 0.9},
				CRNNTime:   0.1,
				Text:       "Hello World",
			},
		},
		Texts:  []string{"Hello World"},
		QRCode: false,
	}

	// 序列化
	jsonBytes, err := json.Marshal(originalData)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// 反序列化
	var parsedData OCRResultData
	err = json.Unmarshal(jsonBytes, &parsedData)
	assert.NoError(t, err)

	// 验证数据完整性
	assert.Equal(t, originalData.DBNetTime, parsedData.DBNetTime)
	assert.Equal(t, originalData.DetectTime, parsedData.DetectTime)
	assert.Equal(t, originalData.QRCode, parsedData.QRCode)
	assert.Len(t, parsedData.TextBlocks, 1)
	assert.Len(t, parsedData.Texts, 1)

	// 验证TextBlock数据
	originalBlock := originalData.TextBlocks[0]
	parsedBlock := parsedData.TextBlocks[0]
	assert.Equal(t, originalBlock.Text, parsedBlock.Text)
	assert.Equal(t, originalBlock.BoxScore, parsedBlock.BoxScore)
	assert.Equal(t, originalBlock.AngleScore, parsedBlock.AngleScore)
	assert.Len(t, parsedBlock.BoxPoint, len(originalBlock.BoxPoint))
	assert.Len(t, parsedBlock.CharScores, len(originalBlock.CharScores))
}

func TestEmptyOCRResultData(t *testing.T) {
	// 测试空的OCRResultData
	emptyData := OCRResultData{}

	// 验证默认值
	assert.Equal(t, float64(0), emptyData.DBNetTime)
	assert.Equal(t, float64(0), emptyData.DetectTime)
	assert.Nil(t, emptyData.TextBlocks)
	assert.Nil(t, emptyData.Texts)
	assert.False(t, emptyData.QRCode)

	// 测试JSON序列化
	jsonData, err := json.Marshal(emptyData)
	assert.NoError(t, err)
	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"texts":null`)
	// text_blocks字段在omitempty的情况下可能不会出现在JSON中
	assert.True(t, strings.Contains(jsonStr, `"text_blocks":null`) || !strings.Contains(jsonStr, `"text_blocks"`))
}

// 注意：这里不测试实际的OCR功能（Init, Detect, CleanUp），
// 因为这些函数依赖于C库和模型文件，在单元测试环境中可能不可用
