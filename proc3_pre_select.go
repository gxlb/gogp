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
	"fmt"
	"strings"
)

// deal with #GOGP_IFDEF and #GOGP_SWITCH for .gp file
func (this *gopgProcessor) step3PretreatGpCodeSelector(gpContent string, section string) (replaced string) {
	this.maps.clear()

	replaced = gpContent
	repCnt := 1 //init first loop
	for depth := 0; repCnt > 0; depth++ {
		replaced, repCnt = this.pretreatSelector(replaced, section, depth)
	}

	return
}

var boolTrueValues = []string{"true", "t", "yes", "y", "1"}

// parse bool value from string, treat unknown strings as false
func parseBoolValue(val string) bool {
	for _, v := range boolTrueValues {
		if strings.EqualFold(val, v) {
			return true
		}
	}
	return false
}

func (this *gopgProcessor) selectPart(section, sel string, depth int) string {
	if depth <= maxRecursionDepth {
		rep, _ := this.pretreatSelector(sel, section, depth+1)
		return gogpExpComment.ReplaceAllString(rep, "")
	}

	return gogpExpComment.ReplaceAllString(sel, "")
}

// " <key> || !<key> || <key> == xxx || <key> != xxx "
func (this *gopgProcessor) checkCondition(section, condition string, predefKey string) bool {
	conds := strings.Split(condition, "||")
	selOk := false
	for _, cond := range conds {
		elems := gogpExpCondition.FindAllStringSubmatch(cond, -1)
		if len(elems) == 0 {
			fmt.Printf("[gogp error]: [%s:%s %s] condition(%s) not match patten\n", relateGoPath(this.gpgPath), relateGoPath(this.gpPath), section, cond)
			return false
		}

		// ["", "NOT", "KEY", "OP","VALUE"]
		elem := elems[0]
		not, key, op, value := elem[1], elem[2], elem[3], elem[4]

		condValCheck := (not == "")
		if s := len(key); s >= 2 && key[0] == '<' && key[s-1] == '>' { // <key> -> key
			key = key[1 : s-1]
		}

		cfg := this.getGpgCfg(section, key, false)
		condResult := false
		switch {
		case predefKey != "":
			if op != "" || value != "" {
				fmt.Printf("[gogp warn]: [%s:%s %s] condition(%s) unexpected operator [%s, %s]\n", relateGoPath(this.gpgPath), relateGoPath(this.gpPath), section, cond, op, value)
			}
			if s := len(predefKey); s >= 2 && predefKey[0] == '<' && predefKey[s-1] == '>' { // <predefKey> -> predefKey
				predefKey = predefKey[1 : s-1]
			}
			predefVal := this.getGpgCfg(section, predefKey, false)
			condResult = (predefVal == key)

		case value != "":
			switch op {
			case "==":
				condResult = (cfg == value)
			case "!=":
				condResult = (cfg != value)
			default:
				fmt.Printf("[gogp error]: [%s:%s %s] condition(%s) undefined operator [%s]\n", relateGoPath(this.gpgPath), relateGoPath(this.gpPath), section, cond, op)
				condResult = false
			}

		default:
			if op != "" {
				fmt.Printf("[gogp warn]: [%s:%s %s] condition(%s) unexpected operator [%s]\n", relateGoPath(this.gpgPath), relateGoPath(this.gpPath), section, cond, op)
			}
			condResult = parseBoolValue(cfg)
		}

		if condResult == condValCheck {
			selOk = true
			break
		}
	}
	return selOk
}

func (this *gopgProcessor) selectByCondition(section, cond, t, f string, depth int) string {
	ret := ""
	if this.checkCondition(section, cond, "") {
		ret = this.selectPart(section, t, depth)
	} else {
		ret = this.selectPart(section, f, depth)
	}
	//fmt.Printf("[xx]selectByCondition section=%s depth=%d ret=%s\n%s, %q, %q\n", section, depth, ret, cond, t, f)
	return ret
}

func (this *gopgProcessor) selectByCases(section, cases string, predefKey string) string {
	defaultContent := ""
	found := false
	repaced := gogpExpCases.ReplaceAllStringFunc(cases, func(src string) string {
		if found { //ignore the rest cases if has found
			//return "" //treat as multi switch
		}
		elem := gogpExpCases.FindAllStringSubmatch(src, -1)[0]
		cond, content := elem[1], elem[2]
		switch {
		case cond == "": //DEFAULT branch
			defaultContent = content
		case this.checkCondition(section, cond, predefKey): //CASE branch
			found = true
			return content
		}
		return ""
	})
	if !found {
		return defaultContent
	}
	return repaced
}

func (this *gopgProcessor) pretreatSelector(gpContent string, section string, depth int) (replaced string, repCnt int) {
	if depth > maxRecursionDepth { //limit recursion depth
		s := fmt.Sprintf("[gogp error]: [%s:%s %s depth=%d] replace recursion too deep\n", relateGoPath(this.gpgPath), relateGoPath(this.gpPath), section, depth)
		fmt.Errorf("%s", s)
		return gpContent, 0
	}
	replaced = gogpExpCodeSelector.ReplaceAllStringFunc(gpContent, func(src string) (rep string) {
		repCnt++
		elem := gogpExpCodeSelector.FindAllStringSubmatch(src, -1)[0] //{"", "IGNORE", "GPONLY", "CONDK", "T", "F"}
		ignore, gponly, condk, condHit, condMiss, condk2, condHit2, condMiss2, mapK, mapV, switchKey, switchCases :=
			elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7], elem[8], elem[9], elem[10], elem[11], elem[12]

		switch {
		case condk != "":
			rep = this.selectByCondition(section, condk, condHit, condMiss, depth)

		case condk2 != "":
			rep = this.selectByCondition(section, condk2, condHit2, condMiss2, depth)

		case switchCases != "":
			return this.selectByCases(section, switchCases, switchKey)

		case mapK != "":
			this.maps.insert(mapK, mapV, false)
			rep = ""
			//println("set2", mapK, mapV)

		case ignore != "" || gponly != "":
			rep = "\n\n"
		default:
			rep = ""
		}
		//fmt.Printf("[$$] rep=%s\ndepth=%d ##src=[%#v]\n ignore=[%s] gponly=[%s] condk=[%s] t=[%q] f=[%q] condk2=[%s] t2=[%q] f2=[%q] map=[%s,%s] switchCases=[%s]\n", rep, depth, src, ignore, gponly, condk, condHit, condMiss, condk2, condHit2, condMiss2, mapK, mapV, switchCases)
		return
	})
	//fmt.Println("[$$$]", replaced)
	return
}
