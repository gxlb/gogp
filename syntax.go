// MIT License
//
// Copyright (c) 2021 @gxlb
// Url:
//     https://github.com/gxlb
//     https://gitee.com/gxlb
// AUTHORS:
//     Ally Dale <vipally@gamil.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gogp

import (
	"bytes"
	"fmt"

	"regexp"
	//regexp "github.com/dlclark/regexp2"
)

var res = []*syntax{
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#comment",
		usage: "make an in line comment in fake .go file.",
		exp:   `(?sm:(?P<COMMENT>/{2,}[ |\t]*#GOGP_COMMENT))`,
		syntax: `
// #GOGP_COMMENT {expected code}
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#if",
		usage: "double-way branch selector by condition",
		exp:   `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF[ |\t]+(?P<IFCOND>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<IFT>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE(?:.*?$[\r|\n]?)[\r|\n]*(?P<IFF>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF.*?$[\r|\n]?)`,
		syntax: `
// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
[// #GOGP_ELSE
	{else content}]
// #GOGP_ENDIF

// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
// #GOGP_ENDIF
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#switch",
		usage: "multi-way branch selector by condition",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:#GOGP_SWITCH)(?:[ |\t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:[ |\t]*?.*?$)[\r|\n]*(?P<SWITCHCONTENT>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDSWITCH.*?$[\r|\n]?)`,
		syntax: `
**** it is multi-switch logic(more than one case brantch can trigger out) ****
// #GOGP_SWITCH [<SwitchKey>] 
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
// #GOGP_GOGP_ENDSWITCH
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#case",
		usage: "branches of switch syntax",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:(?:#GOGP_CASE[ |\t]+(?P<CASEKEY>[[:word:]<>\|!]+))|(?:#GOGP_DEFAULT))(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASECONTENT>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDCASE.*?$[\r|\n]*)`,
		syntax: `
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#require",
		usage: "require another .gp file",
		exp:   `(?sm:(?P<REQ>(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ |\t]*?,[ |\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ |\t]*?\))).*?$[\r|\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r|\n]*)`,
		syntax: `
// #GOGP_REQUIRE(<gp-path> [, <gpgSection>])
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#replace",
		usage: "build-in key-value replace command for generating .gp file",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ |\t]*,[ |\t]*(?P<REPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, literal replacement****
// #GOGP_REPLACE(<src>, <dst>)
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#map",
		usage: "build-in key-value define for generating .gp file. Which can affect brantch of #if and #switch after this code.",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ |\t]*,[ |\t]*(?P<MAPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, which can affect brantch of #GOGP_IFDEF and #GOGP_SWITCH after this code****
// #GOGP_MAP(<src>, <dst>)
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#ignore",
		usage: "txt that will ignore by gogp tool.",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\r|\n]*)`,
		syntax: `
// #GOGP_IGNORE_BEGIN 
     {ignore-content} 
// #GOGP_IGNORE_END
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#gp-only",
		usage: "txt that will stay at .gp file only. Which will ignored at final .go file.",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\r|\n]*)`,
		syntax: `
// #GOGP_GPONLY_BEGIN 
     {gp-only content} 
