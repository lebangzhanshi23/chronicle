package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
)

var exportCmd = &cobra.Command{
	Use:   "export [output-path]",
	Short: "Export the entire database to a file",
	Long: `Export the SQLite database file to the specified path.
If no path is given, exports to chronicle-export-{timestamp}.db in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := config.GetDBPath()

		// Check source exists
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: database not found at %s\n", dbPath)
			os.Exit(1)
		}

		// Determine output path
		outputPath := ""
		if len(args) > 0 {
			outputPath = args[0]
		} else {
			timestamp := time.Now().Format("20060102_150405")
			outputPath = fmt.Sprintf("chronicle-export-%s.db", timestamp)
		}

		// Ensure parent directory exists
		outputDir := filepath.Dir(outputPath)
		if outputDir != "." {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: cannot create directory %s: %v\n", outputDir, err)
				os.Exit(1)
			}
		}

		// Copy file
		sourceFile, err := os.Open(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot open source database: %v\n", err)
			os.Exit(1)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot create output file: %v\n", err)
			os.Exit(1)
		}
		defer destFile.Close()

		written, err := io.Copy(destFile, sourceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: copy failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Database exported successfully\n")
		fmt.Printf("   Source: %s (%d bytes)\n", dbPath, written)
		fmt.Printf("   Dest:   %s\n", outputPath)
	},
}

var importCmd = &cobra.Command{
	Use:   "import <source-path>",
	Short: "Import a database file, replacing the current data",
	Long: `Import a previously exported SQLite database file, replacing all current data.
WARNING: This will overwrite the existing database. All current data will be lost.
Use --force to skip confirmation.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]
		dbPath := config.GetDBPath()

		// Check source exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: source file not found: %s\n", sourcePath)
			os.Exit(1)
		}

		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("⚠️  WARNING: This will REPLACE the current database at:\n")
			fmt.Printf("   %s\n\n", dbPath)
			fmt.Printf("   All current data will be lost!\n")
			fmt.Printf("   Use --force to skip this warning.\n\n")
			fmt.Printf("   Are you sure? (yes/no): ")

			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				fmt.Printf("   Failed to read input: %v\n", err)
				return
			}
			if response != "yes" {
				fmt.Println("Import cancelled.")
				return
			}
		}

		// Create backup of current database before import
		backupPath := dbPath + ".pre-import-backup"
		srcFile, err := os.Open(dbPath)
		if err == nil {
			defer srcFile.Close()
			bakFile, err := os.Create(backupPath)
			if err == nil {
				if _, copyErr := io.Copy(bakFile, srcFile); copyErr == nil {
					fmt.Printf("   Auto-backup saved to: %s\n", backupPath)
				}
				bakFile.Close()
			}
		}

		// Copy source to database path
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot open source file: %v\n", err)
			os.Exit(1)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot write to database path: %v\n", err)
			os.Exit(1)
		}
		defer destFile.Close()

		written, err := io.Copy(destFile, sourceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: import failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✅ Database imported successfully\n")
		fmt.Printf("   Source: %s (%d bytes)\n", sourcePath, written)
		fmt.Printf("   Dest:   %s\n", dbPath)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd, importCmd)
	importCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
