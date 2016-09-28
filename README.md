# gopg [![GoDoc](https://godoc.org/github.com/vipally/gogp?status.svg)](https://godoc.org/github.com/vipally/gogp) ![Version](https://img.shields.io/badge/version-2.9.0-green.svg)
----
	
	gopg is a tool of generic-programming for golang or any other language

----

	CopyRight 2016 @Ally Dale. All rights reserved.
    Author  : Ally Dale(vipally@gmail.com)
    Blog    : http://blog.csdn.net/vipally
    Site    : https://github.com/vipally

	It is available for generate gp code for golang
	See details in example directory

## usage of gogp tool:
1. (Recommend)use cmdline(cmd/gogp):

        Tool gogp is used to generate Generic-Programming code
		Usage:
    		gogp [-r=<reverseWork>] <filePath>]
		-r=<reverseWork>
      		Reverse work, this mod is used to gen .gp file from a real-go file.
      		If set this flag, the filePath flag must be a .gpg file path related to GoPath.
  		<filePath>  string
      		Path that gogp will work, if not set, it will work on GoPath.
	
2. package usage:

		2.1 (Recommend)import gogp package in test file
	 		import (
	 			//"testing"
				"github.com/vipally/gogp"
	 		)
			func init() {
				gogp.WorkOnGoPath() //Recommend
				gogp.ReverseWork(gpgFilePath)
				//gogp.ReverseWork("github.com/vipally/gogp/examples/reverse.gpg")
				//gogp.Work(someDir)
			}
	
		2.2 (Seldom use)import gogp package in normal package
			import (
				"github.com/vipally/gogp"
			)
			func someFunc(){
				gogp.WorkOnGoPath()
				gogp.ReverseWork(gpgFilePath)
				//gogp.Work(someDir)
			}
	
