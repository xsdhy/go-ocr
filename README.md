# Go Ocr

基于[chineseocr_lite](https://github.com/DayBreak-u/chineseocr_lite)的golang实现的ocr。

具体来说使用chineseocr_lite的[OcrLiteOnnx](https://github.com/DayBreak-u/chineseocr_lite/blob/onnx/cpp_projects/OcrLiteOnnx/README.md)
项目，参考[xmcy001122](https://blog.csdn.net/xmcy001122/article/details/119795546)
增加ocr.h和ocr.cpp，导出c风格函数,提供给ocr_wrapper.go使用cgo进行调用

## 使用

```bash
docker run --name ocr -itd --rm -p 8080:8080 xsdhy/go-ocr:1.0
```

## 接口文档

### 识别接口

#### 请求
| 参数            | 类型     | 是否必填           | 说明 |
|---------------|--------|----------------|----|
| image_url     | string | 图片地址和base64二选一 |    |
| image_base_64 | string | 图片地址和base64二选一 |    |
| need_block    | bool   | 否，默认为false     |    |

```bash
curl --location 'http://127.0.0.1:8080/api/ocr' \
--header 'Content-Type: application/json' \
--data '{
    "image_url":"图片地址"
}'
```

#### 响应

```bash
{
    "code": 200,
    "msg": "ok",
    "data": {
        "texts": [
            "第一行识别结果",
            "第二行识别结果",
            "第三行识别结果",
        ]
    }
}
```

### 识别接口(表单)
支持图片上传文件，接口地址为/api/ocr_file，文件key为file，其余的和识别接口相同