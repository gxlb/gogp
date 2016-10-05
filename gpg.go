//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

/*
package gogp is a generic-programming solution for golang or any other languages.

Detail desctription:
    1. .gpg files
      An ini file in fact.It's used to define generic parameters's replacing relation.
      Corresponding .gp file may with the same path and name.
      But we can redirect it by key "GOGP_GpFilePath".
      Section "GOGP_REVERSE" is defined for ReverseWork to generate .gp file from .go file.
      So normal work mode will not generate go code file for this section.

    2. .gp files
      A go-like file, but exists some <xxx> format keys,
      that need to be replaced with which defined in .gpg file.

    3. .go files
      gogp tool auto-generated .go files are exactly normal go code files.
      But never modify it manually, you can see this warning at the first line in every file.
      Auto work on GoPath is recmmended.
      gogp tool will deep travel the path to find all .gpg files for processing.
      If the generated go code file's body has no changes, this file will not be updated.
      So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
      So any manually modification will be restored by this tool.
      Take care of that.

usage of gogp tool:
  1. (Recommend)use cmdline(cmd/gogp):

    Tool gogp is a generic-programming solution for golang or any other languages.
    Usage:
        gogp [-e=<codeExt>] [-r=<reverseWork>] <filePath>
    -e=<codeExt>  string
        Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.
    -r=<reverseWork>
        Reverse work, this mode is used to gen .gp file from a real-go file.
        If set this flag, the filePath flag must be a .gpg file path related to GoPath.
    <filePath>  required  string
        Path that gogp will work. GoPath and WorkPath is allowed.

    usage eg:
       gogp gopath
       gogp .
       gogp -r github.com/vipally/gogp/examples/reverse.gpg

  2. package usage:

    2.1 (Recommend)import gogp package in test file
      import (
          //"testing"
          "github.com/vipally/gogp"
      )
      func init() {
          gogp.WorkOnGoPath() //Recommend
          //gogp.WorkOnWorkPath()
          //gogp.Work(someDir)
          //gogp.ReverseWork("github.com/vipally/gogp/examples/reverse.gpg")
      }

    2.2 (Seldom use)import gogp package in normal package
      import (
          "github.com/vipally/gogp"
      )
      func someFunc(){
          //gogp.WorkOnGoPath()
          //gogp.ReverseWork(gpgFilePath)
          //gogp.WorkOnWorkPath()
          //gogp.Work(someDir)
      }
*/
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

	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp/ini"
)

const (
	gGpgExt          = ".gpg"
	gGpExt           = ".gp"
	gGpFileSuffix    = "gpg"
	gReplaceKeyFmt   = "<%s>"
	gSectionReversse = "GOGP_REVERSE" //gpg section that for gogp reverse only

	//generic-programming flag <XXX>
	gReplaceExpTxt = `\<[[:alpha:]][[:word:]]{0,}\>`

	gkeyGpFilePath = "<GOGP_GpFilePath>" //read gp file from another path
	gThisFilePath  = "github.com/vipally/gogp/gpg.go"

	gLibVersion = "2.9.0"
)

var (
	gReplaceExp    = regexp.MustCompile(gReplaceExpTxt)
	gGogpIgnoreExp = regexp.MustCompile("(?s)\\s+//GOGP_IGNORE_BEGIN.*?//GOGP_IGNORE_END.*?\\n\\s*") //igonre text in code file

	g_map_rep = make(map[string]string)
	gGoPath   = "" //GoPath

	gCopyRightCode = "//    " + strings.Replace(cpright.CopyRight(), "\n", "\n//", strings.Count(cpright.CopyRight(), "\n")-1)
	gCodeExt       = ".go"
)

func init() {
	cmdline.Version(gLibVersion)
	gCopyRightCode = cmdline.ReplaceTags(gCopyRightCode)

	//get GoPath
	s := os.Getenv("GOPATH")
	if ss := strings.Split(s, ";"); ss != nil && len(ss) > 0 {
		gGoPath = formatPath(ss[0]) + "/src/"
	}
}

type replaceCase struct {
	key, value string
}

type replaceList struct {
	list []*replaceCase
}

