package gogp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// gen .go file from .gp and .gpg
func (this *gopgProcessor) procStep3Produce() (err error) {
	//normal process
	gpPath := this.getGpFullPath("")
	gpgDir := filepath.Dir(this.gpgPath)

	gpName := strings.TrimSuffix(filepath.Base(gpPath), gpExt)
	codePath := this.getProductFilePath(gpgDir, gpName, this.getCodeFileSuffix(this.impName))

	this.loadCodeFile(codePath) //load code file, ignore error
	if this.gpPath != gpPath {  //load gp file if needed
		if err = this.loadGpFile(gpPath); err != nil {
			return
		}
	}

	replacedGp := ""
	if replacedGp, err = this.doGpReplace(this.gpPath, this.gpContent, this.impName, 0, false); err != nil {
		// if err = this.saveCodeFile(replacedGp); err != nil { //save code to file
		// 	return
		// }
		return
	}

	if err = this.saveCodeFile(replacedGp); err != nil { //save code to file
		return
	}

	return
}

func (this *gopgProcessor) doPredefReplace(gpPath, content, section string, nDepth int) (rep string) {
	pathIdentify := fmt.Sprintf("%s|%s", relateGoPath(gpPath), relateGoPath(filepath.Dir(this.gpgPath))) //gp file+gpg path=unique
	this.replaces.clear()

	// match "//#GOGP_IFDEF cdk ... //#GOGP_ELSE ... //#GOGP_ENDIF" case
	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END
	// "//#GOGP_REQUIRE(path [, gpgSection])"
	//"(?-sm:(?://)?#GOGP_REPLACE\\((?P<REPSRC>\\S+)[ |\\t]*,[ |\\t]*?(?P<REPDST>\\S+)\\))"
	for _content, needReplace, i := content, true, 0; needReplace && i < 3; _content, i = rep, i+1 {
		needReplace = false
		//fmt.Println("try match case", i, 3, _content)
		rep = gogpExpPretreatAll.ReplaceAllStringFunc(_content, func(src string) (_rep string) {
			//[]string{"", "IGNORE", "REQ", "REQP", "REQN", "REQGPG", "REQCONTENT", "GPGCFG", "ONCE", "REPSRC", "REPDST", "MAPSRC", "MAPDST"}
			elem := gogpExpPretreatAll.FindAllStringSubmatch(src, -1)[0]
			ignore, req, reqp, reqn, reqgpg, reqcontent, gpgcfg, once, repsrc, repdst, mapsrc, mapdst :=
				elem[1], elem[2], elem[3], elem[4], elem[5], elem[6], elem[7], elem[8], elem[9], elem[10], elem[11], elem[12]

			if reqgpg != "" && reqn == "" { //section name is config from gpg file
				reqn = this.getGpgCfg(section, reqgpg, true)
			}

			if !optSilence && i > 1 {
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
				if _, ok := onceMap[pathIdentify]; ok { //check if has processed this file
					_rep = "\n\n"
					if debug {
						fmt.Printf("[gogp debug]: %s GOGP_ONCE(%s:%s) ignore [%#v]\n", this.step, pathIdentify, section, once)
					}

				} else {
					_rep = fmt.Sprintf("\n\n%s\n\n", once)
					if debug {
						fmt.Printf("[gogp debug]: %s GOGP_ONCE(%s:%s) ok [%#v]\n", this.step, pathIdentify, section, once)
					}
				}
			case mapsrc != "":
				_rep = ""
				this.replaces.insert(mapdst, mapsrc, true)
				if debug {
					fmt.Printf("[debug] map %s %s %s map [%s] -> [%s]\n", gpPath, section, src, mapsrc, mapdst)
				}
			case repsrc != "":
				_rep = ""
				this.replaces.insert(repdst, repsrc, true)
				if debug {
					fmt.Printf("[debug]%s %s %s replace [%s] -> [%s]\n", gpPath, section, src, repsrc, repdst)
				}

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
		//fmt.Println(rep)
		//		if section == "GOGP_REVERSE_datadef" {
		//			fmt.Printf("@@i=%d _content=[%s]\ni=%d rep=[%s]\n", i, _content, i, rep)
		//		}
	}

	if this.step == gogpStepPRODUCE { //prevent gen #GOGP_ONCE code twice when gen code
		onceMap[pathIdentify] = true //record processed gp file
	}

	return
}

func (this *gopgProcessor) doGpReplace(gpPath, content, section string, nDepth int, second bool) (replacedGp string, err error) {
	_path := fmt.Sprintf("%s|%s", relateGoPath(gpPath), relateGoPath(filepath.Dir(this.gpgPath))) //gp file+gpg path=unique
	//fmt.Println("doGpReplace", this.step.String(), content)
	replacedGp = content
	this.replaces.clear()

	if this.step == gogpStepPRODUCE {
		//replacedGp = gGogpExpCodeIgnore.ReplaceAllString(replacedGp, "\n\n")
		replacedGp = this.pretreatGpCodeSelector(replacedGp, section)
		//fmt.Println("doGpReplace1", this.step.String(), replacedGp)
	}

	replacedGp = this.doPredefReplace(gpPath, replacedGp, section, nDepth)
	//fmt.Println("doGpReplace2", this.step.String(), content)

	//replaces keys that need be replacing
	if this.replaces.Len() > 0 {
		//fmt.Println(this.replaces.expString())
		noRep := 0
		replacedGp, noRep = this.replaces.doReplacing(replacedGp, this.gpgPath, true)
		if noRep > 0 {
			//for _, v := range this.replaces.list {
			//	fmt.Println(v)
			//}
			//fmt.Println("doGpReplace3", noRep, this.step.String(), gpPath, section, content)
			//os.Exit(1)
		}
		this.replaces.clear()
	}

	//gen code file content
	this.buildMatches(section, gpPath, false, second)
	replist := this.getReplist(second)
	norep := 0
	replacedGp, norep = replist.doReplacing(replacedGp, this.gpgPath, false)
	this.nNoReplaceMathNum += norep

	replacedGp = gogpExpEmptyLine.ReplaceAllString(replacedGp, "\n") //avoid multi empty lines

	//remove more empty line
	replacedGp = goFmt(replacedGp, this.gpgPath)

	if this.nNoReplaceMathNum > 0 { //report error
		s := fmt.Sprintf("[gogp error]: [%s:%s %s depth=%d] not every gp have been replaced\n", relateGoPath(this.gpgPath), relateGoPath(_path), replist.sectionName, nDepth)
		//fmt.Println(s)
		fmt.Printf("----**result is:\n%s\n----**end\n", replacedGp)
		err = fmt.Errorf(s)
	}

	//	if section == "GOGP_REVERSE_datadef" {
	//		fmt.Printf("@@doGpReplace replacedGp=[%s]\n", replacedGp)
	//	}

	return
}

func (this *gopgProcessor) saveCodeFile(body string) (err error) {
	if optRemoveProductsOnly { //remove products only
		this.nCodeFile++
		this.remove(this.codePath)
		return
	}
	if optForceUpdate || !strings.HasSuffix(this.codeContent, body) { //body change then save it,else skip it
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
		if !optSilence {
			fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.codePath))
		}
	} else {
		this.nSkipCodeFile++
		if !optSilence {
			fmt.Printf(">>[gogp][%s] skip\n", relateGoPath(this.codePath))
		}
	}
	return
}
