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
	list        []*replaceCase
	match       map[string]string
	sectionName string
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
func (this *replaceList) expString() (exp string) {
	var b bytes.Buffer
	if this.Len() > 0 {
		for _, v := range this.list {
			s := fmt.Sprintf("\\Q%s\\E", v.value) //match raw letter
			b.WriteString(s)
			b.WriteByte('|')
		}
		b.Truncate(b.Len() - 1) //remove last '|'
		exp = b.String()
	} else {
		//avoid return "", which will match every byte
		exp = "\\QGOGP_DO_NOT_HAVE_ANY_KEY__\\E"
	}

	//fmt.Println(exp)
	return exp
}

func (this *replaceList) doReplacing(content, _path string, reverse bool) (rep string, noRep int) {
	reg := gGogpExpReplace
	if reverse {
		exp := this.expString()
		reg = regexp.MustCompile(exp)
	}
	rep = reg.ReplaceAllStringFunc(content, func(src string) (r string) {
		if v, ok := this.getMatch(src); ok {
			r = v
		} else {
			fmt.Printf("[gogp]  error: [%s] has no replacing.[%s:%s]\n", src, relateGoPath(_path), this.sectionName)
			//this.reportNoReplacing(src, _path)
			r = src
			noRep++
		}
		return
	})

	return
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
	gpgContent        *ini.IniFile //gpg file content
	gpContent         string
	codeContent       string
	impName           string         //current gpg section name
	step              gogp_proc_step //current processing step
	matches2          replaceList    //cases that need replacing, secondary
}

//get file suffix of code file
func (this *gopgProcessor) getCodeFileSuffix(section string) (r string) {
	if section == "" {
		section = this.impName
	}
	if v := this.gpgContent.GetString(section, grawKeyProductName, ""); v != "" {
		r = v
	} else {
		if v := this.gpgContent.GetString(section, grawKeyKeyType, ""); v != "" {
			l := strings.ToLower(v)
			if l != v {
				l = fmt.Sprintf("%s%s", l, get_hash(v))
			}
			r = l
		}
		if v := this.gpgContent.GetString(section, grawKeyValueType, ""); v != "" {
			l := strings.ToLower(v)
			if l != v {
				l = fmt.Sprintf("%s%s", l, get_hash(v))
			}
			if r == "" {
				r = l
			} else {
				r = fmt.Sprintf("%s_%s", r, l)
			}
		}
	}

	if r == "" {
		r = section
	} else {
		r = strings.Replace(r, "*", "#", -1)
	}

	return
}

func (this *gopgProcessor) reportNoReplacing(key, gpfile string) {
	fmt.Printf("[gogp] step%d error: [%s] has no replacing.[%s:%s %s]\n", this.step, key, relateGoPath(this.gpgPath), this.impName, gpfile)
}

//if has set key GOGP_Name, use it, else use section name
func (this *gopgProcessor) getGpName() (r string) {
	if name := this.gpgContent.GetString(this.impName, grawKeySrcPathName, ""); name != "" {
		r = strings.TrimSuffix(filepath.Base(name), gGpgExt)
	} else {
		r = "missing"
		fmt.Printf("error:[gogp]missing %s in %s:%s\n", grawKeySrcPathName, relateGoPath(this.gpgPath), this.impName)
	}
	return
}

func (this *gopgProcessor) checkGpgCfg(section, key string) (ok bool) {
	if v := this.gpgContent.GetString(section, key, ""); v == "true" || v == "1" { //if has ignore key
		ok = true
	}
	return
}

//check if a section is a valid task of step
func (this *gopgProcessor) isValidSection(section string, step gogp_proc_step) (ok bool) {
	if !strings.HasPrefix(section, gSectionIgnore) { //not an ignore section
		if checkReverse := strings.HasPrefix(section, gSectionReverse); checkReverse == step.IsReverse() { //if a proper section
			if !this.checkGpgCfg(section, grawKeyIgnore) { //if has ignore key
				ok = true
			}
		}
	}
	return
}

