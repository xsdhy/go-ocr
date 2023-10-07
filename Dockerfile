FROM ubuntu:18.04
RUN apt-get update && apt-get install -y build-essential wget libssl-dev && apt-get remove cmake

RUN cd /usr/local/src && \
    wget wget https://github.com/Kitware/CMake/releases/download/v3.27.7/cmake-3.27.7.tar.gz && \
    tar xzvf cmake-3.27.7.tar.gz && \
    cd cmake-3.27.7 && \
    ./bootstrap && \
    make && \
    make install

WORKDIR /app
ADD . /app

RUN cd /app/cpp/onnxruntime-static && \
    wget https://github.com/RapidAI/OnnxruntimeBuilder/releases/download/1.9.0/onnxruntime-1.9.0-ubuntu1804-static.7z && \



FROM ubuntu:18.04

RUN apt-get update && apt-get install -y build-essential wget
RUN wget -c https://dl.google.com/go/go1.20.9.linux-amd64.tar.gz -O - | tar -xz -C /usr/local
ENV PATH=$PATH:/usr/local/go/bin
WORKDIR /app
ADD . /app

RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 CGO_LDFLAGS="-g -O2 -Wl,--no-as-needed -ldl" go build -ldflags="-s -w" -installsuffix cgo -o ocr .

ENV LD_LIBRARY_PATH=/app/lib/lib

CMD ["./ocr"]
EXPOSE 8080