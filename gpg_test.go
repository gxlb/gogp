//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

package gogp

import (
	"fmt"
	"testing"
)

//func init() {
//	WorkOnGoPath() //run gogp tool to auto-generate go file(s) in test process
//}

//func TestCallGogp(t *testing.T) {
//	//	r := regexp.MustCompile("hello|he")
//	//	s := r.FindAllString("I think he is say hello to HE hehe", -1)
//	//fmt.Printf("%#v\n", os.Environ())
//	//fmt.Println(os.Getenv("GOPATH"))

//}

func TestRegExps(t *testing.T) {
	s := `
	package stl
	
//#GOGP_IGNORE_BEGIN  oneline_ignore1 online_ignore2 #GOGP_IGNORE_END

//#GOGP_IGNORE_BEGIN //GOGPCommentDummyGoFile
///*
  //#GOGP_IGNORE_END //GOGPCommentDummyGoFile

import (
	"sort"
)

//#GOGP_IGNORE_BEGIN   //GOGPDummyDefine
//these defines will never exists in real go files
type GOGPTreeNodeData int

func (this GOGPTreeNodeData) Less(o GOGPTreeNodeData) bool {
	return this < o
}

 	//#GOGP_IGNORE_END   //GOGPDummyDefine

//tree node
type GOGPTreeNamePrefixTreeNode struct {
	GOGPTreeNodeData
	Children GOGPTreeNamePrefixSortSlice
}

//#GOGP_REQUIRE(this_is_required.xxx)
//#GOGP_REQUIRE(this_is_required 2.xxx , integer)
//#GOGP_REQUIRE(this_is_unchtched
//required3.xxx)

//#GOGP_IFDEF online_cd online_t1 online_t2 #GOGP_ELSE online_f1 online_f2 #GOGP_ENDIF

////////////////////////////////
	
     //#GOGP_IFDEF a //if a
//#GOGP_REQUIRE(this_is_required3.xxx)
	true1


//#GOGP_ELSE //else

 false1


//#GOGP_ENDIF    //end

////////////////////////////////


//#GOGP_IFDEF a2   //if a2

	true2

//#GOGP_ENDIF //end




////////////////////////////////
	    
`
	sCheck := `
	package stl    
    
import (
	"sort"
)    
//tree node
type GOGPTreeNamePrefixTreeNode struct {
	GOGPTreeNodeData
	Children GOGPTreeNamePrefixSortSlice
}this_is_required.xxx    
this_is_required 2.xxx integer   
//#GOGP_REQUIRE(this_is_unchtched
//required3.xxx)  online_cd  online_t1 online_t2  online_f1 online_f2
////////////////////////////////  a //#GOGP_REQUIRE(this_is_required3.xxx)
	true1  false1
////////////////////////////////  a2 	true2 
////////////////////////////////
	    
`
	//	ss := gGogpExpPretreatAll.FindAllStringSubmatch(s, -1)
	//	for _, v := range ss {
	//		fmt.Printf("[%s]\n", v[0])
	//		for j, vv := range v[1:] {
	//			fmt.Printf("%d  [%s]\n", j, vv)
	//		}
	//	}

	tt := gGogpExpPretreatAll.ReplaceAllString(s, "$REQP $REQN $CONDK $T $F\n")
	if tt != sCheck {
		t.Errorf("\n%#v\n%#v\n%#v\n", s, sCheck, tt)
		fmt.Printf("[%s\n]", tt)
	}

	//fmt.Printf("[%#v\n]", gGogpExpPretreatAll.SubexpNames())
	//	fmt.Printf("[%s\n]", gGogpExpPretreatAll.String())
}
