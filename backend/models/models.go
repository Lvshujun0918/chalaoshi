package models

import "time"

// Teacher 教师信息
type Teacher struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"type:varchar(128);index:idx_name" json:"name"`
	Department  string  `gorm:"type:varchar(128);index:idx_department" json:"department"`
	Hotness     int     `gorm:"default:0;index:idx_hotness" json:"hotness"`
	RatingCount int     `gorm:"default:0" json:"rating_count"`
	Rating      float64 `gorm:"type:real;default:0;index:idx_rating" json:"rating"`
	Pinyin      string  `gorm:"type:varchar(256);index:idx_pinyin" json:"pinyin"`
	PinyinAbbr  string  `gorm:"type:varchar(64);index:idx_pinyin_abbr" json:"pinyin_abbr"`
}

func (Teacher) TableName() string {
	return "teachers"
}

// Comment 评论信息
type Comment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TeacherID   uint      `gorm:"index:idx_comment_teacher" json:"teacher_id"`
	TeacherName string    `gorm:"type:varchar(128)" json:"teacher_name"`
	PublishTime time.Time `gorm:"index:idx_comment_time" json:"publish_time"`
	NetLikes    int       `gorm:"default:0" json:"net_likes"`
	Likes       int       `gorm:"default:0" json:"likes"`
	Dislikes    int       `gorm:"default:0" json:"dislikes"`
	Content     string    `gorm:"type:text" json:"content"`
}

func (Comment) TableName() string {
	return "comments"
}

// Course 课程GPA信息
type Course struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	TeacherName string  `gorm:"type:varchar(128);index:idx_course_teacher" json:"teacher_name"`
	CourseName  string  `gorm:"type:varchar(256)" json:"course_name"`
	GPA         float64 `gorm:"type:real" json:"gpa"`
	Count       int     `json:"count"`
	StdDev      float64 `gorm:"type:real" json:"std_dev"`
}

func (Course) TableName() string {
	return "courses"
}

// --- API response types ---

// TeacherListResponse 教师列表响应
type TeacherListResponse struct {
	Teachers   []Teacher `json:"teachers"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	Comments   []Comment `json:"comments"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// TeacherDetailResponse 教师详情响应
type TeacherDetailResponse struct {
	Teacher  Teacher   `json:"teacher"`
	Courses  []Course  `json:"courses"`
	Comments []Comment `json:"comments,omitempty"`
}

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Query      string `form:"q"`
	Department string `form:"department"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	SortBy     string `form:"sort_by"`   // rating, hotness, name
	SortOrder  string `form:"sort_order"` // asc, desc
}

// DepartmentItem 院系列表项
type DepartmentItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// StatsResponse 首页统计数据
type StatsResponse struct {
	TotalTeachers   int64 `json:"total_teachers"`
	TotalComments   int64 `json:"total_comments"`
	TotalDepartments int64 `json:"total_departments"`
}
