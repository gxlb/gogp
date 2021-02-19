//this is a cmdline interface of gogp tool.
package main

import (
	"github.com/gxlb/gogp"
	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
)

func main() {
	var (
		filePath = ""
		codeExt  = ""
		//reverseWork = false
		forceUpdate        = false
		moreInfo           = false
		removeProductsOnly = false
		debug              = false
		exit_code          = 0
	)

	cmdline.Version(gogp.Version())
	cmdline.CopyRight(cpright.CopyRight())

	cmdline.Summary("Tool <thiscmd> is a generic-programming solution for golang or any other languages.")
	cmdline.Details(`Tool Site: https://github.com/vipally/gogp
	Work flow: DummyGoFile  --(GPGFile[1])-->  gp_file  --(GPGFile[2])-->  real_go_files
	
	  In this flow, DummyGoFile and GPGFile are hand making, and the rests are products 
	of gogp tool.
	
	1. DummyGoFiles
	  Sample: https://github.com/vipally/gogp/blob/master/examples/stack.go
	  This is a "normal" go file with WELL-DESIGNED structure.
	  Texts that matches 
	        "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END ...\n"
	case will be ingored by gogp tool when loading.
	  Any identifier who wants to be replaced with is defines as unique dummy 
	word(eg: GOGPStackElem), which is similar to template parameter T in C++.
	  GPG file "GOGP_REVERSE_xxx" style sections defines the cases to replacing 
	them to <key> format "identifiers" in GP file.
	
	2. GPG files(.gpg)
	  GPG file is an ini-format file, that defines key-value replacing cases from 
	source to the product.
	  "GOGP_IGNORE_xxx" style sections will be ignored by gogp tool.
	  "GOGP_REVERSE_xxx" style sections are defined for reverse-mode to generate 
	GP file from DummyGoFiles.
	  So normal work mode will not generate go code file for these sections.
	  "GOGP_xxx" style keys are reserved by gogp tool which will never be replacing with.
	  Corresponding GP file may with the same path and name.
	  But we can redirect it by key "GOGP_GpFilePath".
	  Key "GOGP_Name" is used to specify gp file name in reverse flow.
	  And specify go-file-name-suffix in normal flow.

	3. GP files(.gp)
	  A go-like file, but exists some <xxx> style keys,
	  that need to be replaced with which defined in GPG file.

	4. GO files(.go)
	  gogp tool auto-generated GO files are exactly normal go code files.
	  But never modify it manually, you can see this warning infomation at each file head.
	  Auto work on GoPath is recmmended.
	  gogp tool will deep travel the path to find all gpg files for processing.
	  If the generated go code file's body has no changes, this file will not be updated.
	  So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
	  So any manually modification will be restored by this tool.
	  Take care of that.
	
	usage samples:
	  gogp
	  gogp gopath`)

	cmdline.StringVar(&filePath, "", "filePath", filePath, false, "Path that gogp will work. GoPath and WorkPath is allowed.")
	//	cmdline.BoolVar(&reverseWork, "r", "reverse", reverseWork, false,
	//		`Reverse work, this mode is used to gen .gp file from a real-go file.
	//		If set this flag, the filePath flag must be a .gpg file path related to GoPath.`)
	cmdline.BoolVar(&forceUpdate, "f", "force", forceUpdate, false, "Force update all products.")
	cmdline.StringVar(&codeExt, "e", "Ext", codeExt, false, "Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.")
	cmdline.BoolVar(&moreInfo, "m", "more", moreInfo, false, "More information in working process.")
	cmdline.BoolVar(&debug, "d", "debug", debug, false, "Debug mode.")
	cmdline.BoolVar(&removeProductsOnly, "remove", "remove", removeProductsOnly, false, "Only remove all products.")

	// cmdline.AnotherName("ext", "e")
	// cmdline.AnotherName("force", "f")
	// cmdline.AnotherName("more", "m")
	// cmdline.AnotherName("debug", "d")
	cmdline.Parse()

	gogp.RemoveProductsOnly(removeProductsOnly)
	gogp.Silence(!moreInfo)
	gogp.ForceUpdate(forceUpdate)
	gogp.CodeExtName(codeExt)
	gogp.Debug(debug)
	gogp.Work(filePath)

	cmdline.Exit(exit_code)
}
