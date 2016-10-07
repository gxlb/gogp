//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

/*
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

usage of gogp tool:
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
          gogp.WorkOnGoPath()
      }
*/
package gogp
