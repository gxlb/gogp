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
			rawName = (w == "<VALUE_TYPE>" || w == "<KEY_TYPE>") && (p == "." || s == ":")
			//fmt.Printf("[%s][%s][%s][%s][%v]\n", src, p, w, s, rawName)
		}
		if v, ok := this.getMatch(w); ok {
			if reverse {
				r = v
			} else { //gp replacing
				wv := v
				if rawName {
					wv = gGetRawName(v)
				}
				r = p + wv + s
				//fmt.Printf("[%s][%s]->[%s]\n", w, v, r)
			}
		} else {
			fmt.Printf("[gogp error]: [%s] has no replacing.[%s] [%s : %s]\n", w, relateGoPath(this.gpPath), relateGoPath(this.gpgPath), this.sectionName)
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
	replaces          replaceList    //keys that need replace
}

//get file suffix of code file
func (this *gopgProcessor) getCodeFileSuffix(section string) (r string) {
	if section == "" {
		section = this.impName
	}
	if v := this.getGpgCfg(section, grawKeyProductName, false); v != "" {
		r = v
	} else {
		if v := this.getGpgCfg(section, grawKeyKeyType, false); v != "" {
			l := strings.ToLower(v)
			if l != v {
				l = fmt.Sprintf("%s%s", l, get_hash(v))
			}
			r = l
		}
		if v := this.getGpgCfg(section, grawKeyValueType, false); v != "" {
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
	fmt.Printf("[gogp error]: %s [%s] has no replacing. [%s:%s %s]\n", this.step, key, relateGoPath(this.gpgPath), this.impName, gpfile)
}

//if has set key GOGP_Name, use it, else use section name
func (this *gopgProcessor) getGpName() (r string) {
	if name := this.getGpgCfg(this.impName, grawKeySrcPathName, true); name != "" {
		n := filepath.Base(name)
		idx := 0
		if idx = strings.Index(n, "."); idx < 0 { //split by first '.'
			idx = len(n)
		}
		r = n[:idx]
	} else {
		r = "missing"
		fmt.Printf("[gogp error]: missing %s in %s:%s\n", grawKeySrcPathName, relateGoPath(this.gpgPath), this.impName)
	}
	return
}

func (this *gopgProcessor) checkGpgCfg(section, key string) (ok bool) {
	if v := this.getGpgCfg(section, key, false); v != "" && v != "false" && v != "0" {
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
				fmt.Println(err)
				//return
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

func (this *gopgProcessor) buildMatches(section, gpPath string, reverse, second bool) (ok bool) {
	pmatch := this.getReplist(second)
	if !second {
		this.nNoReplaceMathNum = 0
	}
	pmatch.clear()
	if section == "" || section == "_" {
		section = this.impName
	}
	pmatch.sectionName = section
	pmatch.gpgPath = this.gpgPath
	pmatch.gpPath = gpPath
	if replaceList := this.gpgContent.Keys(section); replaceList != nil {
		//make replace map
		for _, key := range replaceList {
			replace := this.getGpgCfg(section, key, false)
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

func (this *gopgProcessor) getGpgCfg(section, key string, warnEmpty bool) (val string) {
	val = this.gpgContent.GetString(section, key, "")
	if val == "" && warnEmpty {
		fmt.Printf("[gogp warn]: [%s:%s] maybe lost key [%s]\n", relateGoPath(this.gpgPath), section, key)
	}
	return
}

func (this *gopgProcessor) remove(file string) {
	fmt.Printf(">>[gogp]: [%s] removed.\n", relateGoPath(file))
	os.Remove(file)
}

//require a gp file, maybe recursive
func (this *gopgProcessor) procRequireReplacement(statement, section string, nDepth int) (rep string, replaced bool, err error) {

	rep = statement
	if nDepth >= 5 {
		panic(fmt.Sprintf("[gogp error] [%s:%s]maybe loop recursive of #GOGP_REQUIRE(...), %d", relateGoPath(this.gpgPath), section, nDepth))
	}

	elem := gGogpExpRequire.FindAllStringSubmatch(statement, -1)[0] //{"", "REQ", "REQP", "REQN","REQGPG","CONTENT"}
	req, reqp, reqn, reqgpg, content := elem[1], elem[2], elem[3], elem[4], elem[5]

	if gDebug {
		fmt.Printf("[gogp debug] #GOGP_REQUIRE: [%s][%s]\n", reqp, reqn)
	}

	if reqgpg != "" && reqn == "" { //section name is config from gpg file
		reqn = this.getGpgCfg(section, reqgpg, true)
	}

	replaceSection := reqn
	leftFmt := gsTxtRequireResultFmt                          //left required file
	at := len(replaceSection) > 0 && replaceSection[0] == '@' //left result in this file
	if at {
		replaceSection = replaceSection[1:]
		leftFmt = gsTxtRequireAtResultFmt
	}
	sharp := len(replaceSection) > 0 && replaceSection[0] == '#' //do not save
	if sharp {
		replaceSection = replaceSection[1:]
	}
	if replaceSection == "" || replaceSection == "_" {
		replaceSection = section
	}
	gpFullPath := this.getGpFullPath(reqp)
	gpContent := ""

	if gpContent, err = this.rawLoadFile(gpFullPath); err == nil {
		replacedGp := ""
		if this.step == gogp_step_PRODUCE {
			replaced = true
			if at {
				rep = "\n" + content + "\n"
			} else {
				rep = "\n\n"
			}

			if !at && !sharp && reqn != "_" && !this.checkGpgCfg(replaceSection, grawKeyDontSave) { //reqn=="_" will not generate this code file
				gpgDir := filepath.Dir(this.gpgPath)
				gpName := strings.TrimSuffix(filepath.Base(gpFullPath), gGpExt)
				codePath := this.getProductFilePath(gpgDir, gpName, this.getCodeFileSuffix(replaceSection))

				if gRemoveProductsOnly { //remove products only
					this.nCodeFile++
					this.remove(codePath)
					return
				}

				if replacedGp, err = this.doGpReplace(gpFullPath, gpContent, replaceSection, nDepth, true); err != nil {
					return
				}

				if _, ok := gSavedCodeFile[codePath]; ok { //skip saved file
					//					if gDebug {
					//						fmt.Printf("[gogp] debug: step%d Required file [%s] skip\n", this.step, codePath)
					//					}
					return
				} else {
					gSavedCodeFile[codePath] = true //to prevent rewrite this file no matter it chages or not
					//					if gDebug {
					//						fmt.Printf("[gogp] debug: step%d Required file [%s] save ok\n", this.step, codePath)
					//					}
				}

				oldCode, _ := this.rawLoadFile(codePath)

				if gForceUpdate || !strings.HasSuffix(oldCode, replacedGp) { //body change then save it,else skip it
					codeContent := this.fileHead(gpFullPath, this.gpgPath, replaceSection) + "\n" + replacedGp
					codeContent = goFmt(codeContent, codePath)
					if err = this.rawSaveFile(codePath, codeContent); err == nil {
						this.nCodeFile++
						if !gSilence {
							fmt.Printf(">>[gogp] #GOGP_REQUIRE [%s:%s -> %s] ok\n", relateGoPath(gpFullPath), replaceSection, relateGoPath(codePath))
						}
					}
				} else {
					this.nSkipCodeFile++
					if !gSilence {
						fmt.Printf(">>[gogp] #GOGP_REQUIRE [%s:%s -> %s] skip\n", relateGoPath(gpFullPath), replaceSection, relateGoPath(codePath))
					}
				}
			}
		} else {
			if gRemoveProductsOnly {
				rep = fmt.Sprintf("\n\n%s\n\n", req)
				if !gSilence {
					fmt.Printf("%#v\n", rep)
				}
				replaced = true
			} else {
				if nDepth == 0 { //do not let require recursive
					if replacedGp, err = this.doGpReplace(gpFullPath, gpContent, replaceSection, nDepth, true); err != nil {
						return
					}
					//					if section == "GOGP_REVERSE_datadef" {
					//						fmt.Printf("@@procRequireReplacement replacedGp=[%s]\n", replacedGp)
					//					}
					replacedGp = strings.Replace(replacedGp, "package", "//package", -1) //comment package declaration
					replacedGp = strings.Replace(replacedGp, "import", "//import", -1)
					//reqSave := strings.Replace(req, "//#GOGP_REQUIRE", "//##GOGP_REQUIRE", -1)
					reqResult := fmt.Sprintf(leftFmt, reqp, "$CONTENT", reqp)
					out := fmt.Sprintf("\n\n%s\n%s\n\n", req, reqResult)
					replacedGp = gGogpExpTrimEmptyLine.ReplaceAllString(replacedGp, out)
					oldContent := gGogpExpTrimEmptyLine.ReplaceAllString(statement, "$CONTENT")

					rep = goFmt(replacedGp, this.gpPath)

					//check if content changed
					replaced = !strings.Contains(rep, oldContent) //|| !strings.Contains(oldContent, "//#GOGP_IGNORE_BEGIN")
					//					if section == "GOGP_REVERSE_datadef" {
					//						fmt.Printf("@@procRequireReplacement replaced=%v old=[%s] \nreplaced=%v replacedGp=[%s]\n", replaced, oldContent, replaced, rep)
					//					}
					//if replaced {
					//	fmt.Printf("\nrep=[%#v]\nold=[%#v]\n", rep, oldContent)
					//}

				} else {
					rep = "\n\n"
					replaced = true
				}
				//if gDebug {
				//	fmt.Printf("%#v %d\n%#v\n%#v\n", replaced, nDepth, statement, rep)
				//}
			}
		}

	} else {
		fmt.Printf("[gogp error]: [%s:%s] #GOGP_REQULRE(%s) : %s\n", relateGoPath(this.gpgPath), section, reqp, err.Error())
		err = nil
		rep = statement
	}
	//if gDebug {
	//	//fmt.Printf("%#v %d\n[%#v]\n[%#v]\n", replaced, nDepth, statement, rep)
	//}
	return
}

func (this *gopgProcessor) procStepRequire() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	codeFilePath := this.getFakeSrcFilePath(pathWithName)
	this.codePath = codeFilePath

	this.buildMatches(this.impName, "", false, false)

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	replcaceCnt := 0

	//match "//#GOGP_REQUIRE(path [, gpgSection])"
	replacedCode := gGogpExpRequireAll.ReplaceAllStringFunc(this.codeContent, func(src string) (rep string) {
		elem := gGogpExpRequireAll.FindAllStringSubmatch(src, -1)[0] //{"","REQ", "REQP", "REQN","REQGPG","FILEB","OPEN","FILEE"}
		req, reqp, reqn, reqgpg, content, fileb, open, filee := elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7], elem[8]

		reqp, reqn, reqgpg, content = reqp, reqn, reqgpg, content //avoid compile error

		var err error
		var replaced bool
		switch {
		case req != "":
			if rep, replaced, err = this.procRequireReplacement(src, this.impName, 0); err != nil {
				fmt.Println(err)
			}
		case fileb != "":
			if gRemoveProductsOnly {
				rep, replaced = fmt.Sprintf("\n\n%s\n\n", fileb), true
				break
			}
			repContent := fmt.Sprintf(gsTxtGogpIgnoreFmt, " ///gogp_file_begin\n", gGetTxtFileBeginContent(open != ""), " ///gogp_file_begin\n\n")
			if rep, replaced = src, !strings.Contains(src, gGogpExpTrimEmptyLine.ReplaceAllString(repContent, "$CONTENT")); replaced {
				rep = fmt.Sprintf("\n\n%s\n%s", fileb, repContent)
			}
		case filee != "":
			if gRemoveProductsOnly {
				rep, replaced = fmt.Sprintf("\n\n%s\n\n", filee), true
				break
			}
			repContent := fmt.Sprintf(gsTxtGogpIgnoreFmt, " ///gogp_file_end\n", gsTxtFileEndContent, " ///gogp_file_end\n\n")
			if rep, replaced = src, !strings.Contains(src, gGogpExpTrimEmptyLine.ReplaceAllString(repContent, "$CONTENT")); replaced {
				rep = fmt.Sprintf("\n\n%s\n%s", filee, repContent)
			}
		}
		if replaced {
			replcaceCnt++
		}
		return
	})

	if gForceUpdate || replcaceCnt > 0 {
		replacedCode = gGogpExpEmptyLine.ReplaceAllString(replacedCode, "\n\n") //avoid multi empty lines
		replacedCode = goFmt(replacedCode, this.gpPath)

		if err = this.rawSaveFile(this.codePath, replacedCode); err == nil {
			this.nCodeFile++
			if !gSilence {
				fmt.Printf(">>[gogp] %s updated for #GOGP_REQUIRE\n", relateGoPath(this.codePath))
			}
		} else {
			fmt.Println(err)
		}
	}

	return
}

func (this *gopgProcessor) getFakeSrcFilePath(pathWithName string) string {
	return fmt.Sprintf("%s.%s%s", pathWithName, gGpCodeFileSuffix, gCodeExt)
}

func (this *gopgProcessor) getProductFilePath(gpgDir, gpName, codeFileSuffix string) string {
	return fmt.Sprintf("%s/%s.%s_%s%s", gpgDir, gpName, gGpCodeFileSuffix, codeFileSuffix, gCodeExt)

}

func (this *gopgProcessor) procStepReverse() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	gpFilePath := pathWithName + gGpExt
	codeFilePath := this.getFakeSrcFilePath(pathWithName)
	this.codePath = codeFilePath
	this.gpPath = gpFilePath

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	//ignore text format like "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	this.codeContent = gGogpExpReverseIgnoreAll.ReplaceAllString(this.codeContent, "\n\n")

	if this.buildMatches(this.impName, this.gpPath, true, false) {
		this.matches.sort()
		replacedCode, norep := this.matches.doReplacing(this.codeContent, this.gpgPath, true)
		this.nNoReplaceMathNum += norep

		replacedCode = gGogpExpEmptyLine.ReplaceAllString(replacedCode, "\n\n") //avoid multi empty lines

		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("[gogp error]: [%s:%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.impName)
			//fmt.Println(s)
			err = fmt.Errorf(s)
		}

		if err = this.saveGpFile(replacedCode, this.gpPath); err != nil { //save code to file
			return
		}
	} else {
		err = fmt.Errorf("[gogp error]: [%s] must have [%s] section", relateGoPath(this.gpgPath), gSectionReverse)
	}
	return
}

func (this *gopgProcessor) getGpFullPath(gp string) string {
	gpPath := ""
	gpgDir := filepath.Dir(this.gpgPath)
	if "" == gp {
		gp = this.getGpgCfg(this.impName, grawKeySrcPathName, false) //read gp file from another path or name
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
		fmt.Printf("[gogp error]: missing [%s] in [%s:%s]\n", grawKeySrcPathName, relateGoPath(this.gpgPath), this.impName)
	}
	return gpPath
}

func (this *gopgProcessor) doPredefReplace(gpPath, content, section string, nDepth int) (rep string) {
	pathIdentify := fmt.Sprintf("%s|%s", relateGoPath(gpPath), relateGoPath(filepath.Dir(this.gpgPath))) //gp file+gpg path=unique
	this.replaces.clear()

	// match "//#GOGP_IFDEF cdk ... //#GOGP_ELSE ... //#GOGP_ENDIF" case
	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END
	// "//#GOGP_REQUIRE(path [, gpgSection])"
	for _content, needReplace, i := content, true, 0; needReplace && i < 3; _content, i = rep, i+1 {
		needReplace = false
		rep = gGogpExpPretreatAll.ReplaceAllStringFunc(_content, func(src string) (_rep string) {
			elem := gGogpExpPretreatAll.FindAllStringSubmatch(src, -1)[0] //{"", "IGNORE", "REQ", "REQP", "REQN", "REQGPG","CONDK", "T", "F","GPGCFG","ONCE"}
			ignore, req, reqp, reqn, reqgpg, reqcontent, gpgcfg, once, repsrc, repdst := elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7], elem[8], elem[9], elem[10]

			if reqgpg != "" && reqn == "" { //section name is config from gpg file
				reqn = this.getGpgCfg(section, reqgpg, true)
			}

			if !gSilence && i > 1 {
				fmt.Printf("##src=[%#v]\n i=%d ignore=[%s] req=[%s] reqp=[%s] reqn=[%s] reqgpg=[%s] gpgcfg=[%s] once=[%s] repsrc=[%s] repdst=[%s]\n",
					src, i, ignore, req, reqp, reqn, reqgpg, gpgcfg, once, repsrc, repdst)
			}

			needReplace = true

			switch {
			case ignore != "":
				_rep = "\n\n"

			case reqp != "":
				if reqcontent == "" {
					//require process
					if r, _, err := this.procRequireReplacement(src, section, nDepth+1); err == nil {
						_rep = r
					} else {
						fmt.Println(err)
					}
					req, reqn, reqcontent = req, reqn, reqcontent //never use
				} else {
					_rep = reqcontent
				}

			case gpgcfg != "":
				_rep = this.getGpgCfg(section, gpgcfg, true)
			case once != "":
				if _, ok := gOnceMap[pathIdentify]; ok { //check if has processed this file
					_rep = "\n\n"
					if gDebug {
						fmt.Printf("[gogp debug]: %s GOGP_ONCE(%s:%s) ignore [%#v]\n", this.step, pathIdentify, section, once)
					}

				} else {
					_rep = fmt.Sprintf("\n\n%s\n\n", once)
					if gDebug {
						fmt.Printf("[gogp debug]: %s GOGP_ONCE(%s:%s) ok [%#v]\n", this.step, pathIdentify, section, once)
					}
				}
			case repsrc != "":
				_rep = ""
				this.replaces.insert(repdst, repsrc, true)
				//fmt.Printf("%s %s %s [%s] -> [%s]\n", gpPath, section, src, repsrc, repdst)

			default:
				fmt.Printf("[gogp error]: %s invalid predef statement [%#v]\n", this.step, src)
			}
			//			if section == "GOGP_REVERSE_datadef" {
			//				fmt.Printf("##gpPath=[%s] section[%s]\n##%s 2src=[%s]\n##%s 3rep=[%s]##4%s rep=[%s]\n",
			//					gpPath, section, section, src, section, _rep, section, reqcontent)
			//				//				fmt.Printf("ignore=[%s] req=[%s] reqp=[%s] reqn=[%s] reqgpg=[%s] gpgcfg=[%s] once=[%s] repsrc=[%s] repdst=[%s]\n",
			//				//					ignore, req, reqp, reqn, reqgpg, gpgcfg, once, repsrc, repdst)
			//			}
			return
		})
		//		if section == "GOGP_REVERSE_datadef" {
		//			fmt.Printf("@@i=%d _content=[%s]\ni=%d rep=[%s]\n", i, _content, i, rep)
		//		}
	}

	if this.step == gogp_step_PRODUCE { //prevent gen #GOGP_ONCE code twice when gen code
		gOnceMap[pathIdentify] = true //record processed gp file
	}

	return
}

func (this *gopgProcessor) pretreatGpForCode(gpContent string, section string) (replaced string) {
	replaced = gGogpExpCodeIgnore.ReplaceAllStringFunc(gpContent, func(src string) (rep string) {
		elem := gGogpExpCodeIgnore.FindAllStringSubmatch(src, -1)[0] //{"", "IGNORE", "GPONLY", "CONDK", "T", "F"}
		ignore, gponly, condk, t, f := elem[1], elem[2], elem[3], elem[4], elem[5]
		switch {
		case condk != "":
			cfg := this.getGpgCfg(section, condk, false)

			sel := t
			if cfg == "" || cfg == "false" || cfg == "0" {
				sel = f
			}
			sel = strings.Replace(sel, grawStringNotComment, "", -1) //uncomment selected
			rep = fmt.Sprintf("\n\n%s\n\n", sel)
		default:
		case ignore != "" || gponly != "":
			rep = "\n\n"
		}
		return
	})
	return
}

func (this *gopgProcessor) doGpReplace(gpPath, content, section string, nDepth int, second bool) (replacedGp string, err error) {
	_path := fmt.Sprintf("%s|%s", relateGoPath(gpPath), relateGoPath(filepath.Dir(this.gpgPath))) //gp file+gpg path=unique

	replacedGp = content
	this.replaces.clear()

	if this.step == gogp_step_PRODUCE {
		//replacedGp = gGogpExpCodeIgnore.ReplaceAllString(replacedGp, "\n\n")
		replacedGp = this.pretreatGpForCode(replacedGp, section)
	}

	replacedGp = this.doPredefReplace(gpPath, replacedGp, section, nDepth)

	//replaces keys that need be replacing
	if this.replaces.Len() > 0 {
		replacedGp, _ = this.replaces.doReplacing(replacedGp, this.gpgPath, true)
		this.replaces.clear()
	}

	//gen code file content
	this.buildMatches(section, gpPath, false, second)
	replist := this.getReplist(second)
	norep := 0
	replacedGp, norep = replist.doReplacing(replacedGp, this.gpgPath, false)
	this.nNoReplaceMathNum += norep

	replacedGp = gGogpExpEmptyLine.ReplaceAllString(replacedGp, "\n\n") //avoid multi empty lines

	//remove more empty line
	replacedGp = goFmt(replacedGp, this.gpgPath)

	if this.nNoReplaceMathNum > 0 { //report error
		s := fmt.Sprintf("[gogp error]: [%s:%s %s depth=%d] not every gp have been replaced\n", relateGoPath(this.gpgPath), relateGoPath(_path), replist.sectionName, nDepth)
		//fmt.Println(s)
		err = fmt.Errorf(s)
	}

	//	if section == "GOGP_REVERSE_datadef" {
	//		fmt.Printf("@@doGpReplace replacedGp=[%s]\n", replacedGp)
	//	}

	return
}

func (this *gopgProcessor) procStepNormal() (err error) {
	//normal process
	gpPath := this.getGpFullPath("")
	gpgDir := filepath.Dir(this.gpgPath)

	gpName := strings.TrimSuffix(filepath.Base(gpPath), gGpExt)
	codePath := this.getProductFilePath(gpgDir, gpName, this.getCodeFileSuffix(this.impName))

	this.loadCodeFile(codePath) //load code file, ignore error
	if this.gpPath != gpPath {  //load gp file if needed
		if err = this.loadGpFile(gpPath); err != nil {
			return
		}
	}

	replacedGp := ""
	if replacedGp, err = this.doGpReplace(this.gpPath, this.gpContent, this.impName, 0, false); err != nil {
		return
	}

	if err = this.saveCodeFile(replacedGp); err != nil { //save code to file
		return
	}

	return
}

//gen code or gp file
func (this *gopgProcessor) genProduct(id int, impName string) (err error) {
	if 0 == id && !gSilence {
		fmt.Printf(">[gogp] %s: [%s]\n", this.step, relateGoPath(this.gpgPath))
	}

	if !this.isValidSection(impName, this.step) { //not a valid section for this step, do nothing
		return
	}

	this.impName = impName

	if !gSilence {
		fmt.Printf(">[gogp] %s [%s:%s] \n", this.step, relateGoPath(this.gpgPath), this.impName)
	}

	switch this.step {
	case gogp_step_REQUIRE:
		err = this.procStepRequire()
	case gogp_step_REVERSE:
		err = this.procStepReverse()
	case gogp_step_PRODUCE:
		err = this.procStepNormal()
	}
	if err != nil {
		fmt.Printf("[gogp error]: %s [%s:%s] [%s]\n", this.step, relateGoPath(this.gpgPath), this.impName, err.Error())
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

func (this *gopgProcessor) fileHead(srcFile, gpgFile, section string) (h string) {
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
		relateGoPath(srcFile),
		relateGoPath(gpgFile),
		section,
		tool,
		gCopyRightCode,
	)
	return
}

func (this *gopgProcessor) saveGpFile(body, gpFilePath string) (err error) {
	this.gpPath = gpFilePath
	if gRemoveProductsOnly { //remove products only
		this.nCodeFile++
		this.remove(this.gpPath)
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

`, this.fileHead(this.codePath, this.gpgPath, this.impName))
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
		this.remove(this.codePath)
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

		wt.WriteString(this.fileHead(this.gpPath, this.gpgPath, this.impName))
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
