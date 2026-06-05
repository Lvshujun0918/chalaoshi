# ============================================
# 阶段 1：构建前端 (Node.js)
# ============================================
FROM node:22-alpine AS builder-node

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# ============================================
# 阶段 2：构建后端 (Go)
# ============================================
FROM golang:1.21-alpine AS builder-go

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
# 将前端构建产物复制到 Go 能找到的位置
COPY --from=builder-node /app/dist ./dist

RUN CGO_ENABLED=1 go build -o chalaoshi-server .

# ============================================
# 阶段 3：运行阶段
# ============================================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制 Go 二进制
COPY --from=builder-go /app/chalaoshi-server .

# 复制前端静态文件
COPY --from=builder-go /app/dist ./dist

# 复制数据文件（镜像内预置，运行时也可挂载覆盖）
COPY data/ ./data/

# 数据目录（SQLite 数据库存这里）
RUN mkdir -p /app/db

EXPOSE 8080

# 启动服务器，静态文件目录设为 ./dist
CMD ["./chalaoshi-server", "-data", "./data", "-db", "/app/db/chalaoshi.db", "-static", "./dist", "-port", "8080"]
