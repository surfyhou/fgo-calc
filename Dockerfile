# ---- build stage ----
FROM golang:1.25-alpine AS build

# 国内服务器走 goproxy.cn 加速；proxy.golang.org 也可用
ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /src

# 先拉依赖，利用层缓存
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 再拷源码编译
COPY backend/ ./
RUN go build -trimpath -ldflags="-s -w" -o /out/fgo-calc-backend .

# ---- runtime stage ----
FROM alpine:3.20

# 时区（可选）+ 基础证书
RUN apk add --no-cache tzdata ca-certificates

# 目录结构需匹配后端的相对路径约定：
#   工作目录 = /app/backend，二进制内引用 ./static 与 ../data
WORKDIR /app/backend

COPY --from=build /out/fgo-calc-backend ./fgo-calc-backend
COPY backend/static ./static
COPY backend/config.prod.json ./config.prod.json
COPY data /app/data

EXPOSE 30005

CMD ["./fgo-calc-backend", "-config", "config.prod.json"]
