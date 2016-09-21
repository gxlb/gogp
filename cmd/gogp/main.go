package main

import (
	"github.com/vipally/cmdline"
	"github.com/vipally/gogp"
)

func main() {
	workPath := ""
	exit_code := 0

	cmdline.Summary("Tool <thiscmd> is used to generate Generic-Programming code")
	cmdline.StringVar(&workPath, "", "workPath", workPath, false, "path of gp")
	cmdline.Parse()

	if workPath == "" {
		workPath = cmdline.WorkDir()
	}

	gogp.Work(workPath)

	cmdline.Exit(exit_code)
}
