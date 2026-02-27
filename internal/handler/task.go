package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuyudeqiu/chronicle/internal/exporter"
	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", CreateTask)
		v1.GET("/tasks", GetActiveTasks)
		v1.GET("/tasks/:id", GetTask)
		v1.DELETE("/tasks/:id", DeleteTask)
		v1.POST("/tasks/:id/progress", UpdateProgress)
		v1.DELETE("/worklogs/:id", DeleteWorklog)
		v1.GET("/reports/daily-summary", GetDailySummary)
		v1.GET("/exports/daily-markdown", GetDailyMarkdown)
	}
}

func CreateTask(c *gin.Context) {
	var req model.CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp(400, "invalid parameters: "+err.Error()))
		return
	}

	task, err := service.CreateTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to create task: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResp(task))
}

func GetActiveTasks(c *gin.Context) {
	status := c.Query("status")
	if status == "" || status == "in-progress,todo" || status == "todo,in-progress" {
		tasks, err := service.GetActiveTasks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to get tasks: "+err.Error()))
			return
		}
		c.JSON(http.StatusOK, model.SuccessResp(tasks))
		return
	}

	if status == "done" {
		tasks, err := service.GetHistoryTasks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to get tasks: "+err.Error()))
			return
		}
		c.JSON(http.StatusOK, model.SuccessResp(tasks))
		return
	}

	// In case agent asks for other statuses, this is not fully implemented per spec,
	// we just ignore it or return bad request since spec says active tasks API is the agent's eyes.
	c.JSON(http.StatusBadRequest, model.ErrorResp(400, "unsupported status query"))
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResp(400, "missing task id"))
		return
	}

	if err := service.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to delete task: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResp(nil))
}

func DeleteWorklog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResp(400, "missing worklog id"))
		return
	}

	if err := service.DeleteWorklog(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to delete worklog: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResp(nil))
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResp(400, "missing task id"))
		return
	}

	task, err := service.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to get task: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResp(task))
}

func UpdateProgress(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResp(400, "missing task id"))
		return
	}

	var req model.UpdateProgressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp(400, "invalid parameters: "+err.Error()))
		return
	}

	if err := service.UpdateProgress(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to update progress: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResp(nil))
}

func GetDailySummary(c *gin.Context) {
	dateStr := c.Query("date") // Format: YYYY-MM-DD
	summary, err := service.GetDailySummary(dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to get summary: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResp(summary))
}

func GetDailyMarkdown(c *gin.Context) {
	dateStr := c.Query("date") // Format: YYYY-MM-DD
	zipBytes, err := exporter.GenerateDailyMarkdown(dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp(500, "failed to generate markdown: "+err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename=obsidian_tasks.zip")
	c.Data(http.StatusOK, "application/zip", zipBytes)
}
