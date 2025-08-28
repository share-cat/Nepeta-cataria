package repositories

import (
	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

type LessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// 創建課程
func (r *LessonRepository) Create(lesson *models.Lesson) error {
	return r.db.Create(lesson).Error
}

// 根據 ID 查詢課程
func (r *LessonRepository) GetByID(id uint) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.Preload("Teacher").Preload("Students").First(&lesson, id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// 根據教師 ID 查詢課程
func (r *LessonRepository) GetByTeacherID(teacherID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Where("teacher_id = ?", teacherID).
		Preload("Teacher").
		Preload("Students").
		Find(&lessons).Error
	return lessons, err
}

// 根據學生 ID 查詢加入的課程
func (r *LessonRepository) GetByStudentID(studentID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Joins("JOIN user_lessons ON user_lessons.lesson_id = lessons.id").
		Where("user_lessons.user_id = ?", studentID).
		Preload("Teacher").
		Preload("Students").
		Find(&lessons).Error
	return lessons, err
}

// 更新課程
func (r *LessonRepository) Update(lesson *models.Lesson) error {
	return r.db.Save(lesson).Error
}

// 刪除課程
func (r *LessonRepository) Delete(id uint) error {
	return r.db.Delete(&models.Lesson{}, id).Error
}

// 學生加入課程
func (r *LessonRepository) AddStudent(lessonID, studentID uint) error {
	var lesson models.Lesson
	var student models.User

	if err := r.db.First(&lesson, lessonID).Error; err != nil {
		return err
	}

	if err := r.db.First(&student, studentID).Error; err != nil {
		return err
	}

	return r.db.Model(&lesson).Association("Students").Append(&student)
}

// 學生退出課程
func (r *LessonRepository) RemoveStudent(lessonID, studentID uint) error {
	var lesson models.Lesson
	var student models.User

	if err := r.db.First(&lesson, lessonID).Error; err != nil {
		return err
	}

	if err := r.db.First(&student, studentID).Error; err != nil {
		return err
	}

	return r.db.Model(&lesson).Association("Students").Delete(&student)
}

// 檢查學生是否已加入課程
func (r *LessonRepository) IsStudentEnrolled(lessonID, studentID uint) (bool, error) {
	var count int64
	err := r.db.Table("user_lessons").
		Where("lesson_id = ? AND user_id = ?", lessonID, studentID).
		Count(&count).Error
	return count > 0, err
}

// 獲取所有課程（分頁）
func (r *LessonRepository) GetAll(offset, limit int) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Offset(offset).Limit(limit).
		Preload("Teacher").
		Preload("Students").
		Find(&lessons).Error
	return lessons, err
}

// 統計課程總數
func (r *LessonRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Lesson{}).Count(&count).Error
	return count, err
}
