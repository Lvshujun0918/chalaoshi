# 查老师 — 教师评分查询系统

一个面向ZJU的教师评分与课程查询平台，支持按教师姓名、院系、课程进行搜索、排序，查看教师详情、学生评论及课程评分。数据来自于査老师原版离线数据。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite |
| 后端 | Go (Gin + GORM) |
| 数据库 | SQLite |
| 反向代理 | nginx |
| 容器化 | Docker (多阶段构建) |
| CI/CD | GitHub Actions |

## 项目结构

```
chalaoshi/
├── src/                    # 前端源码 (Vue 3)
│   ├── App.vue             # 主应用组件
│   ├── components/         # UI 组件 (SearchBar, TeacherCard, CourseCard, CommentModal 等)
│   ├── composables/        # 组合式 API (useTeacherSearch, useCourseSearch)
│   ├── assets/             # 静态资源
│   └── main.js             # 应用入口
├── backend/                # 后端源码 (Go)
│   ├── main.go             # 服务入口，路由注册
│   ├── handlers/           # API 处理器
│   ├── models/             # 数据模型 (GORM)
│   ├── database/           # 数据库初始化与连接
│   ├── importer/           # CSV 数据导入
│   └── data/               # SQLite 数据库文件存放位置
├── data/                   # 原始数据 (CSV)
│   ├── teachers.csv        # 教师基础数据
│   ├── comment_*.csv       # 各院系评论数据
│   └── gpa.json            # GPA 数据
├── public/                 # 前端静态资源
├── docker-compose.yml      # Docker Compose 部署配置
├── Dockerfile              # 多阶段 Docker 构建
├── nginx.conf              # nginx 配置
├── vite.config.js          # Vite 构建配置
└── package.json            # 前端依赖
```

## 快速开始

### 环境要求

- Node.js ≥ 18
- Go ≥ 1.19
- SQLite 3

### 本地开发

**1. 启动后端**

```bash
cd backend
go mod download
go run main.go -port 8080 -db data/chalaoshi.db -data ../data
```

后端默认监听 `http://localhost:8080`，首次启动会自动从 `../data` 目录导入 CSV 数据到 SQLite。

**2. 启动前端**

```bash
npm install
npm run dev
```

前端开发服务器默认监听 `http://localhost:5173`，API 请求会自动代理到后端 `localhost:8080`。

### Docker 部署

使用 Docker Compose 一键部署：

```bash
docker compose up -d
```

或者手动构建运行：

```bash
# 构建镜像
docker build -t chalaoshi .

# 运行容器
docker run -d \
  -p 2382:80 \
  -v $(pwd)/backend/data/chalaoshi.db:/app/db/chalaoshi.db \
  chalaoshi
```

服务启动后访问 `http://localhost:2382`。

### GitHub Container Registry

项目已配置 GitHub Actions 自动构建并推送镜像到 `ghcr.io/lvshujun0918/chalaoshi:latest`。

---

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/teachers` | 教师列表（支持 `q` 搜索、`department` 筛选、`sort_by` / `sort_order` 排序、分页） |
| GET | `/api/teachers/:id` | 教师详情 |
| GET | `/api/teachers/:id/comments` | 教师评论列表 |
| GET | `/api/search` | 快捷搜索（输入提示） |
| GET | `/api/courses` | 课程列表（支持搜索、排序、分页） |
| GET | `/api/courses/search` | 课程快捷搜索 |
| GET | `/api/departments` | 院系列表 |
| GET | `/api/stats` | 平台统计数据 |
| GET | `/api/version` | 数据版本信息 |

---

## 功能特性

- **教师搜索** — 按姓名、拼音、院系搜索教师，支持评分排序
- **课程搜索** — 按课程名称搜索，查看课程评分
- **教师详情** — 查看教师基本信息、综合评分、热度
- **学生评论** — 浏览学生对教师的评论（含点赞/点踩）
- **院系筛选** — 按学院/部门筛选教师
- **分页浏览** — 支持大数据量下的分页加载
- **版本追踪** — 页面底部展示 Git Commit Hash 和数据版本
- **响应式设计** — 适配桌面与移动端

---

## 许可

MIT License
