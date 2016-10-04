# gogp [![GoDoc](https://godoc.org/github.com/vipally/gogp?status.svg)](https://godoc.org/github.com/vipally/gogp) ![Version](https://img.shields.io/badge/version-2.9.0-green.svg)
----
	
package gogp is a generic-programming solution for golang or any other languages.
	
	Detail desctription:

        1. .gpg files
          An ini file in fact.It's used to define generic parameters's replacing relation.
          Corresponding .gp file may with the same path and name.
          But we can redirect it by key "GOGP_GpFilePath".
          Section "GOGP_REVERSE" is defined for ReverseWork to generate .gp file from .go file.
          So normal work mode will not generate go code file for this section.

        2. .gp files
          A go-like file, but exists some <xxx> format keys,
          that need to be replaced with which defined in .gpg file.

        3. .go files
          gogp tool auto-generated .go files are exactly normal go code files.
          But never modify it manually, you can see this warning at the first line in every file.
          Auto work on GoPath is recmmended.
          gogp tool will deep travel the path to find all .gpg files for processing.
          If the generated go code file's body has no changes, this file will not be updated.
          So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
          So any manually modification will be restored by this tool.
          Take care of that.

----

	CopyRight 2016 @Ally Dale. All rights reserved.
    Author  : Ally Dale(vipally@gmail.com)
    Blog    : http://blog.csdn.net/vipally
    Site    : https://github.com/vipally

## usage of gogp tool:
    1. (Recommend)use cmdline(cmd/gogp):

        Tool gogp is a generic-programming solution for golang or any other languages.
		Usage:
    		gogp [-r=<reverseWork>] <filePath>]
		-e=<codeExt>  string
    		Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.
		-r=<reverseWork>
      		Reverse work, this mode is used to gen .gp file from a real-go file.
      		If set this flag, the filePath flag must be a .gpg file path related to GoPath.
  		<filePath>  string
      		Path that gogp will work. GoPath and WorkPath is allowed.
			
		usage eg:
            gogp gopath
            gogp .
	
    2. package usage:

		2.1 (Recommend)import gogp package in test file
	 		import (
	 			//"testing"
				"github.com/vipally/gogp"
	 		)
			func init() {
				gogp.WorkOnGoPath() //Recommend
				//gogp.WorkOnWorkPath()
				//gogp.Work(someDir)
				//gogp.ReverseWork("github.com/vipally/gogp/examples/reverse.gpg")
			}
	
		2.2 (Seldom use)import gogp package in normal package
			import (
				"github.com/vipally/gogp"
			)
			func someFunc(){
				//gogp.WorkOnGoPath()
				//gogp.ReverseWork(gpgFilePath)
				//gogp.WorkOnWorkPath()
				//gogp.Work(someDir)
			}
	
