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
          gogp tool will deep travel the path to find all gpg files for processing.
          If the generated go code file's body has no changes, this file will not be updated.
          So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
          So any manually modification will be restored by this tool.
          Take care of that.

----

CopyRight 2016 @Ally Dale. All rights reserved.
	
Author  : [Ally Dale(vipally@gmail.com)](mailto://vipally@gmail.com)

Blog    : [http://blog.csdn.net/vipally](http://blog.csdn.net/vipally)

Site    : [https://github.com/vipally](https://github.com/vipally)

----

## usage of gogp tool:
    1. (Recommend)use cmdline(cmd/gogp):
  
        Tool gogp is a generic-programming solution for golang or any other languages.
        Usage:
            gogp [-e|ext=<Ext>] [-f|force=<force>] [-m|more=<more>] [<filePath>]
          -e|ext=<Ext>  string
            Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.
          -f|force=<force>
            Force update all products.
          -m|more=<more>
            More information in working process.
          <filePath>  string
            Path that gogp will work. GoPath and WorkPath is allowed.
  
        usage eg:
           gogp
           gogp gopath
  
    2. package usage:
  
        2.1 (Recommend)import gogp/auto package in test file
          import (
              //"testing"
              _ "github.com/vipally/gogp/auto" //auto runs gogp tool on GoPath when init()
          )
    
        2.2 (Seldom use)import gogp package in test file
          import (
              //"testing"
              "github.com/vipally/gogp"
          )
          func init() {
              gogp.WorkOnGoPath() //Recommend
          }

----		
## More gogp details:

### 1. Working flow:
	
DummyGoFile  --<u>(GPGFile[1])</u>-->  gp_file  --<u>(GPGFile[2])</u>-->  real_go_files
	
	   In the flow, DummyGoFile and GPGFile are hand making, and the left are 
	products of gogp tool.
	
#### 1.1 DummyGoFile
	    Sample: https://github.com/vipally/gogp/blob/master/examples/stack.go
		
	    This is a WELL-DESIGNED structure of "normal" go file.
	    Text that matches 
	         "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END ...\n"
	will be ingored by gogp tool when loading.
	    From line 3~14, we add some help info about this DummyGoFile, and that will
	not exists in products.
	    At line-6 (/*   //<----This line can be...), we setted a whole-file comment
	switch corresponding to line-89 (//*/).If add "//" to head of this line, this
	file comes to a "normal" go file, we can edit,compile,test, and of cause, use
	go-fmt tool to format this file.
	    After that, remove "//" from line-6. This file becomes a big-commented file.
	And will have noting for go-doc tool and no export-symbols.Of cause, this
	does nothing to do with the final products real-go files.
	    But there is one limit, we can not use "/* ... */" style comment in this file
	anywhere again.
	
	    Any more, from line 18~35, we defines some dummy types and methods.For making
	this file LEGAL.What we exactly need is the unique dummy identifiers (GOGPStackElem). 
	Which is similar to template parameter T in C++.
	
	    After that, we have a go-like file, but anywhere we want to be replacing with 
	has been set to a unique legal identifiers.
	
#### 1.2 GPGFile
	   GPGFile is an ini-format file, that defines key-value replacing cases from 
	source to the product.
	   "GOGP_IGNORExxxx" style sections will ignore by gogp tool.
	   "GOGP_REVERSExxxx" style sections are used as GPGFile[1](reverse) flow.
	Which is used to generate .gp file from DummyGoFile.
	   Reverse process replaces value(GOGPStackElem) with <key>(<STACK_ELEM>) in .gp file.
	   So .gp file is a normal-go-like file that exists some <xxx> format template 
	keys, which	need to be replaced with proper txt to generate real-go file.
	   Other styles of gpg sections are used as the last flow: generate go code 
	file from .gp file. It is a mechanical matches from keys to values.
	
	   Moreover, "GOGP_xxxx" style keys are reserved by gogp tool, and they will 
	not do replacing work.
	   "GOGP_Name" is used to specify DummyGoFileName in the first flow, and specify 
	go-file-name-suffix in the second flow.
	   "GOGP_GpFilePath" is used to specify .gp file path in the second flow.
	
	   
	
