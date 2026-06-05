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

# 数据库目录
RUN mkdir -p /app/db

# 复制预构建的 SQLite 数据库（无需 CSV 导入，直接使用）
COPY backend/data/chalaoshi.db /app/db/chalaoshi.db

EXPOSE 8080

# 启动服务器（数据库已预构建，不需要 -data 导入）
CMD ["./chalaoshi-server", "-db", "/app/db/chalaoshi.db", "-static", "./dist", "-port", "8080"]
