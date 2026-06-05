package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"chalaoshi-backend/database"
	"chalaoshi-backend/models"
)

const (
	teacherFile = "teachers.csv"
	commentDir  = "."
	gpaFile     = "gpa.json"
)

// ImportAll 导入所有数据
func ImportAll(dataDir string) error {
	log.Println("[Import] 开始导入数据...")

	// 1. 导入教师数据
	if err := importTeachers(filepath.Join(dataDir, teacherFile)); err != nil {
		return fmt.Errorf("导入教师数据失败: %w", err)
	}

	// 2. 导入评论数据
	if err := importComments(dataDir); err != nil {
		return fmt.Errorf("导入评论数据失败: %w", err)
	}

	// 3. 导入GPA数据（gpa.json 在项目根目录）
	gpaPath := filepath.Join(dataDir, "..", gpaFile)
	if err := importGPA(gpaPath); err != nil {
		// 也尝试在data目录下查找
		if err2 := importGPA(filepath.Join(dataDir, gpaFile)); err2 != nil {
			log.Printf("[Import] 导入GPA数据警告: %v (data目录: %v)", err, err2)
		}
	}

	// 4. 数据导入后无需额外索引操作

	log.Println("[Import] 数据导入完成!")
	return nil
}

func importTeachers(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// 跳过表头
	if _, err := reader.Read(); err != nil {
		return err
	}

	batch := make([]models.Teacher, 0, 500)
	total := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[Import] 读取教师行错误: %v", err)
			continue
		}
		if len(record) < 8 {
			continue
		}

		id, _ := strconv.ParseUint(record[0], 10, 64)
		hotness, _ := strconv.Atoi(record[3])
		ratingCount, _ := strconv.Atoi(record[4])
		rating, _ := strconv.ParseFloat(record[5], 64)

		batch = append(batch, models.Teacher{
			ID:          uint(id),
			Name:        record[1],
			Department:  record[2],
			Hotness:     hotness,
			RatingCount: ratingCount,
			Rating:      rating,
			Pinyin:      record[6],
			PinyinAbbr:  record[7],
		})

		if len(batch) >= 500 {
			if err := database.DB.Create(&batch).Error; err != nil {
				log.Printf("[Import] 批量插入教师错误: %v", err)
			}
			total += len(batch)
			batch = batch[:0]
		}
	}

	// 处理剩余
	if len(batch) > 0 {
		if err := database.DB.Create(&batch).Error; err != nil {
			log.Printf("[Import] 批量插入教师错误: %v", err)
		}
		total += len(batch)
	}

	log.Printf("[Import] 教师: %d 条", total)
	return nil
}

func importComments(dataDir string) error {
	files, err := filepath.Glob(filepath.Join(dataDir, "comment_*.csv"))
	if err != nil {
		return err
	}

	total := 0
	batch := make([]models.Comment, 0, 500)

	for _, filePath := range files {
		// 根据文件名推断院系（仅用于日志）
		// baseName := filepath.Base(filePath)
		// deptName := strings.TrimSuffix(strings.TrimPrefix(baseName, "comment_"), ".csv")

		f, err := os.Open(filePath)
		if err != nil {
			log.Printf("[Import] 打开评论文件失败 %s: %v", filePath, err)
			continue
		}

		reader := csv.NewReader(f)
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1 // 允许变长字段

		// 跳过表头
		if _, err := reader.Read(); err != nil {
			f.Close()
			continue
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				// 跳过有问题的行
				continue
			}
			if len(record) < 8 {
				continue
			}

			commentID, _ := strconv.ParseUint(record[0], 10, 64)
			teacherID, _ := strconv.ParseUint(record[1], 10, 64)
			netLikes, _ := strconv.Atoi(record[4])
			likes, _ := strconv.Atoi(record[5])
			dislikes, _ := strconv.Atoi(record[6])

			pubTime, err := time.Parse("2006-01-02 15:04:05", record[3])
			if err != nil {
				pubTime = time.Now()
			}

			batch = append(batch, models.Comment{
				ID:          uint(commentID),
				TeacherID:   uint(teacherID),
				TeacherName: record[2],
				PublishTime: pubTime,
				NetLikes:    netLikes,
				Likes:       likes,
				Dislikes:    dislikes,
				Content:     record[7],
			})

			if len(batch) >= 500 {
				if err := database.DB.Create(&batch).Error; err != nil {
					log.Printf("[Import] 批量插入评论错误: %v", err)
				}
				total += len(batch)
				batch = batch[:0]
			}
		}
		f.Close()
	}

	// 处理剩余
	if len(batch) > 0 {
		if err := database.DB.Create(&batch).Error; err != nil {
			log.Printf("[Import] 批量插入评论错误: %v", err)
		}
		total += len(batch)
	}

	log.Printf("[Import] 评论: %d 条", total)
	return nil
}

