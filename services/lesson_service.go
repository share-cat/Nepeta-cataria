package services

import (
	"errors"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/repositories"
)

type LessonService struct {
	lessonRepo *repositories.LessonRepository
	userRepo   *repositories.UserRepository
}

func NewLessonService(lessonRepo *repositories.LessonRepository, userRepo *repositories.UserRepository) *LessonService {
	return &LessonService{
		lessonRepo: lessonRepo,
		userRepo:   userRepo,
	}
}

// 創建課程
func (s *LessonService) CreateLesson(lesson *models.Lesson, teacherID uint) error {

	teacher, err := s.userRepo.GetByID(teacherID)
	if err != nil {
		return errors.New("teacher not found")
	}

	if teacher.Permission == nil || !teacher.Permission.IsTeacher {
		return errors.New("user is not a teacher")
	}

	lesson.TeacherId = int(teacherID)
	return s.lessonRepo.Create(lesson)
}

// 獲取課程詳情
func (s *LessonService) GetLessonByID(id uint) (*models.Lesson, error) {
	return s.lessonRepo.GetByID(id)
}

// TODO: 之後切換使用permission service 驗證
// 獲取教師的課程
func (s *LessonService) GetTeacherLessons(teacherID uint) ([]models.Lesson, error) {

	teacher, err := s.userRepo.GetByID(teacherID)
	if err != nil {
		return nil, errors.New("teacher not found")
	}

	if teacher.Permission == nil || !teacher.Permission.IsTeacher {
		return nil, errors.New("user is not a teacher")
	}

	return s.lessonRepo.GetByTeacherID(teacherID)
}

// 獲取學生的課程
func (s *LessonService) GetStudentLessons(studentID uint) ([]models.Lesson, error) {

	student, err := s.userRepo.GetByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	if student.Permission == nil || !student.Permission.IsStudent {
		return nil, errors.New("user is not a student")
	}

	return s.lessonRepo.GetByStudentID(studentID)
}

// 更新課程
func (s *LessonService) UpdateLesson(lesson *models.Lesson, teacherID uint) error {

	existingLesson, err := s.lessonRepo.GetByID(lesson.ID)
	if err != nil {
		return errors.New("lesson not found")
	}

	if existingLesson.TeacherId != int(teacherID) {
		return errors.New("permission denied: not the lesson creator")
	}

	return s.lessonRepo.Update(lesson)
}

// 刪除課程
func (s *LessonService) DeleteLesson(id, teacherID uint) error {

	lesson, err := s.lessonRepo.GetByID(id)
	if err != nil {
		return errors.New("lesson not found")
	}

	if lesson.TeacherId != int(teacherID) {
		return errors.New("permission denied: not the lesson creator")
	}

	return s.lessonRepo.Delete(id)
}

// 學生加入課程
func (s *LessonService) JoinLesson(lessonID, studentID uint) error {

	student, err := s.userRepo.GetByID(studentID)
	if err != nil {
		return errors.New("student not found")
	}

	if student.Permission == nil || !student.Permission.IsStudent {
		return errors.New("user is not a student")
	}

	enrolled, err := s.lessonRepo.IsStudentEnrolled(lessonID, studentID)
	if err != nil {
		return err
	}

	if enrolled {
		return errors.New("student already enrolled in this lesson")
	}

	return s.lessonRepo.AddStudent(lessonID, studentID)
}

// 學生退出課程
func (s *LessonService) LeaveLesson(lessonID, studentID uint) error {
	enrolled, err := s.lessonRepo.IsStudentEnrolled(lessonID, studentID)
	if err != nil {
		return err
	}

	if !enrolled {
		return errors.New("student not enrolled in this lesson")
	}

	return s.lessonRepo.RemoveStudent(lessonID, studentID)
}

// 獲取所有課程（分頁）
func (s *LessonService) GetAllLessons(page, pageSize int) ([]models.Lesson, int64, error) {
	offset := (page - 1) * pageSize
	lessons, err := s.lessonRepo.GetAll(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.lessonRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}
