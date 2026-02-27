package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yuyudeqiu/chronicle/internal/model"
	"gorm.io/gorm"
)

func CreateTask(req model.CreateTaskReq) (*model.Task, error) {
	var localDeadline *time.Time
	if req.Deadline != nil {
		ld := req.Deadline.Local()
		localDeadline = &ld
	}

	task := &model.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Category:    req.Category,
		Description: req.Description,
		Targets:     req.Targets,
		Deadline:    localDeadline,
		Status:      model.TaskStatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func GetActiveTasks() ([]model.ActiveTaskResp, error) {
	var tasks []model.ActiveTaskResp
	err := DB.Model(&model.Task{}).Select("id", "title", "category", "status").
		Where("status IN ?", []string{model.TaskStatusTodo, model.TaskStatusInProgress}).
		Order("created_at desc").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetHistoryTasks() ([]model.ActiveTaskResp, error) {
	var tasks []model.ActiveTaskResp
	err := DB.Model(&model.Task{}).Select("id", "title", "category", "status").
		Where("status = ?", model.TaskStatusDone).
		Order("actual_completed_at desc").
		Limit(50).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(id string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Delete logs first
		if err := tx.Where("task_id = ?", id).Delete(&model.TaskLog{}).Error; err != nil {
			return err
		}
		// Delete task
		if err := tx.Where("id = ?", id).Delete(&model.Task{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func DeleteWorklog(id string) error {
	return DB.Where("id = ?", id).Delete(&model.TaskLog{}).Error
}

func GetTask(id string) (*model.Task, error) {
	var task model.Task
	if err := DB.Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateProgress(taskID string, req model.UpdateProgressReq) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var task model.Task
		if err := tx.First(&task, "id = ?", taskID).Error; err != nil {
			return err
		}

		if task.Status == model.TaskStatusDone {
			return errors.New("task is already done")
		}

		logEntry := model.TaskLog{
			ID:        uuid.New().String(),
			TaskID:    taskID,
			LogText:   req.LogText,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&logEntry).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"updated_at": time.Now(),
		}

		// Update deadline if provided
		if req.Deadline != nil {
			updates["deadline"] = req.Deadline.Local()
		}

		newStatus := task.Status
		if req.MarkAsDone || req.NewStatus == model.TaskStatusDone {
			newStatus = model.TaskStatusDone
			now := time.Now()
			updates["actual_completed_at"] = now
		} else if req.NewStatus != "" && req.NewStatus != task.Status {
			// e.g. move to in-progress
			newStatus = req.NewStatus
		}

		if newStatus != task.Status {
			updates["status"] = newStatus
		}

		if err := tx.Model(&task).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetDailySummary returns the daily summary view for a given date string (YYYY-MM-DD).
// Active tasks and tasks done on this day. Wait, design spec says daily summary
// returns what the agent did today.
// In the design for GetDailySummary:
// { "task_id": "xxx", "task_title": "...", "status": "done", "today_logs": ["10:12 - ...", ...] }
func GetDailySummary(dateStr string) (*model.DailySummaryResp, error) {
	// Parse date
	var targetDate time.Time
	var err error
	if dateStr == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
	}

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	var logs []model.TaskLog
	// Find all logs for the target date
	if err := DB.Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Order("created_at asc").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	taskLogMap := make(map[string][]string)
	var taskIDs []string
	for _, l := range logs {
		timeStr := l.CreatedAt.Format("15:04")
		logLine := timeStr + " - " + l.LogText
		if len(taskLogMap[l.TaskID]) == 0 {
			taskIDs = append(taskIDs, l.TaskID)
		}
		taskLogMap[l.TaskID] = append(taskLogMap[l.TaskID], logLine)
	}

	var activities []model.DailySummaryActivity
	if len(taskIDs) > 0 {
		var tasks []model.Task
		if err := DB.Where("id IN ?", taskIDs).Find(&tasks).Error; err != nil {
			return nil, err
		}
		taskMap := make(map[string]model.Task)
		for _, t := range tasks {
			taskMap[t.ID] = t
		}

		for _, tid := range taskIDs {
			t := taskMap[tid]
			activities = append(activities, model.DailySummaryActivity{
				TaskID:    t.ID,
				TaskTitle: t.Title,
				Status:    t.Status,
				TodayLogs: taskLogMap[tid],
			})
		}
	}

	if activities == nil {
		activities = []model.DailySummaryActivity{}
	}

	resp := &model.DailySummaryResp{
		Date:       startOfDay.Format("2006-01-02"),
		Activities: activities,
	}

	return resp, nil
}
