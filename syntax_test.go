package gogp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const (
	testPrintResult = false
)

func TestAllRegexpSyntax(t *testing.T) {
	expAll := compileMultiRegexps(res...)
	groups := expAll.SubexpNames()
	if testPrintResult {
		fmt.Printf("%#v\n", groups)
	}
	if matched := expAll.MatchString(tstExpSyntaxAll); !matched {
		t.Errorf("test case not match")
		return
	}

	var submatchesExpect = [][]string{
		[]string{"match1", "COMMENT:// #GOGP_COMMENT"},
		[]string{"match2", "COMMENT://#GOGP_COMMENT"},
		[]string{"match3", "IFCOND,0,2:<key>", "IFCOND,1,1:!", "IFCOND,1,2:<key>", "IFCOND,2,2:<key>", "IFCOND,2,3:==", "IFCOND,2,4:xxx", "IFCOND,3,2:<key>", "IFCOND,3,3:!=", "IFCOND,3,4:xxx", "IFT:\t{if-true content}\n", "IFF:\t{if-else content}\n"},
		[]string{"match4", "IFCOND,0,2:<key>", "IFCOND,1,1:!", "IFCOND,1,2:<key>", "IFCOND,2,2:<key>", "IFCOND,2,3:==", "IFCOND,2,4:xxx", "IFCOND,3,2:<key>", "IFCOND,3,3:!=", "IFCOND,3,4:xxx", "IFT:\t{if-true content2}\n"},
		[]string{"match5", "IFCOND,0,2:<key>", "IFCOND,1,1:!", "IFCOND,1,2:<key>", "IFCOND,2,2:<key>", "IFCOND,2,3:==", "IFCOND,2,4:xxx", "IFCOND,3,2:<key>", "IFCOND,3,3:!=", "IFCOND,3,4:xxx", "IFT:\t{if-true content}\n", "IFF:\t{if-else content}\n"},
		[]string{"match6", "IFCOND,0,2:<key>", "IFCOND,1,1:!", "IFCOND,1,2:<key>", "IFCOND,2,2:<key>", "IFCOND,2,3:==", "IFCOND,2,4:xxx", "IFCOND,3,2:<key>", "IFCOND,3,3:!=", "IFCOND,3,4:xxx", "IFT:\t{if-true content2}\n"},
		[]string{"match7", "IFCOND,0,2:x", "IFT://     #GOGP_IFDEF2 yyy\n\t{if-true content}\n//     #GOGP_ELSE2\n\t{if-else content}\n//     #GOGP_ENDIF2\n", "IFF://     #GOGP_IFDEF2 yyy\n\t{if-true content}\n//     #GOGP_ELSE2\n\t{if-else content}\n//     #GOGP_ENDIF2\n"},
		[]string{"match8", "IFCOND2:xx", "IFT2://     #GOGP_IFDEF yyy\n\t      {if-true content}\n//     #GOGP_ELSE //\n\t      {if-else content}\n//     #GOGP_ENDIF //\n", "IFF2://     #GOGP_IFDEF yyy\n\t      {if-true content}\n//     #GOGP_ELSE\n\t      {if-else content}\n//     #GOGP_ENDIF //\n"},
		[]string{"match9", "SWITCHKEY:<SwitchKey>", "SWITCHCONTENT,1:<SwitchKeyValue1>", "SWITCHCONTENT,2:        {case content1}\n", "SWITCHCONTENT,1:<SwitchKeyValue2>", "SWITCHCONTENT,2:        {case content2}\n", "SWITCHCONTENT,2:        {default content1}\n"},
		[]string{"match10", "SWITCHCONTENT,1:<key>", "SWITCHCONTENT,2:        {case content3}\n", "SWITCHCONTENT,1:<key>", "SWITCHCONTENT,2:        {case content4}\n", "SWITCHCONTENT,2:        {default content2}\n"},
		[]string{"match11", "SWITCHKEY:<SwitchKey>", "SWITCHCONTENT,1:<SwitchKeyValue1>", "SWITCHCONTENT,2:        {case content1}\n", "SWITCHCONTENT,1:<SwitchKeyValue2>", "SWITCHCONTENT,2:        {case content2}\n", "SWITCHCONTENT,2:        {default content1}\n"},
		[]string{"match12", "SWITCHCONTENT,1:<key>", "SWITCHCONTENT,2:        {case content3}\n", "SWITCHCONTENT,1:<key>", "SWITCHCONTENT,2:        {case content4}\n", "SWITCHCONTENT,2:        {default content2}\n"},
		[]string{"match13", "CASEKEY,0,2:<key>", "CASECONTENT:        {case content3}\n"},
		[]string{"match14", "CASEKEY,0,2:<key>", "CASECONTENT:        {case content4}\n"},
		[]string{"match15", "CASECONTENT:        {default content2}\n"},
		[]string{"match16", "CASEKEY,0,2:<key>", "CASECONTENT:        {case content3}\n"},
		[]string{"match17", "CASEKEY,0,2:<key>", "CASECONTENT:        {case content4}\n"},
		[]string{"match18", "CASECONTENT:        {default content2}\n"},
		[]string{"match19", "REQ:// #GOGP_REQUIRE(<gp-path> , gpgSection)", "REQP:<gp-path>", "REQN:gpgSection"},
		[]string{"match20", "GPGCFG:<config-name>"},
		[]string{"match21", "REPSRC:<src>", "REPDST:<dst>"},
		[]string{"match22", "MAPSRC:<src>", "MAPDST:<dst>"},
		[]string{"match23", "IGNORE: \n     {ignore content} \n// "},
		[]string{"match24", "GPONLY: \n     {gp-only content} \n// "},
		[]string{"match25", "FILEB:// #GOGP_FILE_BEGIN"},
		[]string{"match26", "FILEE:// #GOGP_FILE_END"},
		[]string{"match27", "ONCE: \n    {only generate once from a gp file} "},
		[]string{"match28", "EMPTY_LINE:\n\n\n\n\n\n\n"},
	}
	var submatches [][]string
	rep := expAll.ReplaceAllStringFunc(tstExpSyntaxAll, func(src string) string {
		if testPrintResult {
			fmt.Println("-----------------------------")
		}

		elem := expAll.FindAllStringSubmatch(src, -1)[0]
		subs := make([]string, 0, 5)
		subs = append(subs, fmt.Sprintf("match%d", len(submatches)+1))
		for i, v := range groups {
			if elem[i] != "" && i > 0 {
				if testPrintResult {
					fmt.Printf("%d %s-------\n%s\n", i, v, elem[i])
				}
				switch {
				case v == "CASEKEY" || v == "IFCOND":
					ss := strings.Split(elem[i], "||")
					for j, vv := range ss {
						c := gogpExpCondition.FindAllStringSubmatch(vv, -1)[0]
						for k, v := range c {
							if testPrintResult {
								fmt.Println(j, k, v)
							}
							if v != "" && i > 0 && k > 0 {
								subs = append(subs, fmt.Sprintf("%s,%d,%d:%s", groups[i], j, k, v))
							}
						}
					}
				case v == "SWITCHCONTENT":
					rep := gogpExpCases.ReplaceAllStringFunc(elem[i], func(src string) string {
						cases := gogpExpCases.FindAllStringSubmatch(src, -1)[0]
						if testPrintResult {
							fmt.Printf("case: %#v\n", cases)
						}
						for j, v := range cases {
							if v != "" && i > 0 && j > 0 {
								subs = append(subs, fmt.Sprintf("%s,%d:%s", groups[i], j, v))
							}
						}
						return ""
					})
					if testPrintResult {
						fmt.Println(rep)
					}
				default:
					if i > 0 {
						subs = append(subs, fmt.Sprintf("%s:%s", groups[i], elem[i]))
					}
				}
			}
		}

		//fmt.Println(subs)
		submatches = append(submatches, subs)

		return ""
	})
	if testPrintResult {
		fmt.Println("replaced:----------------------------------------", rep)
		fmt.Println(testShowStringList(submatches))
	}
	if err := testCheckStrings(submatches, submatchesExpect); err != nil {
		if !testPrintResult {
			fmt.Println(testShowStringList(submatches))
		}
		t.Error(err)
	}
}