func (this *replaceList) push(v *replaceCase) int {
	if v.value != "" {
		this.list = append(this.list, v)
	}
	return this.Len()
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

//set extension of code file, ".go" is default
func CodeExtName(n string) string {
	if n != "" && gCodeExt != n && n != gGpExt && n != gGpgExt {
		gCodeExt = n
	}
	return gCodeExt
}

// reverse work, gen .gp file from code and .gpg file
// gpgFilePath must related from GoPath
func ReverseWork(gpgFilePath string) (err error) {
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("[gogp]Work(%s) end: gpg=%d code=%d skip=%d\n", relateGoPath(dir), nGpg, nCode, nSkip)
	}()

	var p gopgProcessor
	if err = p.reverseWork(gpgFilePath); err != nil {
		return
	}

	return
}

//run work process on current working path
func WorkOnWorkPath() (nGpg, nCode, nSkip int, err error) {
	return Work(workPath())
}

//run work process on GoPath
func WorkOnGoPath() (nGpg, nCode, nSkip int, err error) {
	return Work(gGoPath)
}

// work, gen code from gp file
func Work(dir string) (nGpg, nCode, nSkip int, err error) {
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("[gogp]Work(%s) end: gpg=%d code=%d skip=%d\n", relateGoPath(dir), nGpg, nCode, nSkip)
	}()
	if dir == "" || strings.ToLower(dir) == "gopath" { //if not set a dir,use GoPath
		dir = gGoPath
	} else if dir == "." || strings.ToLower(dir) == "workpath" {
		dir = workPath()
	}
	dir = formatPath(dir)
	var list []string
	if list, err = deepCollectSubFiles(dir, gGpgExt); err == nil {
		if len(list) > 0 {
			fmt.Printf("[gogp]Working at:[%s]\n", relateGoPath(dir))
		}
		for _, gpg := range list {
			nGpg++
			var p gopgProcessor
			if err = p.procGpg(gpg); err != nil {
				return
			}
			nCode += p.nCodeFile
			nSkip += p.nSkipCodeFile
		}
	}
	fmt.Printf("[gogp][%s] %d/%d updated from %d gpg file(s)\n", relateGoPath(dir), nCode, nCode+nSkip, nGpg)
	return
}

//get version of this gogp lib
func Version() string {
	return gLibVersion
}

//object to process gpg file
type gopgProcessor struct {
	gpgPath    string            //gpg file path
	gpPath     string            //gp file path
	codePath   string            //code file path
	replaceMap map[string]string //cases that need replacing

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
		err = fmt.Errorf("[%s] must be %s file at reverse mode", relateGoPath(gpgFilePath), gGpgExt)
		return
	}
	gpgFullPath := formatPath(filepath.Join(gGoPath, gpgFilePath)) //make full path
	this.impName = gSectionReversse
	pathWithName := strings.TrimSuffix(gpgFullPath, gGpgExt)
	gpFilePath := pathWithName + gGpExt
	codeFilePath := pathWithName + gCodeExt
	fmt.Printf("[gogp]ReverseWork:[%s]\n", relateGoPath(gpgFullPath))
	if err = this.loadCodeFile(codeFilePath); err != nil { //load code file
		return
	}

	//ignore text format like "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END"
	this.codeContent = gGogpIgnoreExp.ReplaceAllString(this.codeContent, "\n\n")

	if err = this.loadGpgFile(gpgFullPath); err == nil {
		if keys := this.gpgContent.Keys(gSectionReversse); keys != nil {
			var sortKey replaceList
			this.replaceMap = make(map[string]string) //clear map
			for _, k := range keys {
				v := this.gpgContent.GetString(gSectionReversse, k, "")
				if v != "" {
					sortKey.push(&replaceCase{key: k, value: v})
					this.replaceMap[v] = fmt.Sprintf(gReplaceKeyFmt, k) //match key from value
				}
			}
			sort.Sort(&sortKey)
			exp := sortKey.expString()
			reg := regexp.MustCompile(exp)
			replacedCode := reg.ReplaceAllStringFunc(this.codeContent, func(src string) (rep string) {
				if v, ok := this.getMatch(src); ok {
					rep = v
				} else {
					fmt.Printf("error: %s has no replacing\n", src)
					rep = src
					this.nNoReplaceMathNum++
				}
				return
			})
			if this.nNoReplaceMathNum > 0 { //report error
				s := fmt.Sprintf("error:[%s].[%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), this.impName)
				fmt.Println(s)
				err = fmt.Errorf(s)
			}
			if err = this.saveGpFile(replacedCode, gpFilePath); err != nil { //save code to file
				return
			}
		} else {
			err = fmt.Errorf("[%s] must have [%s] section", relateGoPath(gpgFullPath), gSectionReversse)
		}
	}
	return
}
func (this *gopgProcessor) saveGpFile(body, gpFilePath string) (err error) {
	this.gpPath = gpFilePath
	var fout *os.File
	if fout, err = os.OpenFile(this.gpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
		return
	}
	defer fout.Close()
	fout.WriteString(body)

	this.nCodeFile++
	fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.gpPath))
	return
}

