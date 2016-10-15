//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

package gogp

import (
	"fmt"
	"testing"
)

func init() {
	WorkOnGoPath() //run gogp tool to auto-generate go file(s) in test process
}

func TestCallGogp(t *testing.T) {
	//	r := regexp.MustCompile("hello|he")
	//	s := r.FindAllString("I think he is say hello to HE hehe", -1)
	//fmt.Printf("%#v\n", os.Environ())
	//fmt.Println(os.Getenv("GOPATH"))

	s := `package stl

//GOGP_IGNORE_BEGIN //GOGPCommentDummyGoFile
///*
//GOGP_IGNORE_END //GOGPCommentDummyGoFile

import (
	"sort"
)

//GOGP_IGNORE_BEGIN   //GOGPDummyDefine
//these defines will never exists in real go files
type GOGPTreeNodeData int

func (this GOGPTreeNodeData) Less(o GOGPTreeNodeData) bool {
	return this < o
}

//GOGP_IGNORE_END   //GOGPDummyDefine

//tree node
type GOGPTreeNamePrefixTreeNode struct {
	GOGPTreeNodeData
	Children GOGPTreeNamePrefixSortSlice
}
`
	check := `package stl

import (
	"sort"
)

//tree node
type GOGPTreeNamePrefixTreeNode struct {
	GOGPTreeNodeData
	Children GOGPTreeNamePrefixSortSlice
}
`
	tt := gGogpIgnoreExp.ReplaceAllString(s, "\n\n")
	if tt != check {
		t.Error(tt)
	}
}

func TestChoiceExp(t *testing.T) {
	s := `
//#if a == 5
true
//#else
false
//#endif
	`
	fmt.Printf("%#v\n", gGogpChoiceExp.FindAllStringSubmatch(s, -1))
	s2 := `
//#if a == 5
true
//#endif
`

	fmt.Printf("%#v\n", gGogpChoiceExp.FindAllStringSubmatch(s2, -1))
	fmt.Printf("%#v\n", gGogpChoiceExp.SubexpNames())

	tt := gGogpChoiceExp.ReplaceAllString(s, "$COND$T$F")
	//	tt := gGogpChoiceExp.ReplaceAllStringFunc(s, func(src string) string {
	//		fmt.Println(src)
	//		return "$4"
	//	})
	fmt.Println(tt)
}
