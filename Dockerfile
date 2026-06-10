# ============================================
# chalaoshi — 多阶段 Docker 构建
#
# 阶段 1: Node.js 构建前端 (npm ci + vite build)
# 阶段 2: Go 构建后端 (CGO 静态链接, 多平台 QEMU 支持)
# 阶段 3: Alpine 最小运行镜像
#
# CI 加速策略：
#   - cache mount 避免重复下载 npm / Go 模块
#   - GHA cache (type=gha) 缓存 Docker 层跨工作流复用
#   - 依赖文件不变时 npm ci / go mod download 层直接命中
# ============================================

# ============================================
# 阶段 1：构建前端 (Node.js)
# ============================================
FROM node:22-alpine AS builder-node

WORKDIR /app

# GIT_HASH 由 CI --build-arg 传入；本地构建时 vite.config.js 自动 fallback 到 git 命令
ARG GIT_HASH
ENV GIT_HASH=${GIT_HASH}

# 先复制依赖描述文件 → Docker 层缓存
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci

COPY . .
RUN npm run build

# ============================================
# 阶段 2：构建后端 (Go)
# ============================================
FROM golang:1.23-bullseye AS builder-go

WORKDIR /app

# Go 模块依赖先复制 → 层缓存
COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY backend/ ./
COPY --from=builder-node /app/dist ./dist

# CGO 静态链接 + cache mount 加速
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 \
    go build -trimpath -ldflags='-s -w -linkmode external -extldflags "-static"' \
    -o chalaoshi-server .

# ============================================
# 阶段 3：运行阶段（最小镜像）
# ============================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder-go /app/chalaoshi-server .
COPY --from=builder-go /app/dist ./dist

RUN mkdir -p /app/db /app/data
COPY backend/data/chalaoshi.db /app/db/chalaoshi.db

# 复制数据版本文件（供 /api/version 读取）
COPY data/ver.json /app/data/ver.json

EXPOSE 8080
CMD ["./chalaoshi-server", "-db", "/app/db/chalaoshi.db", "-static", "./dist", "-port", "8080"]
