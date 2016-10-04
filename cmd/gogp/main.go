//this is a cmdline interface of gogp tool.
package main

import (
	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp"
)

func main() {
	filePath := ""
	codeExt := ""
	reverseWork := false
	exit_code := 0

	cmdline.Version(gogp.Version())
	cmdline.CopyRight(cpright.CopyRight())

	cmdline.Summary("Tool <thiscmd> is a generic-programming solution for golang or any other languages.")
	cmdline.Details(`1. .gpg files
		An ini file in fact.It's used to define generic parameters's replacing relation.
		Corresponding .gp file may with the same path and name, but we can redirect it by key "GOGP_GpFilePath".
		Section "GOGP_REVERSE" is defined for ReverseWork mode to auto-generate .gp file from .go file.
		So normal work mode will not generate go code file for this section.
		
	2. .gp files
		A go-like file, but exists some <xxx> format keys, need to be replaced with which defined in .gpg file.
		
	3. .go files
		gogp tool auto-generated .go files can be identification and compiled as well as normal go code files.
		But never modify it manualy, you can see this warning at the first line in every file.
		Auto work on GoPath is recmmended. 
		gogp tool will deep travel the path to find all .gpg files to generate go code files for them.
		If the generated go code file's body has no changes, this file will not be updated.
		So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
		So any manualy modification will be restored by tool.
		Take care of that.`)

	cmdline.StringVar(&filePath, "", "filePath", filePath, true, "Path that gogp will work. GoPath and WorkPath is allowed.")
	cmdline.BoolVar(&reverseWork, "r", "reverseWork", reverseWork, false, "Reverse work, this mode is used to gen .gp file from a real-go file.\nIf set this flag, the filePath flag must be a .gpg file path related to GoPath.")
	cmdline.StringVar(&codeExt, "e", "codeExt", codeExt, false, "Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.")
	cmdline.Parse()

	//	if filePath == "" {
	//		filePath = cmdline.GoPath()
	//	}
	gogp.CodeExtName(codeExt)
	if reverseWork {
		gogp.ReverseWork(filePath)
	} else {
		gogp.Work(filePath)
	}

	cmdline.Exit(exit_code)
}
