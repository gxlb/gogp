package main

import (
	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp"
)

func main() {
	workPath := ""
	exit_code := 0

	cmdline.Version(gogp.Version())
	cmdline.CopyRight(cpright.CopyRight())

	cmdline.Summary("Tool <thiscmd> is used to generate Generic-Programming code")
	cmdline.StringVar(&workPath, "", "workPath", workPath, false, "path of .gpg file")
	cmdline.Parse()

	if workPath == "" {
		workPath = cmdline.WorkDir()
	}

	gogp.Work(workPath)

	cmdline.Exit(exit_code)
}
