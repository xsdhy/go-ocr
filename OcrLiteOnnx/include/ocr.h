/** @file ocr.h
  * @brief  封装给GO调用
  * @author teng.qing
  * @date 8/13/21
  */
#ifndef ONLINE_BASE_OCRLITEONNX_OCR_H
#define ONLINE_BASE_OCRLITEONNX_OCR_H

#ifdef __cplusplus
extern "C"
{
#else
    // c
    typedef enum{
        false, true
    }bool;
#endif

const int kOcrError = 0;
const int kOcrSuccess = 1;
const int kDefaultPadding = 50;
const int kDefaultMaxSideLen = 1024;
const float kDefaultBoxScoreThresh = 0.6f;
const float kDefaultBoxThresh = 0.3f;
const float kDefaultUnClipRatio = 2.0f;
const bool kDefaultDoAngle = true;
const bool kDefaultMostAngle = true;

/**@fn ocr_init
  *@brief 初始化OCR
  *@param numThread: 线程数量，不超过CPU数量
  *@param dbNetPath: dbnet模型路径
  *@param anglePath: 角度识别模型路径
  *@param crnnPath: crnn推理模型路径
  *@param keyPath: keys.txt样本路径
  *@return <0: error, >0: instance
  */
int ocr_init(int numThread, const char *dbNetPath, const char *anglePath, const char *crnnPath, const char *keyPath);

/**@fn ocr_cleanup
  *@brief 清理，退出程序前执行
  */
void ocr_cleanup();

/**@fn ocr_detect
  *@brief 识别图片
  *@param image_path: 图片完整路径，会在同路径下生成图片识别框选效果，便于调试
  *@param out_json_result: 识别结果输出，json格式。
  *@param buffer_len: 输出缓冲区大小
  *@param padding: 50
  *@param maxSideLen: 1024
  *@param boxScoreThresh: 0.6f
  *@param boxThresh: 0.3f
  *@param unClipRatio: 2.0f
  *@param doAngle: true
  *@param mostAngle: true
  *@return 成功与否
  */
int ocr_detect(const char *image_path, char *out_buffer, int *buffer_len, int padding, int maxSideLen,
                float boxScoreThresh, float boxThresh, float unClipRatio, bool doAngle, bool mostAngle);
                
/**@fn ocr_detect
  *@brief 使用默认参数，识别图片
  *@param image_path: 图片完整路径
  *@param out_buffer: 识别结果，json格式。
  *@param buffer_len: 输出缓冲区大小
  *@return 成功与否
  */
int ocr_detect2(const char *image_path, char *out_buffer, int *buffer_len);

#ifdef __cplusplus
}
#endif

#endif //ONLINE_BASE_OCRLITEONNX_OCR_H