func importGPA(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// gpa.json 结构: { "TeacherName": [["courseName", "gpa", "count", "stdDev"], ...] }
	var raw map[string][][]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		// 尝试另一种格式
		var raw2 map[string][][]json.RawMessage
		if err2 := json.Unmarshal(data, &raw2); err2 != nil {
			return fmt.Errorf("解析GPA JSON失败: %w (尝试2: %v)", err, err2)
		}
		return importGPAFromRaw(raw2)
	}

	return importGPAFromInterface(raw)
}

func importGPAFromInterface(raw map[string][][]interface{}) error {
	batch := make([]models.Course, 0, 200)
	total := 0

	for teacherName, courses := range raw {
		for _, c := range courses {
			if len(c) < 4 {
				continue
			}
			courseName, _ := c[0].(string)
			gpa, _ := strconv.ParseFloat(fmt.Sprint(c[1]), 64)
			countStr := strings.TrimSuffix(fmt.Sprint(c[2]), "+")
			count, _ := strconv.Atoi(countStr)
			stdDev, _ := strconv.ParseFloat(fmt.Sprint(c[3]), 64)

			batch = append(batch, models.Course{
				TeacherName: teacherName,
				CourseName:  courseName,
				GPA:         gpa,
				Count:       count,
				StdDev:      stdDev,
			})

			if len(batch) >= 200 {
				if err := database.DB.Create(&batch).Error; err != nil {
					log.Printf("[Import] 批量插入GPA错误: %v", err)
				}
				total += len(batch)
				batch = batch[:0]
			}
		}
	}

	if len(batch) > 0 {
		if err := database.DB.Create(&batch).Error; err != nil {
			log.Printf("[Import] 批量插入GPA错误: %v", err)
		}
		total += len(batch)
	}

	log.Printf("[Import] 课程GPA: %d 条", total)
	return nil
}

func importGPAFromRaw(raw map[string][][]json.RawMessage) error {
	batch := make([]models.Course, 0, 200)
	total := 0

	for teacherName, courses := range raw {
		for _, c := range courses {
			if len(c) < 4 {
				continue
			}
			var courseName string
			json.Unmarshal(c[0], &courseName)

			var gpaStr string
			json.Unmarshal(c[1], &gpaStr)
			gpa, _ := strconv.ParseFloat(gpaStr, 64)

			var countStr string
			json.Unmarshal(c[2], &countStr)
			countStr = strings.TrimSuffix(countStr, "+")
			count, _ := strconv.Atoi(countStr)

			var stdDevStr string
			json.Unmarshal(c[3], &stdDevStr)
			stdDev, _ := strconv.ParseFloat(stdDevStr, 64)

			batch = append(batch, models.Course{
				TeacherName: teacherName,
				CourseName:  courseName,
				GPA:         gpa,
				Count:       count,
				StdDev:      stdDev,
			})

			if len(batch) >= 200 {
				if err := database.DB.Create(&batch).Error; err != nil {
					log.Printf("[Import] 批量插入GPA错误: %v", err)
				}
				total += len(batch)
				batch = batch[:0]
			}
		}
	}

	if len(batch) > 0 {
		if err := database.DB.Create(&batch).Error; err != nil {
			log.Printf("[Import] 批量插入GPA错误: %v", err)
		}
		total += len(batch)
	}

	log.Printf("[Import] 课程GPA: %d 条", total)
	return nil
}


