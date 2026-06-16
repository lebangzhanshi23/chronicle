package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

var jsonOutput bool
var dataDir string

var rootCmd = &cobra.Command{
	Use:   "chronicle",
	Short: "Chronicle is a task management tool",
	Long:  `Chronicle is a task management tool with a CLI and a web interface.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 初始化配置（命令行参数优先级最高）
		if dataDir != "" {
			config.DataDir = dataDir
		}
		// 初始化数据库
		service.InitDB(config.GetDBPath())
	},
}

// JSON output helper
func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "o", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Data directory (default: data/, or use CHRONICLE_DATA_DIR env var)")
}
