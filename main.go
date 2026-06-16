package main

import (
	"log"

	"github.com/yuyudeqiu/chronicle/cmd"
)

// Build variables injected via ldflags
var (
	gitCommit string
	gitDate   string
	buildTime string
)

func init() {
	if gitCommit != "" {
		log.Printf("Build: commit=%s date=%s time=%s", gitCommit, gitDate, buildTime)
	}
}

func main() {
	cmd.Execute()
}
