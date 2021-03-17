package gogp

import (
	"fmt"
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
			if true && elem[i] != "" && i > 0 {
				fmt.Printf("%d %s-------\n%s\n", i, v, elem[i])
			}
		}
		return ""
	})
	fmt.Println("replaced:----------------------------------------", rep)
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
