//this is a cmdline interface of gogp tool.
package main

import (
	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp"
)

func main() {
	var (
		filePath = ""
		codeExt  = ""
		//reverseWork = false
		forceUpdate = false
		moreInfo    = false
		exit_code   = 0
	)

	cmdline.Version(gogp.Version())
	cmdline.CopyRight(cpright.CopyRight())

	cmdline.Summary("Tool <thiscmd> is a generic-programming solution for golang or any other languages.")
	cmdline.Details(`1. .gpg files
	  An ini file in fact.It's used to define generic parameters's replacing relation.
	  "GOGP_IGNORE_xxx" style sections will be ignored by gogp tool.
	  "GOGP_REVERSE_xxx" style sections are defined for reverse-mode to generate .gp file from .go file.
	  So normal work mode will not generate go code file for these sections.
	  "GOGP_xxx" style keys are reserved by gogp tool which will never be replacing with.
	  Corresponding .gp file may with the same path and name.
	  But we can redirect it by key "GOGP_GpFilePath".
	  Key "GOGP_Name" is used to specify gp file name in reverse flow.
	  And specify go-file-name-suffix in normal flow.

	2. .gp files
	  A go-like file, but exists some <xxx> style keys,
	  that need to be replaced with which defined in .gpg file.

	3. .go files
	  gogp tool auto-generated .go files are exactly normal go code files.
	  But never modify it manually, you can see this warning at the first line in every file.
	  Auto work on GoPath is recmmended.
	  gogp tool will deep travel the path to find all gpg files for processing.
	  If the generated go code file's body has no changes, this file will not be updated.
	  So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
	  So any manually modification will be restored by this tool.
	  Take care of that.
	
	usage example:
	  gogp
	  gogp gopath`)

	cmdline.StringVar(&filePath, "", "filePath", filePath, false, "Path that gogp will work. GoPath and WorkPath is allowed.")
	//	cmdline.BoolVar(&reverseWork, "r", "reverse", reverseWork, false,
	//		`Reverse work, this mode is used to gen .gp file from a real-go file.
	//		If set this flag, the filePath flag must be a .gpg file path related to GoPath.`)
	cmdline.BoolVar(&forceUpdate, "f", "force", forceUpdate, false, "Force update all products.")
	cmdline.StringVar(&codeExt, "e", "Ext", codeExt, false, "Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.")
	cmdline.BoolVar(&moreInfo, "m", "more", moreInfo, false, "More information in working process.")
	cmdline.AnotherName("ext", "e")
	//cmdline.AnotherName("reverse", "r")
	cmdline.AnotherName("force", "f")
	cmdline.AnotherName("more", "m")
	cmdline.Parse()

	gogp.Silence(!moreInfo)
	gogp.ForceUpdate(forceUpdate)
	gogp.CodeExtName(codeExt)
	//	if reverseWork {
	//		gogp.ReverseWork(filePath)
	//	} else {
	//		gogp.Work(filePath)
	//	}
	gogp.Work(filePath)

	cmdline.Exit(exit_code)
}
