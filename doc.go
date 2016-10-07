//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

/*
package gogp is a generic-programming solution for golang or any other languages.

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

    usage samples:
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

Detail desctription:
    Tool Site: https://github.com/vipally/gogp
    Work flow: DummyGoFile  --(GPGFile[1])-->  gp_file  --(GPGFile[2])-->  real_go_files

      In this flow, DummyGoFile and GPGFile are hand making, and the rests are products
    of gogp tool.

    1. DummyGoFiles
      Sample: https://github.com/vipally/gogp/blob/master/examples/stack.go
      This is a "normal" go file with WELL-DESIGNED structure.
      Texts that matches
            "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END ...\n"
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

More gogp details:

    1. Working flow:

    DummyGoFile  --(GPGFile[1])-->  gp_file  --(GPGFile[2])-->  real_go_files

           In this flow, DummyGoFile and GPGFile are hand making, and the rests are products
        of gogp tool.

    1.1 DummyGoFile
            Sample: https://github.com/vipally/gogp/blob/master/examples/stack.go

            This is a "normal" go file with WELL-DESIGNED structure.
            Texts that matches
                 "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END ...\n"
        case will be ingored by gogp tool when loading.
            From line 3~14, we add some help info about this DummyGoFile, and that will
        not exists in products.
            At line-6 (/*   //<----This line can be...), we setted a whole-file comment
        switch corresponding to line-89 (// * /).If add "//" to head of this line, this
        file comes to a "normal" go file, we can edit,compile,test, and of cause, use
        go-fmt tool to format this file.
            After that, remove "//" from line-6. This file becomes a big-commented file.
        And will have noting for go-doc tool and no export-symbols.Of cause, this
        does nothing to do with the final products real-go files.
            But there is one limit, we can not use "/* ...  * /" style comment in this file
        anywhere again.

            Any more, from line 18~35, we defines some dummy types and methods.For making
        this file LEGAL.What we exactly need is the unique dummy identifiers (GOGPStackElem).
        Which is similar to template parameter T in C++.

            After that, we have a go-like file, but anywhere we want to be replacing with
        has been set to a unique legal identifiers.

    1.2 GPGFile
           GPGFile is an ini-format file, that defines key-value replacing cases from
        source to the product.
           "GOGP_IGNORExxxx" style sections will ignore by gogp tool.
           "GOGP_REVERSExxxx" style sections are used as GPGFile[1](reverse) flow.
        Which is used to generate .gp file from DummyGoFile.
           Reverse process replaces value(GOGPStackElem) with <key>(<STACK_ELEM>) in .gp file.
           So .gp file is a normal-go-like file that exists some <xxx> format template
        keys, which need to be replaced with proper txt to generate real-go file.
           Other styles of gpg sections are used as the last flow: generate go code
        file from .gp file. It is a mechanical matches from keys to values.

           Moreover, "GOGP_xxx" style keys are reserved by gogp tool which will never
        be replacing with.
           "GOGP_Name" is used to specify DummyGoFileName in the first flow, and specify
        go-file-name-suffix in the second flow.
           "GOGP_GpFilePath" is used to specify .gp file path in the second flow.
*/
package gogp
