FROM xsdhy/chineseocr_lite_cpp as builder

WORKDIR /app
ADD . /app

RUN cd /app/cpp/onnxruntime-static && \
    wget https://github.com/RapidAI/OnnxruntimeBuilder/releases/download/1.9.0/onnxruntime-1.9.0-ubuntu1804-static.7z && \
    7z X onnxruntime-1.9.0-ubuntu1804-static.7z && \
    cd /app/cpp/opencv-static && \
    wget https://github.com/RapidAI/OpenCVBuilder/releases/download/4.5.4/opencv-4.5.4-ubuntu1804.7z && \
    7z X opencv-4.5.4-ubuntu1804.7z

RUN cd /app/cpp && \
    cmake -DCMAKE_INSTALL_PREFIX=install -DCMAKE_BUILD_TYPE=Release -DOCR_OUTPUT=CLIB && \
    cmake --build . --config Release  && \
    cmake --build . --config Release --target install

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y build-essential wget
RUN wget -c https://dl.google.com/go/go1.20.9.linux-amd64.tar.gz -O - | tar -xz -C /usr/local
ENV PATH=$PATH:/usr/local/go/bin

COPY --from=builder /app /app

WORKDIR /app

RUN go mod tidy && GOOS=linux CGO_ENABLED=1 GOARCH=amd64 CGO_LDFLAGS="-g -O2 -Wl,--no-as-needed -ldl" go build -ldflags="-s -w" -installsuffix cgo -o ocr .

ENV LD_LIBRARY_PATH=/app/cpp/install/lib

CMD ["./ocr"]
EXPOSE 8080