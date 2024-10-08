# 使用官方 Go 镜像作为基础
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod ./
# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o my-go-app .

# 创建一个轻量级的运行时镜像
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/my-go-app .
CMD ["./my-go-app"]