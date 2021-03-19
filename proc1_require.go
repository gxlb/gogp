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
	"path/filepath"
	"strings"
)

func (this *gopgProcessor) procStep1Require() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	codeFilePath := this.getFakeSrcFilePath(pathWithName)
	this.codePath = codeFilePath

	this.buildMatches(this.section, "", false, false)

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	replcaceCnt := 0

	// []string{"", "REQ", "REQP", "REQN", "REQGPG", "REQCONTENT", "FILEB", "OPEN", "FILEE"}
	replacedCode := gogpExpRequireAll.ReplaceAllStringFunc(this.codeContent, func(src string) (rep string) {
		elem := gogpExpRequireAll.FindAllStringSubmatch(src, -1)[0]
		req, reqp, reqn, reqgpg, content, fileb, open, filee := elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7], elem[8]

		reqp, reqn, reqgpg, content = reqp, reqn, reqgpg, content //avoid compile error

		var err error
		var replaced bool
		switch {
		case req != "":
			//fmt.Printf("req: %#v\n", elem[1:])
			if rep, replaced, err = this.procRequireReplacement(src, this.section, 0); err != nil {
				fmt.Println(err)
			}
		case fileb != "":
			if optRemoveProductsOnly {
				rep, replaced = fmt.Sprintf("\n\n%s\n\n", fileb), true
				break
			}
			repContent := fmt.Sprintf(txtGogpIgnoreFmt, " ///gogp_file_begin\n", getTxtFileBeginContent(open != ""), " ///gogp_file_begin\n\n")
			if rep, replaced = src, !strings.Contains(src, gogpExpTrimEmptyLine.ReplaceAllString(repContent, "$CONTENT")); replaced {
				rep = fmt.Sprintf("\n\n%s\n%s", fileb, repContent)
			}
		case filee != "":
			if optRemoveProductsOnly {
				rep, replaced = fmt.Sprintf("\n\n%s\n\n", filee), true
				break
			}
			repContent := fmt.Sprintf(txtGogpIgnoreFmt, " ///gogp_file_end\n", txtFileEndContent, " ///gogp_file_end\n\n")
			if rep, replaced = src, !strings.Contains(src, gogpExpTrimEmptyLine.ReplaceAllString(repContent, "$CONTENT")); replaced {
				rep = fmt.Sprintf("\n\n%s\n%s", filee, repContent)
			}
		}
		if replaced {
			replcaceCnt++
		}
		return
	})

	if optForceUpdate || replcaceCnt > 0 {
		replacedCode = gogpExpEmptyLine.ReplaceAllString(replacedCode, "\n\n") //avoid multi empty lines
		replacedCode = goFmt(replacedCode, this.gpPath)

		if err = this.rawSaveFile(this.codePath, replacedCode); err == nil {
			this.nCodeFile++
			if !optSilence {
				fmt.Printf(">>[gogp] %s updated for #GOGP_REQUIRE\n", relateGoPath(this.codePath))
			}
		} else {
			fmt.Println(err)
		}
	}

	return
}