func (this *gopgProcessor) procGpg(file string) (err error) {
	fmt.Printf(">[gogp]Processing:[%s]\n", relateGoPath(file))
	this.gpContent = "" //clear gp content
	if err = this.loadGpgFile(file); err == nil {
		for _, imp := range this.gpgContent.Sections() {
			if err = this.genCode(imp); err != nil {
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

func (this *gopgProcessor) genCode(impName string) (err error) {
	if impName == gSectionReversse { //reverse only section, ignore it
		return
	}
	this.impName = impName
	this.nNoReplaceMathNum = 0
	this.replaceMap = make(map[string]string) //clear map
	if replaceList := this.gpgContent.Keys(impName); replaceList != nil {
		//make replace map
		for _, key := range replaceList {
			replace := this.gpgContent.GetString(impName, key, "")
			match := fmt.Sprintf(gReplaceKeyFmt, key)
			this.replaceMap[match] = replace
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
				fmt.Printf("error: %s has no replacing\n", src)
				rep = src
				this.nNoReplaceMathNum++
			}
			return
		})
		if this.nNoReplaceMathNum > 0 { //report error
			s := fmt.Sprintf("error:[%s].[%s] not every gp have been replaced\n", relateGoPath(this.gpgPath), impName)
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
	match, ok = this.replaceMap[key]
	return
}

func (this *gopgProcessor) loadGpFile(file string) (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(file); err == nil {
		this.gpPath = file
		this.gpContent = string(b)
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

func (this *gopgProcessor) saveCodeFile(body string) (err error) {
	if !strings.HasSuffix(this.codeContent, body) { //body change then save it,else skip it
		var fout *os.File
		if fout, err = os.OpenFile(this.codePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
			return
		}
		defer fout.Close()
		wt := bufio.NewWriter(fout)
		s := fmt.Sprintf(`///////////////////////////////////////////////////////////////////
//
//    !!!!!!!!!!NEVER MODIFY THIS FILE MANUALLY!!!!!!!!!!
//
// This file was auto-generated by tool [%s]
// Last update at: [%s]
// Generate from:
//     [%s]
//     [%s] [%s]
//
//
`, filepath.ToSlash(filepath.Dir(gThisFilePath)), time.Now().Format("Mon Jan 02 2006 15:04:05"), relateGoPath(this.gpPath), relateGoPath(this.gpgPath), this.impName)
		wt.WriteString(s)
		wt.WriteString(gCopyRightCode)
		wt.WriteString("///////////////////////////////////////////////////////////////////\n\n")
		wt.WriteString(body)
		if err = wt.Flush(); err != nil {
			return
		}

		this.nCodeFile++
		fmt.Printf(">>[gogp][%s] ok\n", relateGoPath(this.codePath))
	} else {
		this.nSkipCodeFile++
		fmt.Printf(">>[gogp][%s] skip\n", relateGoPath(this.codePath))
	}
	return
}

func relateGoPath(full string) string {
	return strings.TrimPrefix(formatPath(full), gGoPath)
}
func expadGoPath(path string) (r string) {
	r = path
	if filepath.VolumeName(path) == "" {
		r = filepath.Join(gGoPath, path)
	}
	return
}

func formatPath(path string) string {
	return filepath.ToSlash(filepath.Clean(expadGoPath(path)))
}

func workPath() (p string) {
	if dir, err := os.Getwd(); err == nil {
		p = dir
	} else {
		panic(err)
	}
	return
}

//deep find the file path
func deepCollectSubFiles(_dir string, ext string) (subfiles []string, err error) {
	err = filepath.Walk(_dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (ext == "" || filepath.Ext(path) == ext) {
			subfiles = append(subfiles, path)
		}
		return err
	})
	return
}
