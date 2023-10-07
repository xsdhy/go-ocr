/** @file ocr.h
  * @brief 
  * @author teng.qing
  * @date 8/13/21
  */
#include "ocr.h"
#include "OcrLite.h"
#include "omp.h"
#include "json.hpp"
#include <iostream>
#include <sys/stat.h>

using json = nlohmann::json;

static OcrLite *g_ocrLite = nullptr;

inline bool isFileExists(const char *name) {
    struct stat buffer{};
    return (stat(name, &buffer) == 0);
}

int ocr_init(int numThread, const char *dbNetPath, const char *anglePath, const char *crnnPath, const char *keyPath) {
    omp_set_num_threads(numThread); // 并行计算
    if (g_ocrLite == nullptr) {
        g_ocrLite = new OcrLite();
    }
    g_ocrLite->setNumThread(numThread);
    g_ocrLite->initLogger(
            true,//isOutputConsole
            false,//isOutputPartImg
            true);//isOutputResultImg
    g_ocrLite->Logger(
            "ocr_init numThread=%d, dbNetPath=%s,anglePath=%s,crnnPath=%s,keyPath=%s \n",
            numThread, dbNetPath, anglePath, crnnPath, keyPath);
    if (!isFileExists(dbNetPath) || !isFileExists(anglePath) || !isFileExists(crnnPath) || !isFileExists(keyPath)) {
        g_ocrLite->Logger("invalid file path.\n");
        return kOcrError;
    }
    g_ocrLite->initModels(dbNetPath, anglePath, crnnPath, keyPath);
    return kOcrSuccess;
}

void ocr_cleanup() {
    if (g_ocrLite != nullptr) {
        delete g_ocrLite;
        g_ocrLite = nullptr;
    }
}

int ocr_detect(const char *image_path, char *out_buffer, int *buffer_len, int padding, int maxSideLen,
               float boxScoreThresh, float boxThresh, float unClipRatio, bool doAngle, bool mostAngle) {
    if (g_ocrLite == nullptr) {
        return kOcrError;
    }
    if (!isFileExists(image_path)) {
        return kOcrError;
    }
    g_ocrLite->Logger(
            "padding(%d),maxSideLen(%d),boxScoreThresh(%f),boxThresh(%f),unClipRatio(%f),doAngle(%d),mostAngle(%d)\n",
            padding, maxSideLen, boxScoreThresh, boxThresh, unClipRatio, doAngle, mostAngle);
    OcrResult result = g_ocrLite->detect("", image_path, padding, maxSideLen,
                                         boxScoreThresh, boxThresh, unClipRatio, doAngle, mostAngle);
    json root;
    root["dbNetTime"] = result.dbNetTime;
    root["detectTime"] = result.detectTime;
    for (const auto &item : result.textBlocks) {
        json textBlock;
        for (const auto &boxPoint : item.boxPoint) {
            json point;
            point["x"] = boxPoint.x;
            point["y"] = boxPoint.y;
            textBlock["boxPoint"].push_back(point);
        }
        for (const auto &score : item.charScores) {
            textBlock["charScores"].push_back(score);
        }
        textBlock["text"] = item.text;
        textBlock["boxScore"] = item.boxScore;
        textBlock["angleIndex"] = item.angleIndex;
        textBlock["angleScore"] = item.angleScore;
        textBlock["angleTime"] = item.angleTime;
        textBlock["crnnTime"] = item.crnnTime;
        textBlock["blockTime"] = item.blockTime;
        root["textBlocks"].push_back(textBlock);
        root["texts"].push_back(item.text);
    }
    std::string tempJsonStr = root.dump();
    if (static_cast<int>(tempJsonStr.length()) > *buffer_len) {
        g_ocrLite->Logger("buff_len is too small \n");
        return kOcrError;
    }
    *buffer_len = static_cast<int>(tempJsonStr.length());
    ::memcpy(out_buffer, tempJsonStr.c_str(), tempJsonStr.length());
    return kOcrSuccess;
}

int ocr_detect2(const char *image_path, char *out_buffer, int *buffer_len) {
    return ocr_detect(image_path, out_buffer, buffer_len, kDefaultPadding, kDefaultMaxSideLen, kDefaultBoxScoreThresh,
                      kDefaultBoxThresh, kDefaultUnClipRatio, kDefaultDoAngle, kDefaultMostAngle);
}