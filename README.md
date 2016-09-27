# gopg [![GoDoc](https://godoc.org/github.com/vipally/gogp?status.svg)](https://godoc.org/github.com/vipally/gogp) ![Version](https://img.shields.io/badge/version-2.1.0-green.svg)
----
	
	gopg is a tool of generic-programming for golang or any other language

----

	CopyRight @Ally Dale 2016
    Author  : Ally Dale(vipally@gmail.com)
    Blog    : http://blog.csdn.net/vipally
    Site    : https://github.com/vipally


	it is available for generate gp code for golang
	see details in example directory

## usage of gogp tool:
1. use cmdline(cmd/gogp):

        use "gopg" to generate go code file(s) from .gp file via .gpg file in current directory
        use "gopg <path>" to generate go code file(s) from .gp file via .gpg file in path directory
        use "gopg -h" for help
	
2. package usage:

		2.1 import gogp package in test file, the tool will auto work at current path
	 		import (
	 			"testing"
	 			_ "github.com/vipally/gogp" //auto run gogp tool at current path in test process
	 		)
	
		2.2 import gogp package in normal package
			import (
				"github.com/vipally/gogp"
			)
			func someFunc(){
				gogp.Work(someDir)
			}
	
