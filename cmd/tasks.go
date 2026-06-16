package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

// Task Flags
var (
	category string
	desc     string
	links    string
	targets  string
	deadline string
	status   string
)

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		deadlineTime := parseDeadline(deadline)

		req := model.CreateTaskReq{
			Title:       title,
			Category:    category,
			Description: desc,
			Targets:     targets,
			Links:       links,
			Deadline:    deadlineTime,
		}

		task, err := service.CreateTask(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(task)
		} else {
			fmt.Printf("Task created: %s\n", task.ID)
			printTask(task)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		var tasks []model.ActiveTaskResp
		var err error

		if len(args) > 0 {
			queryStatus := args[0]
			if queryStatus == "done" {
				tasks, err = service.GetHistoryTasks()
			} else {
				// For cobra, we could handle other status specifically if needed, 
				// but following original logic:
				tasks, err = service.GetActiveTasks()
			}
		} else {
			tasks, err = service.GetActiveTasks()
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			if jsonOutput {
				printJSON([]model.ActiveTaskResp{})
			} else {
				fmt.Println("No tasks found")
			}
			return
		}

		if jsonOutput {
			printJSON(tasks)
		} else {
			fmt.Printf("Found %d tasks:\n\n", len(tasks))
			for _, t := range tasks {
				fmt.Printf("  [%s] %s - %s\n", t.Status, t.Title, t.Category)
				fmt.Printf("    ID: %s\n\n", t.ID)
			}
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get task details by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		task, err := service.GetTask(taskID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(task)
		} else {
			printTask(task)

			if len(task.Logs) > 0 {
				fmt.Println("\nWorklogs:")
				for _, log := range task.Logs {
					fmt.Printf("  [%s] %s\n", log.CreatedAt.Format("2006-01-02 15:04"), log.LogText)
				}
			}
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update task by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		deadlineTime := parseDeadline(deadline)

		if status != "" {
			progressReq := model.UpdateProgressReq{
				NewStatus: status,
			}
			err := service.UpdateProgress(taskID, progressReq)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if jsonOutput {
				printJSON(map[string]string{"id": taskID, "status": "updated"})
			} else {
				fmt.Printf("Task status updated: %s\n", taskID)
			}
			return
		}

		req := model.UpdateTaskReq{
			Category:    category,
			Description: desc,
			Targets:     targets,
			Links:       links,
			Deadline:    deadlineTime,
		}

		task, err := service.UpdateTask(taskID, req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(task)
		} else {
			fmt.Printf("Task updated: %s\n", taskID)
			printTask(task)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete task by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		err := service.DeleteTask(taskID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(map[string]string{"id": taskID, "status": "deleted"})
		} else {
			fmt.Printf("Task deleted: %s\n", taskID)
		}
	},
}

var logCmd = &cobra.Command{
	Use:   "log <id> <message>",
	Short: "Add worklog to task",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		logText := strings.Join(args[1:], " ")

		req := model.UpdateProgressReq{
			LogText: logText,
		}

		err := service.UpdateProgress(taskID, req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(map[string]string{"id": taskID, "status": "worklog added"})
		} else {
			fmt.Printf("Worklog added to task: %s\n", taskID)
		}
	},
}

var summaryCmd = &cobra.Command{
	Use:   "summary [date]",
	Short: "Get daily summary (default: today)",
	Run: func(cmd *cobra.Command, args []string) {
		dateStr := ""
		if len(args) > 0 {
			dateStr = args[0]
		}

		summary, err := service.GetDailySummary(dateStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(summary)
		} else {
			fmt.Printf("Daily Summary for %s\n\n", summary.Date)

			if len(summary.Activities) == 0 {
				fmt.Println("No activities today")
				return
			}

			// Group by category
			categoryMap := make(map[string][]model.DailySummaryActivity)
			for _, a := range summary.Activities {
				categoryMap[a.Category] = append(categoryMap[a.Category], a)
			}

			for cat, activities := range categoryMap {
				fmt.Printf("### %s\n", cat)
				for _, a := range activities {
					statusIcon := map[string]string{"todo": "📝", "in-progress": "🔄", "done": "✅"}[a.Status]
					fmt.Printf("%s %s\n", statusIcon, a.TaskTitle)
					for _, log := range a.TodayLogs {
						fmt.Printf("   - %s\n", log)
					}
				}
				fmt.Println()
			}
		}
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get task statistics",
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := service.GetStatsSummary()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			printJSON(stats)
		} else {
			fmt.Println("=== Task Statistics ===")
			fmt.Printf("Total Tasks: %d\n", stats.TotalTasks)
			fmt.Printf("Completed: %d\n", stats.CompletedTasks)
			fmt.Printf("In Progress: %d\n", stats.InProgressTasks)
			fmt.Printf("Todo: %d\n", stats.TodoTasks)
			fmt.Printf("Completion Rate: %.1f%%\n", stats.CompletionRate*100)

			fmt.Println("\n=== By Category ===")
			for cat, count := range stats.ByCategory {
				fmt.Printf("  %s: %d\n", cat, count)
			}

			fmt.Println("\n=== Weekly Stats ===")
			for _, s := range stats.WeeklyStats {
				fmt.Printf("  %s: created=%d, completed=%d\n", s.Date, s.Created, s.Completed)
			}
		}
	},
}

// Helpers
func parseDeadline(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		fmt.Printf("Warning: invalid deadline format, ignoring: %s\n", s)
		return nil
	}
	return &t
}

func printTask(task *model.Task) {
	fmt.Println("\nTask Details:")
	fmt.Printf("  ID: %s\n", task.ID)
	fmt.Printf("  Title: %s\n", task.Title)
	fmt.Printf("  Category: %s\n", task.Category)
	fmt.Printf("  Status: %s\n", task.Status)
	if task.Deadline != nil {
		fmt.Printf("  Deadline: %s\n", task.Deadline.Format("2006-01-02 15:04"))
	}
	if task.Description != "" {
		fmt.Printf("  Description: %s\n", task.Description)
	}
	if task.Links != "" {
		fmt.Printf("  Links:\n")
		for _, link := range strings.Split(task.Links, "\n") {
			fmt.Printf("    - %s\n", link)
		}
	}
}

func init() {
	// Add subcommands to rootCmd
	rootCmd.AddCommand(createCmd, listCmd, getCmd, updateCmd, deleteCmd, logCmd, summaryCmd, statsCmd)

	// Local flags for create and update
	createCmd.Flags().StringVarP(&category, "category", "c", "", "Task category")
	createCmd.Flags().StringVarP(&desc, "desc", "d", "", "Task description")
	createCmd.Flags().StringVarP(&links, "links", "l", "", "Task links (one per line)")
	createCmd.Flags().StringVarP(&targets, "target", "t", "", "Task targets")
	createCmd.Flags().StringVar(&deadline, "deadline", "", "Deadline (ISO8601 format)")

	updateCmd.Flags().StringVarP(&category, "category", "c", "", "Task category")
	updateCmd.Flags().StringVarP(&desc, "desc", "d", "", "Task description")
	updateCmd.Flags().StringVarP(&links, "links", "l", "", "Task links (one per line)")
	updateCmd.Flags().StringVarP(&targets, "target", "t", "", "Task targets")
	updateCmd.Flags().StringVar(&deadline, "deadline", "", "Deadline (ISO8601 format)")
	updateCmd.Flags().StringVar(&status, "new-status", "", "New status")
}
