package main

import (
	"strings"

	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp"
)

var (
	copyRightCode = "//    " + strings.Replace(cpright.CopyRight(), "\n", "\n//", strings.Count(cpright.CopyRight(), "\n")-1)
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

	//fmt.Println(copyRightCode)
	//fmt.Println(cmdline.GetUsage())

	cmdline.Exit(exit_code)
}
