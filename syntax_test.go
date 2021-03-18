package gogp

import (
	"fmt"
	"strings"
	"testing"
)

func TestAllRegexpSyntax(t *testing.T) {
	expAll := compileMultiRegexps(res...)
	groups := expAll.SubexpNames()
	fmt.Printf("%#v\n", groups)
	fmt.Println("MatchString", expAll.MatchString(tstExpSyntaxAll))
	rep := expAll.ReplaceAllStringFunc(tstExpSyntaxAll, func(src string) string {
		fmt.Println("-----------------------------")
		elem := expAll.FindAllStringSubmatch(src, -1)[0]

		for i, v := range groups {

			if true && elem[i] != "" && i >= 0 {
				fmt.Printf("%d %s-------\n%s\n", i, v, elem[i])
				switch {
				case v == "CASEKEY" || v == "IFCOND":
					ss := strings.Split(elem[i], "||")
					for _, vv := range ss {
						c := gogpExpCondition.FindAllStringSubmatch(vv, -1)[0]
						for j, v := range c {
							fmt.Println(j, v)
						}
					}
				case v == "SWITCHCONTENT":
					rep := gogpExpCases.ReplaceAllStringFunc(elem[i], func(src string) string {
						elem := gogpExpCases.FindAllStringSubmatch(src, -1)[0]
						fmt.Printf("case: %#v\n", elem)
						return ""
					})
					rep = rep
					return ""
				}

			}
		}
		return ""
	})
	fmt.Println("replaced:----------------------------------------", rep)
}

func TestMultiRegexp(t *testing.T) {
	fmt.Printf("%#v\n", gogpExpCodeSelector.SubexpNames())
	fmt.Printf("%#v\n", gogpExpPretreatAll.SubexpNames())
	fmt.Printf("%#v\n", gogpExpRequireAll.SubexpNames())
	fmt.Printf("%#v\n", gogpExpReverseIgnoreAll.SubexpNames())
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