//require a gp file, maybe recursive
func (this *gopgProcessor) procRequireReplacement(statement, section string, nDepth int) (rep string, replaced bool, err error) {
	//fmt.Println("statement", statement)
	rep = statement
	if nDepth >= 5 {
		panic(fmt.Sprintf("[gogp error] [%s:%s]maybe loop recursive of #GOGP_REQUIRE(...), %d", relateGoPath(this.gpgPath), section, nDepth))
	}

	elem := gogpExpRequire.FindAllStringSubmatch(statement, -1)[0] //{"", "REQ", "REQP", "REQN","REQGPG","CONTENT"}
	req, reqp, reqn, reqgpg, content := elem[1], elem[2], elem[3], elem[4], elem[5]

	if debug {
		fmt.Printf("[gogp debug] #GOGP_REQUIRE: [%s][%s][%s][%s][%s]\n", req, reqp, reqn, reqgpg, content)
	}

	if reqgpg != "" && reqn == "" { //section name is config from gpg file
		reqn = this.getGpgCfg(section, reqgpg, true)
	}

	replaceSection := reqn
	leftFmt := txtRequireResultFmt                            //left required file
	at := len(replaceSection) > 0 && replaceSection[0] == '@' //left result in this file
	if at {
		replaceSection = replaceSection[1:]
		leftFmt = txtRequireAtResultFmt
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
		//fmt.Println("gpContent", gpContent)
		replacedGp := ""
		if this.step == gogpStepPRODUCE {
			replaced = true
			if at {
				rep = "\n" + content + "\n"
			} else {
				rep = "\n\n"
			}

			if !at && !sharp && reqn != "_" && !this.checkGpgCfg(replaceSection, rawKeyDontSave) { //reqn=="_" will not generate this code file
				gpgDir := filepath.Dir(this.gpgPath)
				gpName := strings.TrimSuffix(filepath.Base(gpFullPath), gpExt)
				codePath := this.getProductFilePath(gpgDir, gpName, this.getCodeFileSuffix(replaceSection))

				if optRemoveProductsOnly { //remove products only
					this.nCodeFile++
					this.remove(codePath)
					return
				}

				if replacedGp, err = this.doGpReplace(gpFullPath, gpContent, replaceSection, nDepth, true); err != nil {
					return
				}

				if _, ok := savedCodeFile[codePath]; ok { //skip saved file
					//					if gDebug {
					//						fmt.Printf("[gogp] debug: step%d Required file [%s] skip\n", this.step, codePath)
					//					}
					return
				} else {
					savedCodeFile[codePath] = true //to prevent rewrite this file no matter it chages or not
					//					if gDebug {
					//						fmt.Printf("[gogp] debug: step%d Required file [%s] save ok\n", this.step, codePath)
					//					}
				}

				oldCode, _ := this.rawLoadFile(codePath)

				if optForceUpdate || !strings.HasSuffix(oldCode, replacedGp) { //body change then save it,else skip it
					codeContent := this.fileHead(gpFullPath, this.gpgPath, replaceSection) + "\n" + replacedGp
					codeContent = goFmt(codeContent, codePath)
					if err = this.rawSaveFile(codePath, codeContent); err == nil {
						this.nCodeFile++
						if !optSilence {
							fmt.Printf(">>[gogp] #GOGP_REQUIRE [%s:%s -> %s] ok\n", relateGoPath(gpFullPath), replaceSection, relateGoPath(codePath))
						}
					}
				} else {
					this.nSkipCodeFile++
					if !optSilence {
						fmt.Printf(">>[gogp] #GOGP_REQUIRE [%s:%s -> %s] skip\n", relateGoPath(gpFullPath), replaceSection, relateGoPath(codePath))
					}
				}
			}
		} else {
			if optRemoveProductsOnly {
				rep = fmt.Sprintf("\n\n%s\n\n", req)
				if !optSilence {
					fmt.Printf("%#v\n", rep)
				}
				replaced = true
			} else {
				if nDepth == 0 { //do not let require recursive
					//fmt.Println("000", gpContent)
					if replacedGp, err = this.doGpReplace(gpFullPath, gpContent, replaceSection, nDepth, true); err != nil {
						return
					}
					//fmt.Println("111", replacedGp)
					//					if section == "GOGP_REVERSE_datadef" {
					//						fmt.Printf("@@procRequireReplacement replacedGp=[%s]\n", replacedGp)
					//					}
					replacedGp = strings.Replace(replacedGp, "package", "//package", -1) //comment package declaration
					replacedGp = strings.Replace(replacedGp, "import", "//import", -1)
					//fmt.Println("222", replacedGp)
					//reqSave := strings.Replace(req, "//#GOGP_REQUIRE", "//##GOGP_REQUIRE", -1)
					reqResult := fmt.Sprintf(leftFmt, reqp, "$CONTENT", reqp)
					out := fmt.Sprintf("\n\n%s\n%s\n\n", req, reqResult)
					//fmt.Println("out", out)
					//fmt.Printf("replacedGp0 %#v\n", replacedGp)
					//fmt.Printf("%t %#v\n", gogpExpTrimEmptyLine.MatchString(replacedGp), gogpExpTrimEmptyLine.FindAllStringSubmatch(replacedGp, -1))
					replacedGp = gogpExpTrimEmptyLine.ReplaceAllString(replacedGp, out)
					//fmt.Println("replacedGp1", replacedGp)

					oldContent := content
					//fmt.Println("333", replacedGp)

					rep = goFmt(replacedGp, this.gpPath)

					//check if content changed
					replaced = oldContent == "" || !strings.Contains(rep, oldContent) //|| !strings.Contains(oldContent, "//#GOGP_IGNORE_BEGIN")
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
	//fmt.Println("statement", statement)
	//fmt.Println("rep", rep)
	return
}
