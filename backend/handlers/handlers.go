package handlers

import (
	"math"
	"net/http"
	"strconv"

	"chalaoshi-backend/database"
	"chalaoshi-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTeachers 获取教师列表（分页+搜索+筛选+排序）
func GetTeachers(c *gin.Context) {
	query := c.Query("q")
	department := c.Query("department")
	sortBy := c.DefaultQuery("sort_by", "rating")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 1 }
	if pageSize > 100 { pageSize = 100 }

	var dbQuery *gorm.DB

	// LIKE 搜索（对于10k数据量已足够快，配合索引）
	if query != "" {
		like := "%" + query + "%"
		dbQuery = database.DB.Model(&models.Teacher{}).Where(
			"name LIKE ? OR pinyin LIKE ? OR pinyin_abbr LIKE ?",
			like, like, like,
		)
	} else {
		dbQuery = database.DB.Model(&models.Teacher{})
	}

	// 院系筛选
	if department != "" {
		dbQuery = dbQuery.Where("department = ?", department)
	}

	// 计数
	var total int64
	dbQuery.Count(&total)

	// 排序
	orderCol := "rating"
	switch sortBy {
	case "hotness":
		orderCol = "hotness"
	case "name":
		orderCol = "name"
	case "rating_count":
		orderCol = "rating_count"
	}
	if sortOrder == "asc" {
		orderCol += " ASC"
	} else {
		orderCol += " DESC"
	}
	// 二级排序：同名次按姓名
	orderCol += ", name ASC"

	// 分页
	offset := (page - 1) * pageSize
	var teachers []models.Teacher
	dbQuery.Order(orderCol).Offset(offset).Limit(pageSize).Find(&teachers)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.TeacherListResponse{
		Teachers:   teachers,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetTeacherDetail 获取教师详情
func GetTeacherDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的教师ID"})
		return
	}

	var teacher models.Teacher
	if err := database.DB.First(&teacher, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "教师未找到"})
		return
	}

	// 获取课程信息
	var courses []models.Course
	database.DB.Where("teacher_name = ?", teacher.Name).Find(&courses)

	// 获取前5条评论预览
	var comments []models.Comment
	database.DB.Where("teacher_id = ?", id).
		Order("publish_time DESC").
		Limit(5).
		Find(&comments)

	c.JSON(http.StatusOK, models.TeacherDetailResponse{
		Teacher:  teacher,
		Courses:  courses,
		Comments: comments,
	})
}

// GetTeacherComments 获取教师评论
func GetTeacherComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的教师ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "time") // time, likes

	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 1 }
	if pageSize > 50 { pageSize = 50 }

	orderCol := "publish_time DESC"
	if sortBy == "likes" {
		orderCol = "likes DESC"
	} else if sortBy == "net_likes" {
		orderCol = "net_likes DESC"
	}

	var total int64
	database.DB.Model(&models.Comment{}).Where("teacher_id = ?", id).Count(&total)

	offset := (page - 1) * pageSize
	var comments []models.Comment
	database.DB.Where("teacher_id = ?", id).
		Order(orderCol).
		Offset(offset).
		Limit(pageSize).
		Find(&comments)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.CommentListResponse{
		Comments:   comments,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetDepartments 获取所有院系列表（含教师数）
func GetDepartments(c *gin.Context) {
	var results []models.DepartmentItem
	database.DB.Model(&models.Teacher{}).
		Select("department as name, count(*) as count").
		Group("department").
		Order("count DESC").
		Scan(&results)

	c.JSON(http.StatusOK, results)
}

// GetStats 获取统计数据
func GetStats(c *gin.Context) {
	var teacherCount, commentCount, deptCount int64
	database.DB.Model(&models.Teacher{}).Count(&teacherCount)
	database.DB.Model(&models.Comment{}).Count(&commentCount)
	database.DB.Model(&models.Teacher{}).Distinct("department").Count(&deptCount)
	// 用 count distinct
	database.DB.Raw("SELECT count(DISTINCT department) FROM teachers").Scan(&deptCount)

	c.JSON(http.StatusOK, models.StatsResponse{
		TotalTeachers:    teacherCount,
		TotalComments:    commentCount,
		TotalDepartments: deptCount,
	})
}

// SearchTeachers 快捷搜索（用于搜索框自动补全）
func SearchTeachers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, []models.Teacher{})
		return
	}

	like := "%" + query + "%"
	var teachers []models.Teacher
	database.DB.Model(&models.Teacher{}).
		Where("name LIKE ? OR pinyin LIKE ? OR pinyin_abbr LIKE ?", like, like, like).
		Order("rating DESC").
		Limit(10).
		Find(&teachers)

	c.JSON(http.StatusOK, teachers)
}

// GetCourses 获取课程列表（分页+搜索+排序）
func GetCourses(c *gin.Context) {
	query := c.Query("q")
	teacher := c.Query("teacher")
	sortBy := c.DefaultQuery("sort_by", "gpa")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 1 }
	if pageSize > 100 { pageSize = 100 }

	var dbQuery *gorm.DB

	if query != "" {
		like := "%" + query + "%"
		dbQuery = database.DB.Model(&models.Course{}).Where(
			"course_name LIKE ? OR teacher_name LIKE ?",
			like, like,
		)
	} else {
		dbQuery = database.DB.Model(&models.Course{})
	}

	// 按教师筛选
	if teacher != "" {
		dbQuery = dbQuery.Where("teacher_name = ?", teacher)
	}

	// 计数
	var total int64
	dbQuery.Count(&total)

	// 排序
	orderCol := "gpa"
	switch sortBy {
	case "count":
		orderCol = "count"
	case "course_name":
		orderCol = "course_name"
	case "teacher_name":
		orderCol = "teacher_name"
	case "std_dev":
		orderCol = "std_dev"
	}
	if sortOrder == "asc" {
		orderCol += " ASC"
	} else {
		orderCol += " DESC"
	}
	orderCol += ", course_name ASC"

	// 分页
	offset := (page - 1) * pageSize
	var courses []models.Course
	dbQuery.Order(orderCol).Offset(offset).Limit(pageSize).Find(&courses)

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.CourseListResponse{
		Courses:    courses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// SearchCourses 快捷搜索课程（用于搜索框自动补全）
func SearchCourses(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, []models.Course{})
		return
	}

	like := "%" + query + "%"
	var courses []models.Course
	database.DB.Model(&models.Course{}).
		Where("course_name LIKE ? OR teacher_name LIKE ?", like, like).
		Order("gpa DESC").
		Limit(10).
		Find(&courses)

	c.JSON(http.StatusOK, courses)
}
