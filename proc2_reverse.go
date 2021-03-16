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
	this.codeContent = gogpExpReverseIgnoreAll.ReplaceAllString(this.codeContent, "\n\n")

	if this.buildMatches(this.impName, this.gpPath, true, false) {
		this.matches.sort()
		replacedCode, norep := this.matches.doReplacing(this.codeContent, this.gpgPath, true)
		this.nNoReplaceMathNum += norep

		replacedCode = gogpExpEmptyLine.ReplaceAllString(replacedCode, "\n\n") //avoid multi empty lines

		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("[gogp error]: [%s:%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.impName)
			//fmt.Println(s)
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
	if !optSilence {
		fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.gpPath))
	}
	return
}