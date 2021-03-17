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

var res = []*syntax{
	&syntax{
		name:   "comment",
		usage:  "make an in line comment in fake .go file.",
		exp:    `(?sm:(?P<COMMENT>/{2,}[ |\t]*#GOGP_COMMENT))`,
		syntax: `// #GOGP_COMMENT`,
	},
	&syntax{
		name:  "if",
		usage: "",
		exp:   `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF[ |\t]+(?P<CONDK>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE(?:.*?$[\r|\n]?)[\r|\n]*(?P<F>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF.*?$[\r|\n]?)`,
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
	&syntax{
		name:  "switch",
		usage: "",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:#GOGP_SWITCH)(?:[ |\t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASES>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDSWITCH.*?$[\r|\n]?)`,
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
	&syntax{
		name:  "case",
		usage: "",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:(?:#GOGP_CASE[ |\t]+(?P<COND>[[:word:]<>\|!]+))|(?:#GOGP_DEFAULT))(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASE>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDCASE.*?$[\r|\n]*)`,
		syntax: `
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
`,
	},
	&syntax{
		name:   "require",
		usage:  "",
		exp:    `(?sm:\s*(?P<REQ>^[ |\t]*(?://)?#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ |\t]*?,[ |\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ |\t]*?\))).*?$[\r|\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r|\n]*)`,
		syntax: `// #GOGP_REQUIRE(<gp-path> [, <gpgSection>])`,
	},
	&syntax{
		name: "replace",
		exp:  `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ |\t]*,[ |\t]*(?P<REPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, literal replacement****
// #GOGP_REPLACE(<src>, <dst>)
`,
	},
	&syntax{
		name:  "map",
		usage: "",
		exp:   `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ |\t]*,[ |\t]*(?P<MAPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, which can affect brantch of #GOGP_IFDEF and #GOGP_SWITCH after this code****
// #GOGP_MAP(<src>, <dst>)
`,
	},
}

// syntax regexp descriptor
type syntax struct {
	name   string
	usage  string
	exp    string
	syntax string
}

func regexpCompile(res ...*syntax) *regexp.Regexp {
	var b bytes.Buffer
	var exp = `\Q#GOGP_DO_NOT_HAVE_ANY_KEY#\E`
	if len(res) > 0 {
		for _, v := range res {
			b.WriteString(v.exp)
			b.WriteByte('|')
		}
		b.Truncate(b.Len() - 1) //remove last '|'
		exp = b.String()

	}
	return regexp.MustCompile(exp)
}

func (st *syntax) Regexp() *regexp.Regexp {
	return regexp.MustCompile(st.exp)
}

func findRE(name string) *syntax {
	for _, v := range res {
		if v.name == name {
			return v
		}
	}
	panic(fmt.Errorf("findRE(%s) not found", name))
	return nil
}
