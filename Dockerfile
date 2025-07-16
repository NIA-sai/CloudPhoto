# Dockerfile
# docker run -d --name cloudphoto -p 23322:23322   -v ./config.yaml:/root/config.yaml  niasai/cloudphoto
# 1. 使用官方 Go 镜像作为构建环境

FROM golang:latest AS builder
ENV GOPROXY=https://goproxy.cn,direct

# 2. 设置工作目录
WORKDIR /scr/go/cp

# 3. 复制代码并下载依赖
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 4. 编译为静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# 5. 创建最终镜像（更小）
FROM alpine

WORKDIR /root/
COPY --from=builder /scr/go/cp/main .

EXPOSE 23322
CMD ["./main"]
