package model

import (
	"time"
)

const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in-progress"
	TaskStatusDone       = "done"
)

type Task struct {
	ID                string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	Title             string     `gorm:"type:varchar(255);not null" json:"title"`
	Category          string     `gorm:"type:varchar(100);not null" json:"category"`
	Description       string     `gorm:"type:text" json:"description,omitempty"`
	Targets           string     `gorm:"type:text" json:"targets"`
	Links             string     `gorm:"type:text" json:"links"`
	Status            string     `gorm:"type:varchar(20);default:'todo';not null" json:"status"`
	Deadline          *time.Time `json:"deadline,omitempty"`
	ActualCompletedAt *time.Time `json:"actual_completed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Logs []TaskLog `gorm:"foreignKey:TaskID" json:"logs,omitempty"`
}

type TaskLog struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID       string    `gorm:"type:varchar(36);index;not null" json:"task_id"`
	LogText      string    `gorm:"type:text;not null" json:"log_text"`
	ProgressNote string    `gorm:"type:varchar(100)" json:"progress_note,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Request and Response DTOs

type CreateTaskReq struct {
	Title       string     `json:"title" binding:"required"`
	Category    string     `json:"category" binding:"required"`
	Description string     `json:"description"`
	Targets     string     `json:"targets"`
	Links       string     `json:"links"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateProgressReq struct {
	LogText    string     `json:"log_text" binding:"required"`
	MarkAsDone bool       `json:"mark_as_done"`
	NewStatus  string     `json:"new_status"`
	Deadline   *time.Time `json:"deadline"`
}

type UpdateTaskReq struct {
	Title       string     `json:"title"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Targets     string     `json:"targets"`
	Links       string     `json:"links"`
	Deadline    *time.Time `json:"deadline"`
}

type ActiveTaskResp struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

type DailySummaryActivity struct {
	TaskID    string   `json:"task_id"`
	TaskTitle string   `json:"task_title"`
	Category  string   `json:"category"`
	Status    string   `json:"status"`
	TodayLogs []string `json:"today_logs"`
}

type DailySummaryResp struct {
	Date       string                 `json:"date"`
	Activities []DailySummaryActivity `json:"activities"`
}

type StatsSummaryResp struct {
	TotalTasks         int                    `json:"total_tasks"`
	CompletedTasks     int                    `json:"completed_tasks"`
	TodoTasks          int                    `json:"todo_tasks"`
	InProgressTasks    int                    `json:"in_progress_tasks"`
	ByCategory         map[string]int         `json:"by_category"`
	CompletionRate     float64                `json:"completion_rate"`
	WeeklyStats        []DailyStats           `json:"weekly_stats"`
}

type DailyStats struct {
	Date       string `json:"date"`
	Completed  int    `json:"completed"`
	Created    int    `json:"created"`
}

type StandardResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func SuccessResp(data interface{}) StandardResponse {
	return StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

func ErrorResp(code int, msg string) StandardResponse {
	return StandardResponse{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
