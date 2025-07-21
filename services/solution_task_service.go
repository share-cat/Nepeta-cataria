package services

import (
	"errors"
	"fmt"
	"strings"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/pkg/catwheel"
	"share-cat/Nepeta-cataria/repositories"
)

type SolutionTaskService struct {
	repo         *repositories.SolutionTaskRepository
	solutionRepo *repositories.SolutionRepository
}

func NewSolutionTaskService(
	repo *repositories.SolutionTaskRepository,
	solutionRepo *repositories.SolutionRepository,
) *SolutionTaskService {
	return &SolutionTaskService{
		repo:         repo,
		solutionRepo: solutionRepo,
	}
}

// CreateTask 創建新的解答任務
func (s *SolutionTaskService) CreateTask(task models.SolutionTask) (models.SolutionTask, error) {
	// 檢查 Solution 是否存在
	_, err := s.solutionRepo.GetByID(task.SolutionID)
	if err != nil {
		return models.SolutionTask{}, errors.New("解答不存在")
	}

	return s.repo.Create(task)
}

// CreateTasksForSolution 為解答創建多個任務
func (s *SolutionTaskService) CreateTasksForSolution(solutionID uint, taskCount int) ([]models.SolutionTask, error) {
	// 檢查 Solution 是否存在
	_, err := s.solutionRepo.GetByID(solutionID)
	if err != nil {
		return nil, errors.New("解答不存在")
	}

	var tasks []models.SolutionTask
	for i := 0; i < taskCount; i++ {
		task := models.SolutionTask{
			SolutionID: solutionID,
			TaskIndex:  i,
			Status:     "pending",
		}
		tasks = append(tasks, task)
	}

	err = s.repo.CreateBatch(tasks)
	if err != nil {
		return nil, err
	}

	return s.repo.GetBySolutionID(solutionID)
}

// GetTaskByID 根據ID獲取任務
func (s *SolutionTaskService) GetTaskByID(id uint) (models.SolutionTask, error) {
	return s.repo.GetByID(id)
}

// GetTasksBySolutionID 根據Solution ID獲取所有任務
func (s *SolutionTaskService) GetTasksBySolutionID(solutionID uint) ([]models.SolutionTask, error) {
	return s.repo.GetBySolutionID(solutionID)
}

// UpdateTaskFromExecution 從執行結果更新任務
func (s *SolutionTaskService) UpdateTaskFromExecution(result catwheel.ExecuteResponse, taskID uint) (models.SolutionTask, error) {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return models.SolutionTask{}, errors.New("任務不存在")
	}

	task.Status = result.Result.Status
	task.ExecutionTime = result.Result.ExecutionTime
	task.MemoryUsage = result.Result.MemoryUsed
	task.Stdout = result.Result.Stdout
	task.Stderr = result.Result.Stderr
	task.ExitCode = result.Result.ExitCode

	err = s.repo.Update(task)
	if err != nil {
		return models.SolutionTask{}, err
	}

	return task, nil
}

// UpdateTaskStatus 更新任務狀態
func (s *SolutionTaskService) UpdateTaskStatus(id uint, status string, executionTime, memoryUsage int64, score int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("任務不存在")
	}

	status = strings.ToLower(status)

	validStatus := map[string]bool{
		"pending":               true,
		"running":               true,
		"accepted":              true,
		"wrong_answer":          true,
		"time_limit_exceeded":   true,
		"memory_limit_exceeded": true,
		"compilation_error":     true,
		"runtime_error":         true,
		"success":               true,
	}

	if !validStatus[status] {
		fmt.Printf("無效的狀態: %s\n", status)
		return errors.New("無效的狀態")
	}

	return s.repo.UpdateStatus(id, status, executionTime, memoryUsage, score)
}

// DeleteTask 刪除任務
func (s *SolutionTaskService) DeleteTask(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("任務不存在")
	}

	return s.repo.Delete(id)
}

// DeleteTasksBySolutionID 刪除解答的所有任務
func (s *SolutionTaskService) DeleteTasksBySolutionID(solutionID uint) error {
	return s.repo.DeleteBySolutionID(solutionID)
}

// GetSolutionSummary 獲取解答的匯總資訊
func (s *SolutionTaskService) GetSolutionSummary(solutionID uint) (map[string]interface{}, error) {
	tasks, err := s.repo.GetBySolutionID(solutionID)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return map[string]interface{}{
			"total_tasks":          0,
			"completed_tasks":      0,
			"total_execution_time": int64(0),
			"total_memory_usage":   int64(0),
			"total_score":          0,
			"overall_status":       "pending",
		}, nil
	}

	var totalExecutionTime int64 = 0
	var totalMemoryUsage int64 = 0
	var totalScore int = 0
	var completedTasks int = 0
	var acceptedTasks int = 0

	for _, task := range tasks {
		totalExecutionTime += task.ExecutionTime
		totalMemoryUsage += task.MemoryUsage
		totalScore += task.Score

		if task.Status != "pending" && task.Status != "running" {
			completedTasks++
		}
		if task.Status == "accepted" {
			acceptedTasks++
		}
	}

	overallStatus := "pending"
	if completedTasks == len(tasks) {
		if acceptedTasks == len(tasks) {
			overallStatus = "accepted"
		} else {
			overallStatus = "partial_accepted"
		}
	} else if completedTasks > 0 {
		overallStatus = "running"
	}

	return map[string]interface{}{
		"total_tasks":          len(tasks),
		"completed_tasks":      completedTasks,
		"accepted_tasks":       acceptedTasks,
		"total_execution_time": totalExecutionTime,
		"total_memory_usage":   totalMemoryUsage,
		"total_score":          totalScore,
		"overall_status":       overallStatus,
	}, nil
}
