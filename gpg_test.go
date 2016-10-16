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
////////////////////////////////
	
	
	
     //#if a == 5 //if a

true1

//#else //else

 false1

//#endif    //end

////////////////////////////////


//#if a2  5 //if a2

	true2

//#endif //end

   //#if a3

	true3

//#endif //end


////////////////////////////////
	    
`
	sCheck := `
////////////////////////////////
a 5

true1



 false1


////////////////////////////////
a2 5

	true2



a3 

	true3


////////////////////////////////
	    
`
	//	ss := gGogpChoiceExp.FindAllStringSubmatch(s, -1)
	//	for _, v := range ss {
	//		fmt.Printf("[%s]\n", v[0])
	//		for j, vv := range v[1:] {
	//			fmt.Printf("%d  [%s]\n", j, vv)
	//		}
	//	}

	tt := gGogpChoiceExp.ReplaceAllString(s, "\n$CONDK $CONDV$T$F\n")
	if tt != sCheck {
		t.Error(tt)
	}
	//fmt.Printf("[%s]", tt)
}
