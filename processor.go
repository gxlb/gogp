//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

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

func (this *replaceList) sort() {
	sort.Sort(this)
}

func (this *replaceList) insert(k, v string, reverse bool) int {
	return this.push(&replaceCase{key: k, value: v}, reverse)
}

func (this *replaceList) push(v *replaceCase, reverse bool) int {
	if strings.HasPrefix(v.key, gKeyReservePrefix) { //do not match reserved keys
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
func (this *replaceList) expString() string {
	var b bytes.Buffer
	for _, v := range this.list {
		s := fmt.Sprintf("\\Q%s\\E", v.value) //match raw letter
		b.WriteString(s)
		b.WriteByte('|')
	}
	if this.Len() > 0 {
		b.Truncate(b.Len() - 1) //remove last '|'
	}

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
	step              gogp_proc_step
}

//func (this *gopgProcessor) reverseWork(gpgFilePath string) (err error) {
//
//	if !strings.HasSuffix(gpgFilePath, gGpgExt) { //must .gpg file
//		err = fmt.Errorf("[gogp]error:[%s] must be %s file at reverse mode", relateGoPath(gpgFilePath), gGpgExt)
//		return
//	}
//
//	gpgFullPath := expadGoPath(gpgFilePath) //make full path
//
//	if err = this.procGpg(gpgFullPath, true); err != nil {
//		return
//	}
//
//	if this.nCodeFile+this.nSkipCodeFile <= 0 { //no reverse tasks
//		err = fmt.Errorf("[gogp]error:[%s] must has %s leaded sections", relateGoPath(gpgFilePath), gSectionReverse)
//	}
//
//	return
//}

//if has set key GOGP_Name, use it, else use section name
func (this *gopgProcessor) getGpName() (r string) {
	if name := this.gpgContent.GetString(this.impName, grawKeyName, ""); name != "" {
		r = strings.TrimSuffix(name, gGpgExt)
	} else {
		r = strings.TrimSuffix(filepath.Base(this.gpgPath), gGpgExt)
	}
	return
}

func (this *gopgProcessor) hasTask(step gogp_proc_step) bool {
	for _, imp := range this.gpgContent.Sections() {
		if !strings.HasPrefix(imp, gSectionIgnore) {
			if checkReverse := strings.HasPrefix(imp, gSectionReverse); checkReverse == step.IsReverse() {
				return true
			}
		}
	}
	return false
}

func (this *gopgProcessor) procGpg(file string, step gogp_proc_step) (err error) {
	this.gpContent = "" //clear gp content
	this.step = step
	if err = this.loadGpgFile(file); err == nil && this.hasTask(step) {
		for i, imp := range this.gpgContent.Sections() {
			if err = this.genProduct(i, imp); err != nil {
				return
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

//if has set key GOGP_Name, use it, else use section name
func (this *gopgProcessor) getCodeSuffix() (r string) {
	if name := this.gpgContent.GetString(this.impName, grawKeyName, ""); name != "" {
		r = name
	} else {
		r = this.impName
	}
	return
}

func (this *gopgProcessor) buildMatches(reverse bool) (ok bool) {
	this.matches.clear() //clear matches
	this.nNoReplaceMathNum = 0
	if replaceList := this.gpgContent.Keys(this.impName); replaceList != nil {
		//make replace map
		for _, key := range replaceList {
			replace := this.gpgContent.GetString(this.impName, key, "")
			match := fmt.Sprintf(gReplaceKeyFmt, key)
			this.matches.insert(match, replace, reverse)
		}
		ok = true
	}
	return
}

func (this *gopgProcessor) rawLoadFile(file string) (content string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(file); err == nil {
		//deal with new line
		b = bytes.Replace(b, []byte("\r\n"), []byte("\n"), -1)

		content = string(b)
	}
	return
}
func (this *gopgProcessor) rawSaveFile(file, content string) (err error) {
	var fout *os.File
	if fout, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
		return
	}
	defer fout.Close()

	fout.WriteString(content)
	fout.Sync()
	return
}

//func (this *gopgProcessor) procRequireOp(statement, reqHead, reqGpPath, codeFileSuffix string, save bool, nDepth int) (rep string, err error) {
//	if this.step == gogp_step_REQUIRE && reqHead == "##" {
//		rep = statement
//		return
//	}

//	gpFullPath := this.getGpFullPath(reqGpPath)
//	gpContent := ""
//	if gpContent, err = this.rawLoadFile(gpFullPath); err == nil {
//		replacedGp := ""
//		if replacedGp, err = this.doGpReplace(gpContent, nDepth); err == nil {
//			if save {
//				rep = "\n\n"
//				if codeFileSuffix == "" {
//					if v, ok := this.getMatch("<VALUE_TYPE>"); ok {
//						l := strings.ToLower(v)
//						if l != v {
//							l = fmt.Sprintf("%s_%s", l, get_hash(v))
//						}
//						codeFileSuffix = l
//					}
//				}

//				gpgDir := filepath.Dir(this.gpgPath)
//				gpName := strings.TrimSuffix(filepath.Base(gpFullPath), gGpExt)
//				codePath := fmt.Sprintf("%s/%s.%s_%s%s",
//					gpgDir, gpName, gGpCodeFileSuffix, this.getCodeSuffix(), gCodeExt)
//				if err = this.rawSaveFile(codePath, replacedGp); err == nil {
//					//todo
//				}

//			} else {
//				replacedGp = strings.Replace(replacedGp, "package", "//package", -1) //comment package declaration
//				statement = strings.Replace(statement, "//#GOGP_REQUIRE", "//##GOGP_REQUIRE", -1)
//				rep = fmt.Sprintf("%s\n//#GOGP_IGNORE_BEGIN\n%s\n//#GOGP_IGNORE_END\n", statement, replacedGp)
//				rep = gGogpExpEmptyLine.ReplaceAllString(rep, "\n\n") //remove more empty lines
//			}
//		}
//	}
//	return
//}

//require a gp file, maybe recursive
func (this *gopgProcessor) procRequireReplacement(statement string, nDepth int) (rep string, replaced bool, err error) {
	if nDepth >= 10 {
		panic(fmt.Sprintf("[%s:%s]maybe loop recursive of require, %d", relateGoPath(this.gpgPath), this.impName, nDepth))
	}

	elem := gGogpExpRequire.FindAllStringSubmatch(statement, -1)[0] //{"", "REQH", "REQP", "REQN",}
	reqh, reqp, reqn := elem[1], elem[2], elem[3]

	if this.step == gogp_step_REQUIRE && reqh == "##" { //ignore replaced require
		rep = statement
		return
	} else {
		replaced = true
	}

	codeFileSuffix := reqn
	gpFullPath := this.getGpFullPath(reqp)
	gpContent := ""
	if gpContent, err = this.rawLoadFile(gpFullPath); err == nil {
		replacedGp := ""
		if replacedGp, err = this.doGpReplace(gpContent, nDepth); err == nil {
			if this.step == gogp_step_PRODUCE {
				rep = "\n\n"
				if codeFileSuffix == "" {
					if v, ok := this.getMatch("<VALUE_TYPE>"); ok {
						l := strings.ToLower(v)
						if l != v {
							l = fmt.Sprintf("%s_%s", l, get_hash(v))
						}
						codeFileSuffix = l
					}
					if codeFileSuffix == "" {
						codeFileSuffix = "unknown"
					}
				}

				gpgDir := filepath.Dir(this.gpgPath)
				gpName := strings.TrimSuffix(filepath.Base(gpFullPath), gGpExt)
				codePath := fmt.Sprintf("%s/%s.%s_%s%s",
					gpgDir, gpName, gGpCodeFileSuffix, codeFileSuffix, gCodeExt)
				if gRemoveProductsOnly { //remove products only
					this.nCodeFile++
					os.Remove(this.codePath)
					return
				}
				oldCode, _ := this.rawLoadFile(codePath)
				if gForceUpdate || !strings.HasSuffix(oldCode, replacedGp) { //body change then save it,else skip it
					codeContent := this.fileHead(false) + "\n" + replacedGp
					if err = this.rawSaveFile(codePath, codeContent); err == nil {
						this.nCodeFile++
						if !gSilence {
							fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.codePath))
						}
					}
				} else {
					this.nSkipCodeFile++
					if !gSilence {
						fmt.Printf(">>[gogp][%s] skip\n", relateGoPath(this.codePath))
					}
				}
			} else {
				replacedGp = strings.Replace(replacedGp, "package", "//package", -1) //comment package declaration
				statement = strings.Replace(statement, "//#GOGP_REQUIRE", "//##GOGP_REQUIRE", -1)
				rep = fmt.Sprintf("%s\n//#GOGP_IGNORE_BEGIN\n%s\n//#GOGP_IGNORE_END\n", statement, replacedGp)
				rep = gGogpExpEmptyLine.ReplaceAllString(rep, "\n\n") //remove more empty lines
			}
		}
	}
	return

	//	switch reqh {
	//	case "#":
	//		rep, err = this.procRequireOp(statement, reqp, reqn, false, nDepth)
	//	case "##":
	//		rep, err = this.procRequireOp(statement, reqp, reqn, true, nDepth)
	//	}
	//	return
}

func (this *gopgProcessor) procStepRequire() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	codeFilePath := pathWithName + gCodeExt
	this.codePath = codeFilePath

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	replcaceCnt := 0

	//match "//#GOGP_REQUIRE(path [, nameSuffix])"
	replacedCode := gGogpExpRequire.ReplaceAllStringFunc(this.gpContent, func(src string) (rep string) {
		var err error
		var replaced bool
		if rep, replaced, err = this.procRequireReplacement(src, 0); err != nil {
			fmt.Println(err)
		}
		if replaced {
			replcaceCnt++
		}
		return
	})

	if replcaceCnt > 0 {
		if err = this.rawSaveFile(this.codePath, replacedCode); err != nil {
			fmt.Println(err)
		}
	}
	return
}

func (this *gopgProcessor) procStepReverse() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	gpFilePath := pathWithName + gGpExt
	codeFilePath := pathWithName + gCodeExt
	this.codePath = codeFilePath
	this.gpPath = gpFilePath

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	//ignore text format like "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	this.codeContent = gGogpExpIgnore.ReplaceAllString(this.codeContent, "\n\n")

	if this.buildMatches(true) {
		this.matches.sort()
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
		err = fmt.Errorf("[gogp]error:[%s] must have [%s] section", relateGoPath(this.gpgPath), gSectionReverse)
	}
	return
}

func (this *gopgProcessor) getGpFullPath(gp string) string {
	gpPath := ""
	gpgDir := filepath.Dir(this.gpgPath)
	if "" == gp {
		gp = this.gpgContent.GetString(this.impName, grawKeyGpFilePath, "") //read gp file from another path or name
	}
	if gp != "" { //read gp file from another path or name
		if !strings.HasPrefix(gp, gGpExt) {
			gp += gGpExt
		}
		if p, _ := filepath.Split(gp); p == "" || '.' == gp[0] { //if only config gp name, or lead with ".", use gpg dir
			gpPath = filepath.Join(gpgDir, gp)
		} else {
			gpPath = filepath.Join(gGoPath, gp)
		}
	} else {
		panic(fmt.Sprintf("error:[gogp]missing %s in %s:%s", grawKeyGpFilePath, relateGoPath(this.gpgPath), this.impName))
	}

	return gpPath
}

func (this *gopgProcessor) doGpReplace(content string, nDepth int) (rep string, err error) {
	// match "//#GOGP_IFDEF cdk ... //#GOGP_ELSE ... //#GOGP_ENDIF" case
	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END
	// "//#GOGP_REQUIRE(path [, nameSuffix])"
	replacedGp := gGogpExpPretreatAll.ReplaceAllStringFunc(content, func(src string) (rep string) {
		elem := gGogpExpPretreatAll.FindAllStringSubmatch(src, -1)[0] //{"", "IGNORE", "REQH", "REQP", "REQN", "CONDK", "T", "F"}
		ignore, reqh, reqp, reqn, condk, t, f := elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7]
		switch {
		case ignore != "":
			rep = "\n\n"
		case condk != "":
			cfg := this.gpgContent.GetString(this.impName, condk, gFalse)
			if cfg == gFalse || cfg == "0" {
				rep = fmt.Sprintf("\n\n%s\n\n", f)
			} else {
				rep = fmt.Sprintf("\n\n%s\n\n", t)
			}
		case reqp != "":
			//require process
			fmt.Println("[gogp]todo:", this.gpPath, "require", reqp)
			if r, _, err := this.procRequireReplacement(src, nDepth+1); err == nil {
				rep = r
			} else {
				fmt.Println(err)
			}
			reqh, reqn = reqn, reqh //never use

		default:
			fmt.Println("[gogp]invalid predef statement", src)
		}
		return
	})

	//gen code file content
	replacedGp = gGogpExpReplace.ReplaceAllStringFunc(replacedGp, func(src string) (rep string) {
		if v, ok := this.getMatch(src); ok {
			rep = v
		} else {
			fmt.Printf("[gogp]error: %s has no replacing\n", src)
			rep = src
			this.nNoReplaceMathNum++
		}
		return
	})

	//remove more empty line
	replacedGp = gGogpExpEmptyLine.ReplaceAllString(replacedGp, "\n\n")

	if this.nNoReplaceMathNum > 0 { //report error
		s := fmt.Sprintf("[gogp]error:[%s].[%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.impName)
		fmt.Println(s)
		err = fmt.Errorf(s)
	}

	return
}

func (this *gopgProcessor) procStepNormal() (err error) {
	//normal process
	if this.buildMatches(false) {
		gpPath := this.getGpFullPath("")
		gpgDir := filepath.Dir(this.gpgPath)

		gpName := strings.TrimSuffix(filepath.Base(gpPath), gGpExt)
		codePath := fmt.Sprintf("%s/%s.%s_%s%s",
			gpgDir, gpName, gGpCodeFileSuffix, this.getCodeSuffix(), gCodeExt)

		this.loadCodeFile(codePath) //load code file, ignore error
		if this.gpPath != gpPath {  //load gp file if needed
			if err = this.loadGpFile(gpPath); err != nil {
				return
			}
		}

		replacedGp := ""
		if replacedGp, err = this.doGpReplace(this.gpContent, 0); err != nil {
			return
		}

		if err = this.saveCodeFile(replacedGp); err != nil { //save code to file
			return
		}
	}
	return
}

//gen code or gp file
func (this *gopgProcessor) genProduct(id int, impName string) (err error) {
	if 0 == id && !gSilence {
		if this.step.IsReverse() {
			fmt.Printf("[gogp]ReverseWork:[%s]\n", relateGoPath(this.gpgPath))
		} else {
			fmt.Printf(">[gogp]Processing:[%s]\n", relateGoPath(this.gpgPath))
		}
	}

	if strings.HasPrefix(impName, gSectionIgnore) { //never deal with this section
		return
	}

	checkReverse := strings.HasPrefix(impName, gSectionReverse)
	if checkReverse != this.step.IsReverse() { //not proper section, do nothing
		return
	}

	this.impName = impName

	switch this.step {
	case gogp_step_REQUIRE:
		err = this.procStepRequire()
	case gogp_step_REVERSE:
		err = this.procStepReverse()
	case gogp_step_PRODUCE:
		err = this.procStepNormal()
	}

	return
}

func (this *gopgProcessor) getMatch(key string) (match string, ok bool) {
	match, ok = this.matches.getMatch(key)
	return
}

func (this *gopgProcessor) loadGpFile(file string) (err error) {
	this.gpContent = ""
	if this.gpContent, err = this.rawLoadFile(file); err == nil {
		this.gpPath = file
		//ignore text format like "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
		this.gpContent = gGogpExpIgnore.ReplaceAllString(this.gpContent, "")
	}
	return
}

func (this *gopgProcessor) loadCodeFile(file string) (err error) {
	this.codeContent = ""
	if this.codeContent, err = this.rawLoadFile(file); err == nil {
		this.codePath = file
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
	if gRemoveProductsOnly { //remove products only
		this.nCodeFile++
		os.Remove(this.gpPath)
		return
	}
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
	h := fmt.Sprintf(`//#GOGP_IGNORE_BEGIN
%s//#GOGP_IGNORE_END

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
	if gRemoveProductsOnly { //remove products only
		this.nCodeFile++
		os.Remove(this.codePath)
		return
	}
	if gForceUpdate || !strings.HasSuffix(this.codeContent, body) { //body change then save it,else skip it
		//		fmt.Printf("[%#v]\n", this.gpContent)
		//		fmt.Printf("[%#v]\n", this.codeContent)
		//		fmt.Printf("[%#v]\n", body)

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
