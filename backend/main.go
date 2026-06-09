package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"chalaoshi-backend/database"
	"chalaoshi-backend/handlers"
	"chalaoshi-backend/importer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 命令行参数
	port := flag.String("port", "8080", "监听端口")
	dbPath := flag.String("db", "data/chalaoshi.db", "SQLite数据库路径")
	dataDir := flag.String("data", "../data", "CSV数据目录")
	staticDir := flag.String("static", "../dist", "前端静态文件目录")
	flag.Parse()

	// 初始化数据库
	dbFullPath := *dbPath
	if !filepath.IsAbs(dbFullPath) {
		// 相对于可执行文件所在目录
		execDir, _ := os.Getwd()
		dbFullPath = filepath.Join(execDir, dbFullPath)
	}

	log.Printf("[Main] 数据库路径: %s", dbFullPath)
	if err := database.Init(dbFullPath); err != nil {
		log.Fatalf("[Main] 数据库初始化失败: %v", err)
	}

	// 检查是否需要导入数据
	if database.NeedsImport() {
		dataFullPath := *dataDir
		if !filepath.IsAbs(dataFullPath) {
			execDir, _ := os.Getwd()
			dataFullPath = filepath.Join(execDir, dataFullPath)
		}
		log.Printf("[Main] 数据目录: %s", dataFullPath)
		if err := importer.ImportAll(dataFullPath); err != nil {
			log.Fatalf("[Main] 数据导入失败: %v", err)
		}
	} else {
		log.Println("[Main] 数据库已存在数据，跳过导入")
	}

	// 初始化 Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS（允许前端跨域访问）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// ============ API 路由 ============
	api := r.Group("/api")
	{
		// 教师相关
		api.GET("/teachers", handlers.GetTeachers)           // 列表+搜索
		api.GET("/teachers/:id", handlers.GetTeacherDetail)  // 详情
		api.GET("/teachers/:id/comments", handlers.GetTeacherComments) // 评论

		// 搜索（简化版，用于输入提示）
		api.GET("/search", handlers.SearchTeachers)

		// 课程相关
		api.GET("/courses", handlers.GetCourses)                  // 列表+搜索+排序
		api.GET("/courses/search", handlers.SearchCourses)        // 快捷搜索

		// 院系列表
		api.GET("/departments", handlers.GetDepartments)

		// 统计
		api.GET("/stats", handlers.GetStats)
	}

	// 静态文件服务（生产环境下前端打包文件）
	staticFullPath := *staticDir
	if !filepath.IsAbs(staticFullPath) {
		execDir, _ := os.Getwd()
		staticFullPath = filepath.Join(execDir, staticFullPath)
	}
	r.Static("/assets", filepath.Join(staticFullPath, "assets"))
	r.StaticFile("/favicon.svg", filepath.Join(staticFullPath, "favicon.svg"))
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(staticFullPath, "index.html"))
	})

	log.Printf("[Main] 服务启动于 :%s", *port)
	if err := r.Run(":" + *port); err != nil {
		log.Fatalf("[Main] 服务启动失败: %v", err)
	}
}
