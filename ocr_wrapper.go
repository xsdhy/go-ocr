package main

// -I: 配置编译选项
// -L: 依赖库路径

/*
#cgo CFLAGS: -I ./cpp/install/include
#cgo LDFLAGS: -L ./cpp/install/lib -lOcrLiteOnnx -lstdc++

#include <stdlib.h>
#include <string.h>
#include "ocr.h"
*/
import "C"
import (
	"encoding/json"
	"runtime"
	"unsafe"
)

const kModelDbNet = "./models/dbnet.onnx"
const kModelAngle = "./models/angle_net.onnx"
const kModelCRNN = "./models/crnn_lite_lstm.onnx"
const kModelKeys = "./models/keys.txt"

const kDefaultBufferLen = 10 * 1024

var (
	buffer [kDefaultBufferLen]byte
)

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
}

func Init() int {
	// dbNet, angle, crnn, keys string
	threadNum := runtime.NumCPU()
	cDbNet := C.CString(kModelDbNet) // to c char*
	cAngle := C.CString(kModelAngle) // to c char*
	cCRNN := C.CString(kModelCRNN)   // to c char*
	cKeys := C.CString(kModelKeys)   // to c char*

	ret := C.ocr_init(C.int(threadNum), cDbNet, cAngle, cCRNN, cKeys)

	C.free(unsafe.Pointer(cDbNet))
	C.free(unsafe.Pointer(cAngle))
	C.free(unsafe.Pointer(cCRNN))
	C.free(unsafe.Pointer(cKeys))
	return int(ret)
}

func Detect(imagePath string) (bool, *OCRResultData) {
	resultLen := C.int(kDefaultBufferLen)

	// 构造C的缓冲区
	cTempBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	cImagePath := C.CString(imagePath)
	defer C.free(unsafe.Pointer(cImagePath))

	isSuccess := C.ocr_detect2(cImagePath, cTempBuffer, &resultLen)
	if int(isSuccess) != 1 {
		return false, nil
	}
	result := C.GoStringN(cTempBuffer, resultLen)
	var vo OCRResultData

	err := json.Unmarshal([]byte(result), &vo)
	if err != nil {
		return false, nil
	}

	return true, &vo
}

func CleanUp() {
	C.ocr_cleanup()
}
