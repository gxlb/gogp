package main

import (
	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp"
)

func main() {
	filePath := ""
	reverseWork := false
	exit_code := 0

	cmdline.Version(gogp.Version())
	cmdline.CopyRight(cpright.CopyRight())

	cmdline.Summary("Tool <thiscmd> is used to generate Generic-Programming code")
	cmdline.StringVar(&filePath, "", "filePath", filePath, false, "Path that gogp will work, if not set, it will work on GoPath.")
	cmdline.BoolVar(&reverseWork, "r", "reverseWork", reverseWork, false, "Reverse work, this mod is used to gen .gp file from a real-go file.\n      If set this flag, the filePath flag must be a .gpg file path related to GoPath.")
	cmdline.Parse()

	//	if filePath == "" {
	//		filePath = cmdline.GoPath()
	//	}
	if reverseWork {
		gogp.ReverseWork(filePath)
	} else {
		gogp.Work(filePath)
	}

	cmdline.Exit(exit_code)
}
