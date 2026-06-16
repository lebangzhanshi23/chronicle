package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// BuildInfo holds version information.
// Populated either by ldflags (make build) or runtime/debug.ReadBuildInfo (go install).
var BuildInfo = struct {
	Commit string
	Date   string
	Time   string
}{}

// SetBuildInfo is called from main.go to pass ldflags-injected values.
func SetBuildInfo(commit, date, time string) {
	BuildInfo.Commit = commit
	BuildInfo.Date = date
	BuildInfo.Time = time
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	// Skip root's PersistentPreRun (no DB needed for version)
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		// Try runtime/debug.ReadBuildInfo first — works with go install
		versionStr := "(unknown)"
		if info, ok := debug.ReadBuildInfo(); ok {
			versionStr = info.Main.Version
			if versionStr == "" {
				versionStr = "(devel)"
			}
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					if BuildInfo.Commit == "" {
						// Truncate to short hash for consistency with ldflags
						if len(s.Value) > 7 {
							BuildInfo.Commit = s.Value[:7]
						} else {
							BuildInfo.Commit = s.Value
						}
					}
				case "vcs.time":
					if BuildInfo.Date == "" {
						if len(s.Value) >= 10 {
							BuildInfo.Date = s.Value[:10]
						} else {
							BuildInfo.Date = s.Value
						}
					}
				}
			}
		}

		// If ldflags provided a commit but version is (devel), show a nicer string
		if versionStr == "(devel)" && BuildInfo.Commit != "" {
			versionStr = "git-" + BuildInfo.Commit
		}

		fmt.Printf("chronicle %s\n", versionStr)

		if BuildInfo.Commit != "" {
			fmt.Printf("Commit: %s\n", BuildInfo.Commit)
		}
		if BuildInfo.Date != "" {
			fmt.Printf("Date:   %s\n", BuildInfo.Date)
		}
		if BuildInfo.Time != "" {
			fmt.Printf("Build:  %s\n", BuildInfo.Time)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
