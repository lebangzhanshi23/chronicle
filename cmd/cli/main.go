package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func main() {
	// Initialize database
	service.InitDB("data/app.db")

	// Parse subcommand
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "create":
		handleCreate(args)
	case "list":
		handleList(args)
	case "get":
		handleGet(args)
	case "update":
		handleUpdate(args)
	case "delete":
		handleDelete(args)
	case "log":
		handleLog(args)
	case "summary":
		handleSummary(args)
	case "stats":
		handleStats(args)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Chronicle CLI - Task Management Tool

Usage:
  chronicle create <title> [options]
  chronicle list [status]
  chronicle get <id>
  chronicle update <id> [options]
  chronicle delete <id>
  chronicle log <id> <message>
  chronicle summary [date]
  chronicle stats
  chronicle help

Commands:
  create           Create a new task
  list             List tasks (optional: todo/in-progress/done)
  get              Get task details by ID
  update           Update task by ID
  delete           Delete task by ID
  log              Add worklog to task
  summary          Get daily summary (default: today)
  stats            Get task statistics
  help             Show this help message

Options:
  -c, --category <category>    Task category
  -d, --desc <description>     Task description
  -l, --links <links>          Task links (one per line)
  -t, --target <targets>       Task targets
  --deadline <date>            Deadline (ISO8601 format)

Examples:
  chronicle create "完成BCS文档" -c BCS -d "撰写部署文档"
  chronicle list todo
  chronicle log <id> "完成了需求分析"
  chronicle summary 2026-02-28`)
}

func handleCreate(args []string) {
	// Simple manual parsing
	category, desc, links, targets, deadline := "", "", "", "", ""
	title := ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-c" || arg == "--category" {
			if i+1 < len(args) {
				category = args[i+1]
				i++
			}
		} else if arg == "-d" || arg == "--desc" || arg == "--description" {
			if i+1 < len(args) {
				desc = args[i+1]
				i++
			}
		} else if arg == "-l" || arg == "--links" {
			if i+1 < len(args) {
				links = args[i+1]
				i++
			}
		} else if arg == "-t" || arg == "--target" || arg == "--targets" {
			if i+1 < len(args) {
				targets = args[i+1]
				i++
			}
		} else if arg == "--deadline" {
			if i+1 < len(args) {
				deadline = args[i+1]
				i++
			}
		} else if strings.HasPrefix(arg, "-") {
			// Skip unknown flags
		} else {
			title = arg
		}
	}

	if title == "" {
		fmt.Println("Error: title is required")
		fmt.Println("Usage: chronicle create <title> [options]")
		os.Exit(1)
	}

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

	fmt.Printf("Task created: %s\n", task.ID)
	printTask(task)
}

func handleList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Parse(args)

	var tasks []model.ActiveTaskResp
	var err error

	if fs.NArg() > 0 {
		status := fs.Args()[0]
		if status == "done" {
			tasks, err = service.GetHistoryTasks()
		} else {
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
		fmt.Println("No tasks found")
		return
	}

	fmt.Printf("Found %d tasks:\n\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("  [%s] %s - %s\n", t.Status, t.Title, t.Category)
		fmt.Printf("    ID: %s\n\n", t.ID)
	}
}

func handleGet(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Println("Error: task ID is required")
		fmt.Println("Usage: chronicle get <id>")
		os.Exit(1)
	}

	taskID := fs.Args()[0]
	task, err := service.GetTask(taskID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	printTask(task)

	if len(task.Logs) > 0 {
		fmt.Println("\nWorklogs:")
		for _, log := range task.Logs {
			fmt.Printf("  [%s] %s\n", log.CreatedAt.Format("2006-01-02 15:04"), log.LogText)
		}
	}
}

func handleUpdate(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	var category, desc, links, targets, deadline, newStatus string
	fs.StringVar(&category, "c", "", "Category")
	fs.StringVar(&category, "category", "", "Category")
	fs.StringVar(&desc, "d", "", "Description")
	fs.StringVar(&desc, "desc", "", "Description")
	fs.StringVar(&links, "l", "", "Links")
	fs.StringVar(&links, "links", "", "Links")
	fs.StringVar(&targets, "t", "", "Targets")
	fs.StringVar(&targets, "target", "", "Targets")
	fs.StringVar(&deadline, "deadline", "", "Deadline")
	fs.StringVar(&newStatus, "new-status", "", "New status")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Println("Error: task ID is required")
		fmt.Println("Usage: chronicle update <id> [options]")
		os.Exit(1)
	}

	taskID := fs.Args()[0]
	deadlineTime := parseDeadline(deadline)

	if newStatus != "" {
		progressReq := model.UpdateProgressReq{
			NewStatus: newStatus,
		}
		err := service.UpdateProgress(taskID, progressReq)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task status updated: %s\n", taskID)
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

	fmt.Printf("Task updated: %s\n", taskID)
	printTask(task)
}

func handleDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Println("Error: task ID is required")
		fmt.Println("Usage: chronicle delete <id>")
		os.Exit(1)
	}

	taskID := fs.Args()[0]
	err := service.DeleteTask(taskID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task deleted: %s\n", taskID)
}

func handleLog(args []string) {
	fs := flag.NewFlagSet("log", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Println("Error: task ID and log message are required")
		fmt.Println("Usage: chronicle log <id> <message>")
		os.Exit(1)
	}

	taskID := fs.Args()[0]
	logText := strings.Join(fs.Args()[1:], " ")

	req := model.UpdateProgressReq{
		LogText: logText,
	}

	err := service.UpdateProgress(taskID, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Worklog added to task: %s\n", taskID)
}

func handleSummary(args []string) {
	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	fs.Parse(args)

	dateStr := ""
	if fs.NArg() > 0 {
		dateStr = fs.Args()[0]
	}

	summary, err := service.GetDailySummary(dateStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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

	for category, activities := range categoryMap {
		fmt.Printf("### %s\n", category)
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

func handleStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	fs.Parse(args)

	stats, err := service.GetStatsSummary()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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
