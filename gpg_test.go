package gogp

import (
	//"fmt"
	//"os"
	//"regexp"
	"testing"

	//"github.com/vipally/gogp"
)

func init() {
	//ReverseWork("github.com/vipally/gogp/examples/reverse.gpg")
}

//run gogp tool to auto-generate go file(s) in test process
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
