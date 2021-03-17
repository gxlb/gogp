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

	tt := gogpExpPretreatAll.ReplaceAllString(s, "$REQP $REQN $REQGPG $CONDK $T $F $GPGCFG $ONCE $REPSRC $REPDST $MAPSRC $MAPDST\n")
	if tt != sCheck {
		t.Errorf("\n%#v\n%#v\n%#v\n", s, sCheck, tt)
		fmt.Printf("[%s\n]", tt)
		fmt.Printf("[%#v\n]", gogpExpPretreatAll.SubexpNames())
	}

	//fmt.Printf("[%s\n]", gGogpExpPretreatAll.String())

	ss := gogpExpCodeSelector.FindAllStringSubmatch(`xxx //  #GOGP_MAP(GOGP_IfIsPointerFlagValue, yes) yyy
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
	expSwitch := regexp.MustCompile(expTxtSwitch)
	fmt.Println("ssSwitch", expSwitch.MatchString(txt))
	fmt.Printf("%#v\n", expSwitch.SubexpNames())

	rep := expSwitch.ReplaceAllStringFunc(txt, func(src string) string {
		cases := expSwitch.FindAllStringSubmatch(src, -1)[0][1]
		fmt.Printf("match: %#v\n", cases)
		if true {
			exp := regexp.MustCompile(expTxtCase)
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
	exp := regexp.MustCompile(fmt.Sprintf("%s|%s", expTxtIf, expTxtIf2))
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

/*
//  // to remove code

	// // // #GOGP_COMMENT
	// expTxtGogpComment = `(?sm:(?P<COMMENT>/{2,}[ |\t]*#GOGP_COMMENT))`

	// //generic-programming flag <XXX>
	// expTxtTodoReplace = `(?P<P>.?)(?P<W>\<[[:alpha:]_][[:word:]]*\>)(?P<S>.?)`

	// // ignore all text format:
	// // //#GOGP_IGNORE_BEGIN <content> //#GOGP_IGNORE_END
	// expTxtIgnore = `(?sm:\s*//#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\r|\n]*)`
	// expTxtGPOnly = `(?sm:\s*//#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\r|\n]*)`

	// // select by condition <cd> defines in gpg file:
	// // //#GOGP_IFDEF <cd> <true_content> //#GOGP_ELSE <false_content> //#GOGP_ENDIF
	// // "<key> || ! <key> || <key> == xxx || <key> != xxx"
	// expTxtIf  = `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF[ |\t]+(?P<CONDK>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE(?:.*?$[\r|\n]?)[\r|\n]*(?P<F>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF.*?$[\r|\n]?)`
	// expTxtIf2 = `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF2[ |\t]+(?P<CONDK2>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T2>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE2(?:.*?$[\r|\n]?)[\r|\n]*(?P<F2>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF2.*?$[\r|\n]?)`

	// // " <key> || !<key> || <key> == xxx || <key> != xxx "
	// // [<NOT>] <KEY> [<OP><VALUE>]
	// expCondition = `(?sm:^[ |\t]*(?P<NOT>!)?[ |\t]*(?P<KEY>[[:word:]<>]+)[ |\t]*(?:(?P<OP>==|!=)[ |\t]*(?P<VALUE>[[:word:]]+))?[ |\t]*)`

	// //#GOGP_SWITCH [<SWITCHKEY>] <CASES> #GOGP_GOGP_ENDSWITCH
	// expTxtSwitch = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:#GOGP_SWITCH)(?:[ |\t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASES>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDSWITCH.*?$[\r|\n]?)`

	// //#GOGP_CASE <COND> <CASE> #GOGP_ENDCASE
	// //#GOGP_DEFAULT <CASE> #GOGP_ENDCASE
	// expTxtCase = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:(?:#GOGP_CASE[ |\t]+(?P<COND>[[:word:]<>\|!]+))|(?:#GOGP_DEFAULT))(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASE>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDCASE.*?$[\r|\n]*)`

	// // require another gp file:
	// // //#GOGP_REQUIRE(<gpPath> [, <gpgSection>])
	// expTxtRequire   = `(?sm:\s*(?P<REQ>^[ |\t]*(?://)?#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ |\t]*?,[ |\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ |\t]*?\))).*?$[\r|\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r|\n]*)`
	// expTxtEmptyLine = `(?sm:(?P<EMPTY_LINE>[\r|\n]{3,}))`

	// //must be "-sm", otherwise it with will repeat every line
	// expTxtTrimEmptyLine = `(?-sm:^[\r|\n]*(?P<CONTENT>.*?)[\r|\n]*$)`

	// // get gpg config string:
	// // #GOGP_GPGCFG(<cfgName>)
	// expTxtGetGpgCfg = `(?sm:(?://)?#GOGP_GPGCFG\((?P<GPGCFG>[[:word:]]+)\))`

	// // #GOGP_REPLACE(<src>,<dst>)
	// expTxtReplaceKey = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ |\t]*,[ |\t]*(?P<REPDST>\S+)\))`
	// expTxtMapKey     = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ |\t]*,[ |\t]*(?P<MAPDST>\S+)\))`

	// //remove "*" from value type such as "*string -> string"
	// // #GOGP_RAWNAME(<strValueType>)
	// //gsExpTxtRawName = "(?-sm:(?://)?#GOGP_RAWNAME\((?P<RAWNAME>\S+)\))"

	// // only generate <content> once from a gp file:
	// // //#GOGP_ONCE <content> //#GOGP_END_ONCE
	// expTxtOnce = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)//#GOGP_ONCE(?:[ |\t]*?//.*?$)?[\r|\n]*(?P<ONCE>.*?)[\r|\n]*[ |\t]*?(?://)??#GOGP_END_ONCE.*?$[\r|\n]*)`

	// expTxtFileBegin = `(?sm:\s*(?P<FILEB>//#GOGP_FILE_BEGIN(?:[ |\t]+(?P<OPEN>[[:word:]]+))?).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_begin.*?(?://)?#GOGP_IGNORE_END ///gogp_file_begin.*?$)?[\r|\n]*)`
	// expTxtFileEnd   = `(?sm:\s*(?P<FILEE>//#GOGP_FILE_END).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_end.*?(?://)?#GOGP_IGNORE_END ///gogp_file_end.*?$)?[\r|\n]*)`

	// gogpExpTodoReplace      = regexp.MustCompile(expTxtTodoReplace)
	// gogpExpPretreatAll      = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtRequire, expTxtGetGpgCfg, expTxtOnce, expTxtReplaceKey, expTxtGogpComment))
	// gogpExpIgnore           = regexp.MustCompile(expTxtIgnore)
	// gogpExpCodeSelector     = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtGPOnly, expTxtIf, expTxtIf2, expTxtMapKey, expTxtSwitch))
	// gogpExpCases            = regexp.MustCompile(expTxtCase)
	// gogpExpEmptyLine        = regexp.MustCompile(expTxtEmptyLine)
	// gogpExpTrimEmptyLine    = regexp.MustCompile(expTxtTrimEmptyLine)
	// gogpExpRequire          = regexp.MustCompile(expTxtRequire)
	// gogpExpRequireAll       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtRequire, expTxtFileBegin, expTxtFileEnd))
	// gogpExpReverseIgnoreAll = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtFileBegin, expTxtFileEnd, expTxtIgnore))
	// gogpExpCondition        = regexp.MustCompile(expTxtRequire)
	// gogpExpComment          = regexp.MustCompile(expTxtGogpComment)
 */
