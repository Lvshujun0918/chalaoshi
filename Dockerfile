# ---- 构建阶段 ----
FROM node:22-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存层（需要 devDependencies 才能构建）
COPY package.json package-lock.json ./
RUN npm ci

# 复制源码并构建
COPY . .
RUN npm run build

# ---- 运行阶段 ----
FROM nginx:stable-alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
