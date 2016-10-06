package gogp

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/vipally/gogp/ini"
)

//cases of template replacing
type replaceCase struct {
	key, value string
}

//template replacing list, sort by value decend
type replaceList struct {
	list  []*replaceCase
	match map[string]string
}

func (this *replaceList) insert(k, v string, reverse bool) int {
	return this.push(&replaceCase{key: k, value: v}, reverse)
}

func (this *replaceList) push(v *replaceCase, reverse bool) int {
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
func (this *replaceList) expString() string {
	var b bytes.Buffer
	for _, v := range this.list {
		b.WriteString(v.value)
		b.WriteByte('|')
	}
	b.Truncate(b.Len() - 1) //remove last '|'
	exp := b.String()
	//fmt.Println(exp)
	return exp
}

//object to process gpg file
type gopgProcessor struct {
	gpgPath  string      //gpg file path
	gpPath   string      //gp file path
	codePath string      //code file path
	matches  replaceList //cases that need replacing

	nNoReplaceMathNum int //number of math that has no replace string
	nCodeFile         int
	nSkipCodeFile     int
	gpgContent        *ini.IniFile
	gpContent         string
	codeContent       string
	impName           string
}

func (this *gopgProcessor) reverseWork(gpgFilePath string) (err error) {

	if !strings.HasSuffix(gpgFilePath, gGpgExt) { //must .gpg file
		err = fmt.Errorf("[gogp]error:[%s] must be %s file at reverse mode", relateGoPath(gpgFilePath), gGpgExt)
		return
	}
	gpgFullPath := formatPath(filepath.Join(gGoPath, gpgFilePath)) //make full path

	if err = this.loadGpgFile(gpgFullPath); err != nil {
		fmt.Println(err)
		return
	}
	if err = this.reverseProcess(); err != nil {
		return
	}

	return
}

func (this *gopgProcessor) reverseProcess() (err error) {
	if !gSilence {
		fmt.Printf("[gogp]ReverseWork:[%s]\n", relateGoPath(this.gpgPath))
	}

	this.impName = gSectionReversse
	pathWithName := strings.TrimSuffix(this.gpgPath, gGpgExt)
	gpFilePath := pathWithName + gGpExt
	codeFilePath := pathWithName + gCodeExt
	this.codePath = codeFilePath
	this.gpPath = gpFilePath

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	//ignore text format like "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END"
	this.codeContent = gGogpIgnoreExp.ReplaceAllString(this.codeContent, "\n\n")

	if keys := this.gpgContent.Keys(gSectionReversse); keys != nil {
		this.matches.clear()
		for _, k := range keys {
			v := this.gpgContent.GetString(gSectionReversse, k, "")
			if v != "" {
				this.matches.insert(fmt.Sprintf(gReplaceKeyFmt, k), v, true)
			}
		}
		sort.Sort(&this.matches)
		exp := this.matches.expString()
		reg := regexp.MustCompile(exp)
		replacedCode := reg.ReplaceAllStringFunc(this.codeContent, func(src string) (rep string) {
			if v, ok := this.getMatch(src); ok {
				rep = v
			} else {
				fmt.Printf("[gogp]error: %s has no replacing\n", src)
				rep = src
				this.nNoReplaceMathNum++
			}
			return
		})
		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("[gogp]error:[%s].[%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.impName)
			fmt.Println(s)
			err = fmt.Errorf(s)
		}
		if err = this.saveGpFile(replacedCode, this.gpPath); err != nil { //save code to file
			return
		}
	} else {
		err = fmt.Errorf("[gogp]error:[%s] must have [%s] section", relateGoPath(this.gpgPath), gSectionReversse)
	}
	return
}

func (this *gopgProcessor) procGpg(file string, reverse bool) (err error) {
	this.gpContent = "" //clear gp content
	if err = this.loadGpgFile(file); err == nil {
		if reverse { //reverse work
			if this.gpgContent.Keys(gSectionReversse) != nil {
				if err = this.genCode(gSectionReversse); err != nil {
					return
				}
			}
		} else { //normal work
			if !gSilence {
				fmt.Printf(">[gogp]Processing:[%s]\n", relateGoPath(file))
			}
			for _, imp := range this.gpgContent.Sections() {
				if imp != gSectionReversse {
					if err = this.genCode(imp); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func (this *gopgProcessor) loadGpgFile(file string) (err error) {
	file = formatPath(file)
	this.gpPath = ""
	this.gpgPath = formatPath(file)
	this.gpgContent, err = ini.New(this.gpgPath)
	return
}

func (this *gopgProcessor) genCode(impName string) (err error) {
	if impName == gSectionIgnore { //never deal with this section
		return
	}
	if impName == gSectionReversse { //reverse only section, do reverse option
		if err = this.reverseProcess(); err != nil {
			//do nothing
		}
		return
	}
	this.impName = impName
	this.nNoReplaceMathNum = 0
	this.matches.clear() //clear matches
	if replaceList := this.gpgContent.Keys(impName); replaceList != nil {
		//make replace map
		for _, key := range replaceList {
			replace := this.gpgContent.GetString(impName, key, "")
			match := fmt.Sprintf(gReplaceKeyFmt, key)
			this.matches.insert(match, replace, false)
		}

		pathWithName := strings.TrimSuffix(this.gpgPath, gGpgExt)
		codePath := fmt.Sprintf("%s_%s_%s%s",
			pathWithName, gGpFileSuffix, impName, gCodeExt)
		gpPath := ""
		if gp, ok := this.getMatch(gkeyGpFilePath); ok { //read gp file from another path
			gpPath = filepath.Join(gGoPath, gp+gGpExt)
			this.gpPath = "" //clear gp content
		} else {
			gpPath = pathWithName + gGpExt
		}
		this.loadCodeFile(codePath) //load code file
		if this.gpPath != gpPath {  //load gp file if needed
			if err = this.loadGpFile(gpPath); err != nil {
				return
			}
		}
		//gen code file content
		replacedGp := gReplaceExp.ReplaceAllStringFunc(this.gpContent, func(src string) (rep string) {
			if v, ok := this.getMatch(src); ok {
				rep = v
			} else {
				fmt.Printf("[gogp]error: %s has no replacing\n", src)
				rep = src
				this.nNoReplaceMathNum++
			}
			return
		})
		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("[gogp]error:[%s].[%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), impName)
			fmt.Println(s)
			err = fmt.Errorf(s)
		}
		if err = this.saveCodeFile(replacedGp); err != nil { //save code to file
			return
		}
	}
	return
}

func (this *gopgProcessor) getMatch(key string) (match string, ok bool) {
	match, ok = this.matches.getMatch(key)
	return
}

func (this *gopgProcessor) loadGpFile(file string) (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(file); err == nil {
		this.gpPath = file
		this.gpContent = string(b)

		//ignore text format like "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END"
		this.gpContent = gGogpIgnoreExp.ReplaceAllString(this.gpContent, "")
	}
	return
}

func (this *gopgProcessor) loadCodeFile(file string) (err error) {
	var b []byte
	this.codeContent = ""
	this.codePath = file
	if b, err = ioutil.ReadFile(file); err == nil {
		this.codeContent = string(b)
	}
	return
}

func (this *gopgProcessor) getSrcFile(gp bool) string {
	if gp {
		return this.gpPath
	}
	return this.codePath
}

func (this *gopgProcessor) fileHead(gp bool) (h string) {
	tool := filepath.ToSlash(filepath.Dir(gThisFilePath))
	h = fmt.Sprintf(`///////////////////////////////////////////////////////////////////
//
// !!!!!!!!!!!! NEVER MODIFY THIS FILE MANUALLY !!!!!!!!!!!!
//
// This file was auto-generated by tool [%s]
// Last update at: [%s]
// Generate from:
//   [%s]
//   [%s] [%s]
//
// Tool [%s] info:
%s
///////////////////////////////////////////////////////////////////
`,
		tool,
		time.Now().Format("Mon Jan 02 2006 15:04:05"),
		relateGoPath(this.getSrcFile(gp)),
		relateGoPath(this.gpgPath),
		this.impName,
		tool,
		gCopyRightCode,
	)
	return
}

func (this *gopgProcessor) saveGpFile(body, gpFilePath string) (err error) {
	this.gpPath = gpFilePath
	if !gForceUpdate && this.loadGpFile(gpFilePath) == nil { //check if need update
		if this.gpContent == body { //body not change
			this.nSkipCodeFile++
			if !gSilence {
				fmt.Printf(">>[gogp][%s] skip\n", relateGoPath(this.gpPath))
			}
			return
		} else {
			//			fmt.Println("[%s]", xhash.MD5.StringHash(this.gpContent))
			//			fmt.Println("[%s]", xhash.MD5.StringHash(body))
			//			fmt.Println("[%s]", xhash.MD5.StringHash(strings.TrimSuffix(this.gpContent, body)))
		}
	}

	var fout *os.File
	if fout, err = os.OpenFile(this.gpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
		return
	}
	defer fout.Close()

	wt := bufio.NewWriter(fout)
	h := fmt.Sprintf(`//GOGP_IGNORE_BEGIN
%s//GOGP_IGNORE_END

`, this.fileHead(false))
	wt.WriteString(h)
	wt.WriteString(body)
	if err = wt.Flush(); err != nil {
		return
	}

	this.nCodeFile++
	if !gSilence {
		fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.gpPath))
	}
	return
}

func (this *gopgProcessor) saveCodeFile(body string) (err error) {
	if gForceUpdate || !strings.HasSuffix(this.codeContent, body) { //body change then save it,else skip it

		//		fmt.Println("[%s]", xhash.MD5.StringHash(this.codeContent))
		//		fmt.Println("[%s]", xhash.MD5.StringHash(body))
		//		fmt.Println("[%s]", xhash.MD5.StringHash(strings.TrimSuffix(this.codeContent, body)))

		var fout *os.File
		if fout, err = os.OpenFile(this.codePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
			return
		}
		defer fout.Close()
		wt := bufio.NewWriter(fout)

		wt.WriteString(this.fileHead(true))
		wt.WriteByte('\n')
		wt.WriteString(body)

		if err = wt.Flush(); err != nil {
			return
		}

		this.nCodeFile++
		if !gSilence {
			fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.codePath))
		}
	} else {
		this.nSkipCodeFile++
		if !gSilence {
			fmt.Printf(">>[gogp][%s] skip\n", relateGoPath(this.codePath))
		}
	}
	return
}
