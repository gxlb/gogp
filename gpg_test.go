//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

package gogp

import (
	"fmt"
	"regexp"
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
//#GOGP_REQUIRE(this_is_required4.xxx)
//#GOGP_REQUIRE(this_is_required5.xxx,section)
//#GOGP_REQUIRE(this_is_required6.xxx,#GOGP_GPGCFG(cfg))

//#GOGP_IFDEF online_cd online_t1 online_t2 #GOGP_ELSE online_f1 online_f2 #GOGP_ENDIF

////////////////////////////////
	
     //#GOGP_IFDEF a //if a

	true1


//#GOGP_ELSE //else

 false1


//#GOGP_ENDIF    //end

////////////////////////////////

//#GOGP_ONCE one-line once #GOGP_END_ONCE

//#GOGP_ONCE 
	multi-line 
	once 
//#GOGP_END_ONCE

//#GOGP_GPGCFG(gpgCfg)


//#GOGP_IFDEF a2   //if a2

	true2

//#GOGP_ENDIF //end

//    #GOGP_REPLACE(repSrc1, repDst1)
//#GOGP_REPLACE(<repSrc2>, <repDst2>)

//    #GOGP_MAP(mapSrc1, mapDst1)
//#GOGP_MAP(<mapSrc2>, <mapDst2>)




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
//required3.xxx)this_is_required4.xxx       
this_is_required5.xxx section      
this_is_required6.xxx  cfg     
   online_cd  online_t1 online_t2  online_f1 online_f2  
////////////////////////////////   a 	true1  false1  
////////////////////////////////        one-line once
        
	multi-line 
	once 
      gpgCfg 
   a2 	true2   
////////////////////////////////
	    
`
	//	ss := gGogpExpPretreatAll.FindAllStringSubmatch(s, -1)
	//	for _, v := range ss {
	//		fmt.Printf("[%s]\n", v[0])
	//		for j, vv := range v[1:] {
	//			fmt.Printf("%d  [%s]\n", j, vv)
	//		}
	//	}

	tt := gGogpExpPretreatAll.ReplaceAllString(s, "$REQP $REQN $REQGPG $CONDK $T $F $GPGCFG $ONCE $REPSRC $REPDST $MAPSRC $MAPDST\n")
	if tt != sCheck {
		t.Errorf("\n%#v\n%#v\n%#v\n", s, sCheck, tt)
		fmt.Printf("[%s\n]", tt)
		fmt.Printf("[%#v\n]", gGogpExpPretreatAll.SubexpNames())
	}

	//fmt.Printf("[%s\n]", gGogpExpPretreatAll.String())

	ss := gGogpExpCodeIgnore.FindAllStringSubmatch(`xxx //  #GOGP_MAP(GOGP_IfIsPointerFlagValue, yes) yyy
// #GOGP_REPLACE(GOGP_IfIsPointerFlagValue, yes)
	`, -1)
	fmt.Println("ss", ss)
}

func TestRegCase(t *testing.T) {
	txt := `
head
// #GOGP_SWITCH
// #GOGP_CASE x  //xxx
xxx
xxx	
// #GOGP_ENDCASE //x
// #GOGP_CASE y 
yyy
yyy	
// #GOGP_ENDCASE		//y
 /// #GOGP_CASE z //aa
zzz
zzz	
 // #GOGP_ENDCASE //xxx
 // #GOGP_DEFAULT
000
000	 
 // #GOGP_ENDCASE //default
// #GOGP_ENDSWITCH
tail
`
	expSwitch := regexp.MustCompile(gsExpTxtSwitch)
	fmt.Println("ssSwitch", expSwitch.MatchString(txt))
	fmt.Printf("%#v\n", expSwitch.SubexpNames())

	rep := expSwitch.ReplaceAllStringFunc(txt, func(src string) string {
		cases := expSwitch.FindAllStringSubmatch(src, -1)[0][1]
		fmt.Printf("match: %#v\n", cases)
		if true {
			exp := regexp.MustCompile(gsExpTxtCase)
			//ss := exp.FindAllStringSubmatch(txt, -1)
			fmt.Println("ssCase", exp.MatchString(txt))
			fmt.Printf("%#v\n", exp.SubexpNames())
			rep := exp.ReplaceAllStringFunc(cases, func(src string) string {
				elem := exp.FindAllStringSubmatch(src, -1)[0]
				fmt.Println("--match------------------------------------------------")
				fmt.Printf("%s", elem[0])
				fmt.Println("--cond--------------------")
				fmt.Printf("%s\n", elem[1])
				fmt.Println("--content--------------------")
				fmt.Printf("%s", elem[2])
				return ""
			})
			//fmt.Println("rep", rep)
			rep = rep
		}
		return ""
	})
	fmt.Println("rep", rep)
	rep = rep

}

func TestRegIf(t *testing.T) {
	txt := `
head
	//#GOGP_IFDEF true /
aaat
	//#GOGP_ELSE
bbbt
  //#GOGP_ENDIF


//#GOGP_IFDEF false
aaaf
//#GOGP_ELSE
bbbf
//#GOGP_ENDIF


// #GOGP_IFDEF true
ccct
// #GOGP_ENDIF


// #GOGP_IFDEF false
cccf
// #GOGP_ENDIF

//#GOGP_IFDEF2 true //outer
//	  #GOGP_IFDEF true //inner
aaatt
//	  #GOGP_ELSE //a
bbbtt
//	  #GOGP_ENDIF //b
//#GOGP_ELSE2 //x
//	  #GOGP_IFDEF false //inner
aaatf
//    #GOGP_ELSE //c
bbbtf
//	  #GOGP_ENDIF //d
//#GOGP_ENDIF2 //outer
tail
`
	exp := regexp.MustCompile(fmt.Sprintf("%s|%s", gsExpTxtIf, gsExpTxtIf2))
	fmt.Println("ssCase", exp.MatchString(txt))
	fmt.Printf("%#v\n", exp.SubexpNames())
	rep := exp.ReplaceAllStringFunc(txt, func(src string) string {
		elem := exp.FindAllStringSubmatch(src, -1)[0]
		fmt.Println("match------------------------------------------------")
		fmt.Printf("%s", elem[0])
		switch {
		case elem[1] != "":
			fmt.Println("cond--------------------")
			fmt.Printf("%s\n", elem[1])
			fmt.Println("t--------------------")
			fmt.Printf("%s", elem[2])
			fmt.Println("f--------------------")
			fmt.Printf("%s", elem[3])
		case elem[4] != "":
			fmt.Println("cond2--------------------")
			fmt.Printf("%s\n", elem[4])
			fmt.Println("t--------------------")
			fmt.Printf("%s", elem[5])
			fmt.Println("f--------------------")
			fmt.Printf("%s", elem[6])
		}

		return ""
	})
	fmt.Println("------------------------------------------------")
	fmt.Println("rep", rep)
	rep = rep
}
