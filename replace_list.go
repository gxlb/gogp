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
	"sort"
	"strings"
	//regexp "github.com/dlclark/regexp2"
)

//remove "*" from src
// func getRawName(src string) (r string) {
// 	r = strings.Replace(src, "*", "", -1)
// 	return
// }

//cases of template replacing
type replaceCase struct {
	key, value string
}

//template replacing list, sort by value decend
type replaceList struct {
	list        []*replaceCase
	match       map[string]string
	sectionName string
	gpgPath     string
	gpPath      string
}

func (this *replaceList) sort() {
	sort.Sort(this)
}

func (this *replaceList) insert(k, v string, reverse bool) int {
	return this.push(&replaceCase{key: k, value: v}, reverse)
}

func (this *replaceList) push(v *replaceCase, reverse bool) int {
	if strings.HasPrefix(v.key, keyReservePrefix) { //do not match reserved keys
		return 0
	}
	if v.value != "" {
		this.list = append(this.list, v)
	}
	if reverse {
		this.match[v.value] = v.key //match key from value
	} else {
		this.match[v.key] = v.value //match value from key
	}

	return this.Len()
}

func (this *replaceList) clear() {
	this.list = nil
	this.match = make(map[string]string)
	return
}

func (this *replaceList) getMatch(key string) (match string, ok bool) {
	match, ok = this.match[key]
	return
}

func (this *replaceList) Len() int {
	return len(this.list)
}

//sort by value descend
//so in regexp, with the same prefix, the longer will match first
//eg: hello|hehe|he, "he" has the lowest priority and "hello" has the highest
func (this *replaceList) Less(i, j int) bool {
	l, r := this.list[i], this.list[j]
	return l.value > r.value
}

func (this *replaceList) Swap(i, j int) {
	this.list[i], this.list[j] = this.list[j], this.list[i]
}

func (this *replaceList) expString() (exp string) {
	var b bytes.Buffer
	if this.Len() > 0 {
		for _, v := range this.list {
			s := fmt.Sprintf(`\Q%s\E|`, v.value) //match raw letter
			b.WriteString(s)
		}
		b.Truncate(b.Len() - 1) //remove last '|'
		exp = b.String()
	} else {
		//avoid return "", which will match every byte
		exp = `\Q#GOGP_DO_NOT_HAVE_ANY_REPLACE_KEY#\E`
	}

	return exp
}

func (this *replaceList) doReplacing(content, _path string, reverse bool) (rep string, noRep int) {
	reg := gogpExpTodoReplace
	if reverse {
		exp := this.expString()
		reg = regexp.MustCompile(exp)
	}

	rep = reg.ReplaceAllStringFunc(content, func(src string) (r string) {
		w := src
		if !reverse {
			elem := reg.FindAllStringSubmatch(src, 1)[0]
			w = elem[1]
		}
		if v, ok := this.getMatch(w); ok {
			if reverse {
				r = v
			} else { //gp replacing
				wv := v
				r = wv
			}
		} else {
			fmt.Printf("[gogp error]: [%s] has no replacing.[%s] [%s : %s]\n", w, relateGoPath(this.gpPath), relateGoPath(this.gpgPath), this.sectionName)

			r = src
			noRep++
		}
		return
	})

	return
}
