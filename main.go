package main

import (
	"github.com/yuyudeqiu/chronicle/cmd"
)

// Build variables injected via ldflags (make build) or runtime/debug.ReadBuildInfo (go install)
var (
	gitCommit string
	gitDate   string
	buildTime string
)

func main() {
	cmd.SetBuildInfo(gitCommit, gitDate, buildTime)
	cmd.Execute()
}
