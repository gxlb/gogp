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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gxlb/gogp/ini"
)

const maxRecursionDepth = 3

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
	impName           string          //current gpg section name
	step              gogpProcessStep //current processing step
	matches2          replaceList     //cases that need replacing, secondary
	replaces          replaceList     //keys that need replace
	maps              replaceList     //keys that need replace
}

func (this *gopgProcessor) procGpg(file string, step gogpProcessStep) (err error) {
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

//gen code or gp file
func (this *gopgProcessor) genProduct(id int, impName string) (err error) {
	if 0 == id && !optSilence {
		fmt.Printf(">[gogp] %s: [%s]\n", this.step, relateGoPath(this.gpgPath))
	}

	if !this.isValidSection(impName, this.step) { //not a valid section for this step, do nothing
		return
	}

	this.impName = impName

	if !optSilence {
		fmt.Printf(">[gogp] %s [%s:%s] \n", this.step, relateGoPath(this.gpgPath), this.impName)
	}

	switch this.step {
	case gogpStepREQUIRE:
		err = this.procStep1Require()
	case gogpStepREVERSE:
		err = this.procStep2Reverse()
	case gogpStepPRODUCE:
		err = this.procStep3Produce()
	}
	if err != nil {
		fmt.Printf("[gogp error]: %s [%s:%s] [%s]\n", this.step, relateGoPath(this.gpgPath), this.impName, err.Error())
	}

	return
}

//get file suffix of code file
func (this *gopgProcessor) getCodeFileSuffix(section string) (r string) {
	if section == "" {
		section = this.impName
	}
	if v := this.getGpgCfg(section, rawKeyProductName, false); v != "" {
		r = v
	} else {
		if v := this.getGpgCfg(section, rawKeyKeyType, false); v != "" {
			l := strings.ToLower(v)
			if l != v {
				l = fmt.Sprintf("%s%s", l, getHash(v))
			}
			r = l
		}
		if v := this.getGpgCfg(section, rawKeyValueType, false); v != "" {
			l := strings.ToLower(v)
			if l != v {
				l = fmt.Sprintf("%s%s", l, getHash(v))
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
		// do not use * as filename if <VALUE_TYPE> is a pointer
		r = strings.Replace(r, "*", "#", -1)
	}

	return
}

func (this *gopgProcessor) reportNoReplacing(key, gpfile string) {
	fmt.Printf("[gogp error]: %s [%s] has no replacing. [%s:%s %s]\n", this.step, key, relateGoPath(this.gpgPath), this.impName, gpfile)
}

//if has set key GOGP_Name, use it, else use section name
func (this *gopgProcessor) getGpName() (r string) {
	if name := this.getGpgCfg(this.impName, rawKeySrcPathName, true); name != "" {
		n := filepath.Base(name)
		idx := 0
		if idx = strings.Index(n, "."); idx < 0 { //split by first '.'
			idx = len(n)
		}
		r = n[:idx]
	} else {
		r = "missing"
		fmt.Printf("[gogp error]: missing %s in %s:%s\n", rawKeySrcPathName, relateGoPath(this.gpgPath), this.impName)
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
func (this *gopgProcessor) isValidSection(section string, step gogpProcessStep) (ok bool) {
	if !strings.HasPrefix(section, txtSectionIgnore) { //not an ignore section
		if checkReverse := strings.HasPrefix(section, txtSectionReverse); checkReverse == step.IsReverse() { //if a proper section
			if !this.checkGpgCfg(section, rawKeyIgnore) { //if has ignore key
				ok = true
			}
		}
	}
	return
}

func (this *gopgProcessor) hasTask(step gogpProcessStep) bool {
	for _, imp := range this.gpgContent.Sections() {
		if this.isValidSection(imp, step) {
			return true
		}
	}
	return false
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
	//this.replaces.sectionName = section
	//this.replaces.gpgPath = this.gpgPath
	//this.replaces.gpPath = gpPath
	//println("buildMatches", section, gpPath, reverse, second, pmatch.sectionName, pmatch.gpgPath, pmatch.gpPath)
	if replaceList := this.gpgContent.Keys(section); replaceList != nil {
		//make replace map
		for _, key := range replaceList {
			replace := this.getGpgCfg(section, key, false)
			match := fmt.Sprintf(txtReplaceKeyFmt, key)
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
	if val == "" {
		if match, ok := this.maps.getMatch(key); ok {
			val = match
			return
		}
		if warnEmpty {
			fmt.Printf("[gogp warn]: [%s:%s] maybe lost key [%s]\n", relateGoPath(this.gpgPath), section, key)
		}
	}
	return
}

func (this *gopgProcessor) remove(file string) {
	fmt.Printf(">>[gogp]: [%s] removed.\n", relateGoPath(file))
	os.Remove(file)
}

func (this *gopgProcessor) getFakeSrcFilePath(pathWithName string) string {
	return fmt.Sprintf("%s.%s%s", pathWithName, gpCodeFileSuffix, codeExt)
}

func (this *gopgProcessor) getProductFilePath(gpgDir, gpName, codeFileSuffix string) string {
	return fmt.Sprintf("%s/%s.%s_%s%s", gpgDir, gpName, gpCodeFileSuffix, codeFileSuffix, codeExt)

}

func (this *gopgProcessor) getGpFullPath(gp string) string {
	gpPath := ""
	gpgDir := filepath.Dir(this.gpgPath)
	if "" == gp {
		gp = this.getGpgCfg(this.impName, rawKeySrcPathName, false) //read gp file from another path or name
	}
	if gp != "" { //read gp file from another path or name
		if !strings.HasPrefix(gp, gpExt) {
			gp += gpExt
		}
		if p, _ := filepath.Split(gp); p == "" || '.' == gp[0] { //if only config gp name, or lead with ".", use gpg dir
			gpPath = filepath.Join(gpgDir, gp)
		} else {
			gpPath = filepath.Join(goPath, gp)
		}
	} else {
		fmt.Printf("[gogp error]: missing [%s] in [%s:%s]\n", rawKeySrcPathName, relateGoPath(this.gpgPath), this.impName)
	}
	return gpPath
}

func (this *gopgProcessor) loadGpFile(file string) (err error) {
	this.gpContent = ""
	this.gpPath = file
	if this.gpContent, err = this.rawLoadFile(file); err == nil {
		//ignore text format like "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
		this.gpContent = gogpExpIgnore.ReplaceAllString(this.gpContent, "")
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
	tool := filepath.ToSlash(filepath.Dir(thisFilePath))
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
		time.Now().Format("Mon Jan 02 2006 15:04 MST"),
		relateGoPath(srcFile),
		relateGoPath(gpgFile),
		section,
		tool,
		copyRightCode,
	)
	return
}
