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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// generate .gp file
func (this *gopgProcessor) procStep2Reverse() (err error) {
	pathWithName := filepath.Join(filepath.Dir(this.gpgPath), this.getGpName())
	gpFilePath := pathWithName + gpExt
	codeFilePath := this.getFakeSrcFilePath(pathWithName)
	this.codePath = codeFilePath
	this.gpPath = gpFilePath

	if err = this.loadCodeFile(this.codePath); err != nil { //load code file
		return
	}

	//ignore text format like "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	// []string{"", "FILEB", "OPEN", "FILEE", "IGNORE"}
	this.codeContent = gogpExpReverseIgnoreAll.ReplaceAllString(this.codeContent, "\n\n")

	if this.buildMatches(this.section, this.gpPath, true, false) {
		this.matches.sort()
		replacedCode, norep := this.matches.doReplacing(this.codeContent, this.gpgPath, true)
		this.nNoReplaceMathNum += norep

		replacedCode = gogpExpEmptyLine.ReplaceAllString(replacedCode, "\n\n") //avoid multi empty lines

		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("[gogp error]: [%s:%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.section)
			err = fmt.Errorf(s)
		}

		if err = this.saveGpFile(replacedCode, this.gpPath); err != nil { //save code to file
			return
		}
	} else {
		err = fmt.Errorf("[gogp error]: [%s] must have [%s] section", relateGoPath(this.gpgPath), txtSectionReverse)
	}
	return
}

func (this *gopgProcessor) saveGpFile(body, gpFilePath string) (err error) {
	this.gpPath = gpFilePath
	if optRemoveProductsOnly { //remove products only
		this.nCodeFile++
		this.remove(this.gpPath)
		return
	}
	if !optForceUpdate && this.loadGpFile(gpFilePath) == nil { //check if need update
		if this.gpContent == body { //body not change
			this.nSkipCodeFile++
			if !optSilence {
				fmt.Printf(">>[gogp][%s] skip\n", relateGoPath(this.gpPath))
			}
			return
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

`, this.fileHead(this.codePath, this.gpgPath, this.section))
	wt.WriteString(h)
	wt.WriteString(body)
	if err = wt.Flush(); err != nil {
		return
	}

	this.nCodeFile++
	if !optSilence {
		fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.gpPath))
	}
	return
}