func TestMultiRegexp(t *testing.T) {
	var subNamesExpect = [][]string{
		[]string{"", "IGNORE", "GPONLY", "MAPSRC", "MAPDST", "SWITCHKEY", "SWITCHCONTENT", "IFCOND", "IFT", "IFF", "IFCOND2", "IFT2", "IFF2"},
		[]string{"", "IGNORE", "REQ", "REQP", "REQN", "REQGPG", "REQCONTENT", "GPGCFG", "ONCE", "REPSRC", "REPDST", "COMMENT"},
		[]string{"", "REQ", "REQP", "REQN", "REQGPG", "REQCONTENT", "FILEB", "OPEN", "FILEE"},
		[]string{"", "FILEB", "OPEN", "FILEE", "IGNORE"},
	}
	var subNames = [][]string{
		gogpExpCodeSelector.SubexpNames(),
		gogpExpPretreatAll.SubexpNames(),
		gogpExpRequireAll.SubexpNames(),
		gogpExpReverseIgnoreAll.SubexpNames(),
	}
	if testPrintResult {
		for i, v := range subNames {
			fmt.Printf("expr-%d: %#v\n", i+1, v)
		}
	}
	if err := testCheckStrings(subNames, subNamesExpect); err != nil {
		t.Error(err)
	}
}

