package gogp

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

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
		exp = `\Q#GOGP_DO_NOT_HAVE_ANY_KEY#\E`
	}

	//fmt.Println(exp)
	return exp
}

func (this *replaceList) doReplacing(content, _path string, reverse bool) (rep string, noRep int) {
	reg := gogpExpTodoReplace
	if reverse {
		exp := this.expString()
		reg = regexp.MustCompile(exp)
	}
	//	if gDebug {
	//		if this.sectionName == "tree_sort_slice" {
	//			for i, v := range this.list {
	//				fmt.Printf("%d %#v\n", i, v)
	//			}
	//		}
	//	}
	rep = reg.ReplaceAllStringFunc(content, func(src string) (r string) {
		p, w, s := "", src, ""
		rawName := false
		if !reverse {
			elem := reg.FindAllStringSubmatch(src, 1)[0]
			p, w, s = elem[1], elem[2], elem[3]
			//.<VALUE_TYPE> <VALUE_TYPE>:
			rawName = (w == "<VALUE_TYPE>" || w == "<KEY_TYPE>") && (p == "." || s == ":")
			//fmt.Printf("[%s][%s][%s][%s][%v]\n", src, p, w, s, rawName)
		}
		if v, ok := this.getMatch(w); ok {
			if reverse {
				r = v
			} else { //gp replacing
				wv := v
				if rawName {
					wv = getRawName(v)
				}
				r = p + wv + s
				//fmt.Printf("[%s][%s]->[%s]\n", w, v, r)
			}
		} else {
			fmt.Printf("[gogp error]: [%s] has no replacing.[%s] [%s : %s]\n", w, relateGoPath(this.gpPath), relateGoPath(this.gpgPath), this.sectionName)
			//println(_path, src, this.gpgPath, this.gpPath, this.sectionName)
			r = src
			noRep++
		}
		return
	})

	return
}
