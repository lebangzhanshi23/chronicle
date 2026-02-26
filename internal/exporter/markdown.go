package exporter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

type LogView struct {
	Time string
	Text string
	Note string
}

type TaskView struct {
	Title              string
	Category           string
	Status             string
	Description        string
	Targets            string
	CreatedAt          string
	CompletedAt        string
	Deadline           string
	LogsByDate         map[string][]LogView
	ReverseSortedDates []string
}

func GenerateDailyMarkdown(dateStr string) ([]byte, error) {
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

	// query tasks completed today
	var tasks []model.Task
	if err := service.DB.Where("status = ? AND actual_completed_at >= ? AND actual_completed_at <= ?", model.TaskStatusDone, startOfDay, endOfDay).Find(&tasks).Error; err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles("templates/obsidian_task.tmpl")
	if err != nil {
		return nil, err
	}

	// Create a zip buffer
	var zipBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuf)

	for _, t := range tasks {
		var logs []model.TaskLog
		if err := service.DB.Where("task_id = ?", t.ID).Order("created_at desc").Find(&logs).Error; err != nil {
			return nil, err
		}

		logsByDate := make(map[string][]LogView)
		for _, l := range logs {
			d := l.CreatedAt.Format("2006-01-02")
			logsByDate[d] = append(logsByDate[d], LogView{
				Time: l.CreatedAt.Format("15:04"),
				Text: l.LogText,
				Note: l.ProgressNote,
			})
		}

		completedAt := ""
		if t.ActualCompletedAt != nil {
			completedAt = t.ActualCompletedAt.Format("2006-01-02 15:04:05")
		}

		deadlineAt := ""
		if t.Deadline != nil {
			deadlineAt = t.Deadline.Format("2006-01-02 15:04:05")
		}

		tv := TaskView{
			Title:       t.Title,
			Category:    t.Category,
			Status:      t.Status,
			Description: t.Description,
			Targets:     t.Targets,
			CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
			CompletedAt: completedAt,
			Deadline:    deadlineAt,
			LogsByDate:  logsByDate,
		}

		var dates []string
		for k := range logsByDate {
			dates = append(dates, k)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(dates)))
		tv.ReverseSortedDates = dates

		var taskBuf bytes.Buffer
		if err := tmpl.Execute(&taskBuf, tv); err != nil {
			continue
		}

		safeTitle := strings.ReplaceAll(t.Title, "/", "-")
		safeTitle = strings.ReplaceAll(safeTitle, "\\", "-")
		fWriter, err := zipWriter.Create(fmt.Sprintf("%s.md", safeTitle))
		if err != nil {
			continue
		}
		fWriter.Write(taskBuf.Bytes())
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return zipBuf.Bytes(), nil
}