func testCheckStrings(got, expect [][]string) error {
	if len(got) != len(expect) {
		return fmt.Errorf("regexp num mismatch, expect %d, got %d", len(expect), len(got))
	}
	for i, v := range got {
		w := expect[i]
		if len(v) != len(w) {
			return fmt.Errorf("regexp %d submatch-num mismatch, expect %d, got %d", i+1, len(w), len(v))
		}

		for j, vv := range v {
			if vv != w[j] {
				return fmt.Errorf(`regexp %d,%d submatch mismatch, expect "%s", got "%s"`, i+1, j+1, w[j], vv)
			}
		}
	}
	return nil
}

func testShowStringList(ss [][]string) string {
	var b bytes.Buffer
	b.WriteString("[][]string{\n")
	for _, v := range ss {
		b.WriteString(fmt.Sprintf("  %#v,\n", v))
	}
	b.WriteString("}")
	return b.String()
}

const tstExpSyntaxAll = `
head
// #GOGP_COMMENT {comment code}
//#GOGP_COMMENT {comment code2}

// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{if-true content}
// #GOGP_ELSE
	{if-else content}
// #GOGP_ENDIF

// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{if-true content2}
// #GOGP_ENDIF

//#GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{if-true content}
//#GOGP_ELSE
	{if-else content}
//#GOGP_ENDIF

//#GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{if-true content2}
//#GOGP_ENDIF

// #GOGP_IFDEF x
//     #GOGP_IFDEF2 yyy
	{if-true content}
//     #GOGP_ELSE2
	{if-else content}
//     #GOGP_ENDIF2
// #GOGP_ELSE
//     #GOGP_IFDEF2 yyy
	{if-true content}
//     #GOGP_ELSE2
	{if-else content}
//     #GOGP_ENDIF2
// #GOGP_ENDIF

// #GOGP_IFDEF2 xx
//     #GOGP_IFDEF yyy
	      {if-true content}
//     #GOGP_ELSE //
	      {if-else content}
//     #GOGP_ENDIF //
// #GOGP_ELSE2
//     #GOGP_IFDEF yyy
	      {if-true content}
//     #GOGP_ELSE
	      {if-else content}
//     #GOGP_ENDIF //
// #GOGP_ENDIF2

// #GOGP_SWITCH <SwitchKey>
//    #GOGP_CASE <SwitchKeyValue1>
        {case content1}
//    #GOGP_ENDCASE
//    #GOGP_CASE <SwitchKeyValue2>
        {case content2}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content1}
//    #GOGP_ENDCASE
// #GOGP_ENDSWITCH

// #GOGP_SWITCH
//    #GOGP_CASE <key>
        {case content3}
//    #GOGP_ENDCASE
//    #GOGP_CASE <key> != val
        {case content4}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content2}
//    #GOGP_ENDCASE
// #GOGP_ENDSWITCH

//#GOGP_SWITCH <SwitchKey>
//    #GOGP_CASE <SwitchKeyValue1>
        {case content1}
//    #GOGP_ENDCASE
//    #GOGP_CASE <SwitchKeyValue2>
        {case content2}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content1}
//    #GOGP_ENDCASE
//#GOGP_ENDSWITCH

//#GOGP_SWITCH 
//    #GOGP_CASE <key>
        {case content3}
//    #GOGP_ENDCASE
//    #GOGP_CASE <key> != val
        {case content4}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content2}
//    #GOGP_ENDCASE
//#GOGP_ENDSWITCH

//    #GOGP_CASE <key>
        {case content3}
//    #GOGP_ENDCASE
//    #GOGP_CASE <key> != val
        {case content4}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content2}
//    #GOGP_ENDCASE

//#GOGP_CASE <key>
        {case content3}
//#GOGP_ENDCASE
//#GOGP_CASE <key> != val
        {case content4}
//#GOGP_ENDCASE
//#GOGP_DEFAULT
        {default content2}
//#GOGP_ENDCASE

// #GOGP_REQUIRE(<gp-path> , gpgSection)

#GOGP_GPGCFG(<config-name>)

// #GOGP_REPLACE(<src>, <dst>)

// #GOGP_MAP(<src>, <dst>)

// #GOGP_IGNORE_BEGIN 
     {ignore content} 
// #GOGP_IGNORE_END

// #GOGP_GPONLY_BEGIN 
     {gp-only content} 
// #GOGP_GPONLY_END

// #GOGP_FILE_BEGIN

// #GOGP_FILE_END

// #GOGP_ONCE 
    {only generate once from a gp file} 
// #GOGP_END_ONCE 

---






tail
`