// #GOGP_GPONLY_END
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#empty-line",
		usage: "empty line.",
		exp:   `(?-sm:(?P<EMPTY_LINE>[\r|\n]{3,}))`,
		syntax: `
{empty-lines} 
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#trim-empty-line",
		usage: "trim empty line",
		exp:   `(?-sm:^[\r|\n]*(?P<CONTENT>.*?)[\r|\n]*$)`,
		syntax: `
{empty-lines} 
{contents}
{empty-lines} 
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#gpg-config",
		usage: "refer .gpg config",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)?#GOGP_GPGCFG\((?P<GPGCFG>[[:word:]<\->]+)\))`,
		syntax: `
[//] #GOGP_GPGCFG(<GPGCFG>)
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#once",
		usage: "code that will generate once during one .gp file processing.",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ONCE(?:[ |\t]*?//.*?$)?[\r|\n]*(?P<ONCE>.*?)[\r|\n]?(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_END_ONCE.*?$[\r|\n]?)`,
		syntax: `
// #GOGP_ONCE 
    {only generate once from a gp file} 
// #GOGP_END_ONCE 
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#file-begin",
		usage: "file head of a fake .go file.",
		exp:   `(?sm:(?P<FILEB>(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_FILE_BEGIN(?:[ |\t]+(?P<OPEN>[[:word:]]+))?).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_begin.*?(?://)?#GOGP_IGNORE_END ///gogp_file_begin.*?$)?[\r|\n]*)`,
		syntax: `
// #GOGP_FILE_BEGIN
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#file-end",
		usage: "file tail of a fake .go file.",
		exp:   `(?sm:(?P<FILEE>(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_FILE_END).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_end.*?(?://)?#GOGP_IGNORE_END ///gogp_file_end.*?$)?[\r|\n]*)`,
		syntax: `
// #GOGP_FILE_END
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: true,
		name:         "#to-replace",
		usage:        "literal that waiting to replacing.",
		exp:          `(?P<REPPREFIX>.?)(?P<REPLACEKEY>\<[[:alpha:]_][[:word:]]*\>)(?P<REPSUFFIX>.?)`,
		syntax: `
<{to-replace}>
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: true,
		name:         "#condition",
		usage:        "txt that for #if or #case condition field parser.",
		exp:          `(?sm:^[ |\t]*(?P<NOT>!)?[ |\t]*(?P<KEY>[[:word:]<>]+)[ |\t]*(?:(?P<OP>==|!=)[ |\t]*(?P<VALUE>[[:word:]]+))?[ |\t]*)`,
		syntax: `
<key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
`,
	},
}

// syntax regexp descriptor
type syntax struct {
	name         string
	usage        string
	exp          string
	syntax       string
	ignoreInList bool
}

func compileMultiRegexps(res ...*syntax) *regexp.Regexp {
	var b bytes.Buffer
	var exp = `\Q#GOGP_DO_NOT_HAVE_ANY_REGEXP_SYNTAX#\E`
	if len(res) > 0 {
		for _, v := range res {
			if !v.ignoreInList {
				b.WriteString(v.exp)
				b.WriteByte('|')
			}
		}
		b.Truncate(b.Len() - 1) //remove last '|'
		exp = b.String()
	}
	return regexp.MustCompile(exp)
}

func (st *syntax) Regexp() *regexp.Regexp {
	return regexp.MustCompile(st.exp)
}

func findSyntax(name string) *syntax {
	for _, v := range res {
		if v.name == name {
			return v
		}
	}
	panic(fmt.Errorf("findSyntax(%s) not found", name))
	return nil
}

var (
	gogpExpTodoReplace   = findSyntax("#replace").Regexp()
	gogpExpIgnore        = findSyntax("#ignore").Regexp()
	gogpExpCases         = findSyntax("#case").Regexp()
	gogpExpEmptyLine     = findSyntax("#empty-line").Regexp()
	gogpExpTrimEmptyLine = findSyntax("#trim-empty-line").Regexp()
	gogpExpRequire       = findSyntax("#require").Regexp()
	gogpExpCondition     = findSyntax("#condition").Regexp()
	gogpExpComment       = findSyntax("#comment").Regexp()

	gogpExpCodeSelector = compileMultiRegexps(
		findSyntax("#ignore"),
		findSyntax("#gp-only"),
		findSyntax("#map"),
		findSyntax("#if"),
		findSyntax("#switch"),
	)

	gogpExpPretreatAll = compileMultiRegexps(
		findSyntax("#ignore"),
		findSyntax("#require"),
		findSyntax("#gpg-config"),
		findSyntax("#once"),
		findSyntax("#replace"),
		findSyntax("#comment"),
	)

	gogpExpRequireAll = compileMultiRegexps(
		findSyntax("#require"),
		findSyntax("#file-begin"),
		findSyntax("#file-end"),
	)

	gogpExpReverseIgnoreAll = compileMultiRegexps(
		findSyntax("#file-begin"),
		findSyntax("#file-end"),
		findSyntax("#ignore"),
	)
)