func (this *gopgProcessor) hasTask(step gogp_proc_step) bool {
	for _, imp := range this.gpgContent.Sections() {
		if this.isValidSection(imp, step) {
			return true
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

func (this *gopgProcessor) getReplist(second bool) *replaceList {
	pmatch := &this.matches
	if second {
		pmatch = &this.matches2
	}
	return pmatch
}

func (this *gopgProcessor) buildMatches(section string, reverse, second bool) (ok bool) {
	pmatch := this.getReplist(second)
	if !second {
		this.nNoReplaceMathNum = 0
	}
	pmatch.clear()
	if section == "" || section == "_" {
		section = this.impName
	}
	pmatch.sectionName = section
	if replaceList := this.gpgContent.Keys(section); replaceList != nil {
		//make replace map
		for _, key := range replaceList {
			replace := this.gpgContent.GetString(section, key, "")
			match := fmt.Sprintf(gReplaceKeyFmt, key)
			pmatch.insert(match, replace, reverse)
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

//require a gp file, maybe recursive
func (this *gopgProcessor) procRequireReplacement(statement string, nDepth int) (rep string, replaced bool, err error) {

	rep = statement

	if nDepth >= 10 {
		panic(fmt.Sprintf("[%s:%s]maybe loop recursive of #GOGP_REQUIRE(...), %d", relateGoPath(this.gpgPath), this.impName, nDepth))
	}

	elem := gGogpExpRequire.FindAllStringSubmatch(statement, -1)[0] //{"", "REQ", "REQP", "REQN"}
	req, reqp, reqn := elem[1], elem[2], elem[3]

	replaceSection := reqn
	if replaceSection == "" || replaceSection == "_" {
		replaceSection = this.impName
	}
	gpFullPath := this.getGpFullPath(reqp)
	gpContent := ""

	if gpContent, err = this.rawLoadFile(gpFullPath); err == nil {
		replacedGp := ""
		this.buildMatches(replaceSection, false, true)
		if replacedGp, err = this.doGpReplace(gpFullPath, gpContent, nDepth, true); err == nil {
			if this.step == gogp_step_PRODUCE {
				rep = "\n\n"
				replaced = true
				if reqn != "_" && !this.checkGpgCfg(replaceSection, grawKeyDontSave) { //reqn=="_" will not generate this code file
					codeFileSuffix := this.getCodeFileSuffix(replaceSection)

					gpgDir := filepath.Dir(this.gpgPath)
					gpName := strings.TrimSuffix(filepath.Base(gpFullPath), gGpExt)
					codePath := fmt.Sprintf("%s/%s.%s_%s%s",
						gpgDir, gpName, gGpCodeFileSuffix, codeFileSuffix, gCodeExt)

					if gRemoveProductsOnly { //remove products only
						this.nCodeFile++
						os.Remove(codePath)
						return
					}

					oldCode, _ := this.rawLoadFile(codePath)

					if gForceUpdate || !strings.HasSuffix(oldCode, replacedGp) { //body change then save it,else skip it
						codeContent := this.fileHead(false) + "\n" + replacedGp
						codeContent = goFmt(codeContent, codePath)
						if err = this.rawSaveFile(codePath, codeContent); err == nil {
							this.nCodeFile++
							if !gSilence {
								fmt.Printf(">>[gogp][%s] require ok\n", relateGoPath(codePath))
							}
						}
					} else {
						this.nSkipCodeFile++
						if !gSilence {
							fmt.Printf(">>[gogp][%s] require skip\n", relateGoPath(codePath))
						}
					}
				}
			} else {
				if gRemoveProductsOnly {
					rep = fmt.Sprintf("\n\n%s\n\n", req)
					replaced = true
				} else {
					if nDepth == 0 { //do not let require recursive
						replacedGp = strings.Replace(replacedGp, "package", "//package", -1) //comment package declaration
						replacedGp = strings.Replace(replacedGp, "import", "//import", -1)
						//reqSave := strings.Replace(req, "//#GOGP_REQUIRE", "//##GOGP_REQUIRE", -1)
						reqResult := fmt.Sprintf(gsTxtRequireResultFmt, reqp, "$CONTENT", reqp)
						out := fmt.Sprintf("\n\n%s\n%s\n\n", req, reqResult)
						replacedGp = gGogpExpTrimEmptyLine.ReplaceAllString(replacedGp, out)
						oldContent := gGogpExpTrimEmptyLine.ReplaceAllString(statement, "$CONTENT")

						rep = goFmt(replacedGp, this.gpPath)
						replaced = !strings.Contains(rep, oldContent) //check if content changed
						//if replaced {
						//	fmt.Printf("\n%#v\n%#v\n", rep, statement)
						//}

					} else {
						rep = "\n\n"
						replaced = true
					}
					//fmt.Printf("\n%#v\n%#v\n", rep, statement)
				}
			}
		}
	} else {
		fmt.Printf("[gogp][error][%s:%s] #GOGP_REQULRE(%s) : %s\n", relateGoPath(this.gpgPath), this.impName, reqp, err.Error())
		err = nil
		rep = statement
	}
	return
}

func (this *gopgProcessor) procStepRequire() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	codeFilePath := pathWithName + gCodeExt
	this.codePath = codeFilePath

	this.buildMatches(this.impName, false, false)

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	replcaceCnt := 0

	//match "//#GOGP_REQUIRE(path [, gpgSection])"
	replacedCode := gGogpExpRequire.ReplaceAllStringFunc(this.codeContent, func(src string) (rep string) {
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

	replacedCode = goFmt(replacedCode, this.gpPath)

	if replcaceCnt > 0 {
		if err = this.rawSaveFile(this.codePath, replacedCode); err != nil {
			fmt.Println(err)
		} else {
			if !gSilence {
				this.nCodeFile++
				fmt.Printf(">>[gogp]%s updated for require\n", relateGoPath(this.codePath))
			}
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

	if this.buildMatches(this.impName, true, false) {
		this.matches.sort()
		replacedCode, norep := this.matches.doReplacing(this.codeContent, this.codePath, true)
		this.nNoReplaceMathNum += norep

		replacedCode = gGogpExpEmptyLine.ReplaceAllString(replacedCode, "\n\n") //avoid multi empty lines

		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("[gogp]error:[%s].[%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.impName)
			//fmt.Println(s)
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
		gp = this.gpgContent.GetString(this.impName, grawKeySrcPathName, "") //read gp file from another path or name
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
		fmt.Printf("error:[gogp]missing %s in %s:%s\n", grawKeySrcPathName, relateGoPath(this.gpgPath), this.impName)
	}
	return gpPath
}

func (this *gopgProcessor) doGpReplace(_path, content string, nDepth int, second bool) (replacedGp string, err error) {
	replist := this.getReplist(second)
	_path = fmt.Sprintf("%s|%s", _path, filepath.Dir(this.gpgPath)) //gp file+gpg path=unique

	// match "//#GOGP_IFDEF cdk ... //#GOGP_ELSE ... //#GOGP_ENDIF" case
	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END
	// "//#GOGP_REQUIRE(path [, gpgSection])"
	replacedGp = gGogpExpPretreatAll.ReplaceAllStringFunc(content, func(src string) (rep string) {
		elem := gGogpExpPretreatAll.FindAllStringSubmatch(src, -1)[0] //{"", "IGNORE", "REQ", "REQP", "REQN", "CONDK", "T", "F","GPGCFG","ONCE"}
		ignore, req, reqp, reqn, condk, t, f, gpgcfg, once := elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7], elem[8], elem[9]
		switch {
		case ignore != "":
			rep = "\n\n"
		case condk != "":
			cfg := this.gpgContent.GetString(replist.sectionName, condk, "")
			if cfg == "" || cfg == "false" || cfg == "0" {
				rep = fmt.Sprintf("\n\n%s\n\n", f)
			} else {
				rep = fmt.Sprintf("\n\n%s\n\n", t)
			}
		case reqp != "":
			//require process
			if r, _, err := this.procRequireReplacement(src, nDepth+1); err == nil {
				rep = r
			} else {
				fmt.Println(err)
			}
			req, reqn = req, reqn //never use
		case gpgcfg != "":
			rep = this.gpgContent.GetString(replist.sectionName, gpgcfg, "")
		case once != "":
			if _, ok := gOnceMap[_path]; ok { //check if has processed this file
				rep = "\n\n"
			} else {
				rep = fmt.Sprintf("\n\n%s\n\n", once)
			}
		default:
			fmt.Println("error:[gogp]invalid predef statement", src)
		}
		return
	})

	//gen code file content
	norep := 0
	replacedGp, norep = replist.doReplacing(replacedGp, _path, false)
	this.nNoReplaceMathNum += norep

	//remove more empty line
	replacedGp = goFmt(replacedGp, this.gpgPath)

	if this.nNoReplaceMathNum > 0 { //report error
		s := fmt.Sprintf("[gogp]error:[%s:%s %s depth=%d] not every gp have been replaced\n", relateGoPath(this.gpgPath), relateGoPath(_path), replist.sectionName, nDepth)
		//fmt.Println(s)
		err = fmt.Errorf(s)
	}

	gOnceMap[_path] = true //record processed gp file

	return
}

func (this *gopgProcessor) procStepNormal() (err error) {
	//normal process
	if this.buildMatches(this.impName, false, false) {
		gpPath := this.getGpFullPath("")
		gpgDir := filepath.Dir(this.gpgPath)

		gpName := strings.TrimSuffix(filepath.Base(gpPath), gGpExt)
		codePath := fmt.Sprintf("%s/%s.%s_%s%s",
			gpgDir, gpName, gGpCodeFileSuffix, this.getCodeFileSuffix(this.impName), gCodeExt)

		this.loadCodeFile(codePath) //load code file, ignore error
		if this.gpPath != gpPath {  //load gp file if needed
			if err = this.loadGpFile(gpPath); err != nil {
				return
			}
		}

		replacedGp := ""
		if replacedGp, err = this.doGpReplace(this.gpPath, this.gpContent, 0, false); err != nil {
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
			fmt.Printf("[gogp]ReverseWork %d:[%s]\n", this.step, relateGoPath(this.gpgPath))
		} else {
			fmt.Printf(">[gogp]Processing:[%s]\n", relateGoPath(this.gpgPath))
		}
	}

	if !this.isValidSection(impName, this.step) { //not a valid section for this step, do nothing
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
	if err != nil {
		fmt.Printf("[gogp] error:[%s:%s step%d] [%s]\n", relateGoPath(this.gpgPath), this.impName, this.step, err.Error())
	}

	return
}

func (this *gopgProcessor) loadGpFile(file string) (err error) {
	this.gpContent = ""
	this.gpPath = file
	if this.gpContent, err = this.rawLoadFile(file); err == nil {
		//ignore text format like "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
		this.gpContent = gGogpExpIgnore.ReplaceAllString(this.gpContent, "")
	}
	return
}

func (this *gopgProcessor) loadCodeFile(file string) (err error) {
	this.codeContent = ""
	this.codePath = file
	if this.codeContent, err = this.rawLoadFile(file); err == nil {
		//do nothing
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
			//fmt.Println("[%s]", xhash.MD5.StringHash(this.gpContent))
			//fmt.Println("[%s]", xhash.MD5.StringHash(body))
			//fmt.Println("[%s]", xhash.MD5.StringHash(strings.TrimSuffix(this.gpContent, body)))
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
		//fmt.Printf("[%#v]\n", this.gpContent)
		//fmt.Printf("[%#v]\n", this.codeContent)
		//fmt.Printf("[%#v]\n", body)

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
