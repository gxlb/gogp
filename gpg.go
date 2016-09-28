//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

//todo:
//reverse replace?
//walk on all sub dir?
//work on GoPath, find all .gpg file, and auto generate go code for them
//backup old code file, if only time changes never make the change

//package gogp implement a way to generate go-gp code from *.gp+*.gpg file
package gogp

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
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
	gCodeExt         = ".go"
	gGpFileSuffix    = "gpg"
	gReplaceKeyFmt   = "<%s>"
	gSectionReversse = "GOGP_REVERSE" //gpg section that for gogp reverse only

	//generic-programming flag <XXX>
	gReplaceExpTxt = `\<[[:alpha:]][[:word:]]{0,}\>`

	gkeyGpFilePath = "<GpFilePath>" //read gp file from another path
	gThisFilePath  = "github.com/vipally/gogp/gpg.go"

	gLibVersion = "2.9.0"
)

var (
	gReplaceExp = regexp.MustCompile(gReplaceExpTxt)
	g_map_rep   = make(map[string]string)
	//g_match_no_rep = false
	//g_proc_line    = 0
	gGoPath = "" //GoPath

	gCopyRightCode = "//    " + strings.Replace(cpright.CopyRight(), "\n", "\n//", strings.Count(cpright.CopyRight(), "\n")-1)
)

func init() {
	cmdline.Version(gLibVersion)
	gCopyRightCode = cmdline.ReplaceTags(gCopyRightCode)

	//get GoPath
	if _, __file, _, __ok := runtime.Caller(0); __ok { //0 means init func itself
		__file = filepath.ToSlash(__file)
		gGoPath = strings.TrimSuffix(__file, gThisFilePath)
	}

	//Work(workPath()) //auto work at working path
}

type replaceCase struct {
	key, value string
}
type replaceList struct {
	list []*replaceCase
}

func (this *replaceList) push(v *replaceCase) int {
	if v.value != "" {

	}
	this.list = append(this.list, v)
	return this.Len()
}

func (this *replaceList) Len() int {
	return len(this.list)
}

//sort by value descend
//so in regexp, with the same prefix, the longer will match first
//eg: hello|hehe|he, "he" has the lowest priority but "hello" has the highest
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

// reverse work, gen .gp file from code & .gpg file
// gpgFilePath must related from GoPath
func ReverseWork(gpgFilePath string) (err error) {
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("[gogp]Work(%s) end: gpg=%d code=%d skip=%d\n", relateGoPath(dir), nGpg, nCode, nSkip)
	}()

	var p gpgProcessor
	if err = p.reverseWork(gpgFilePath); err != nil {
		return
	}

	return
}

func (this *gpgProcessor) reverseWork(gpgFilePath string) (err error) {

	if !strings.HasSuffix(gpgFilePath, gGpgExt) { //must .gpg file
		err = fmt.Errorf("[%s] must be %s file at reverse mode", relateGoPath(gpgFilePath), gGpgExt)
		return
	}

	gpgFullPath := formatPath(filepath.Join(gGoPath, gpgFilePath)) //make full path
	this.impName = gSectionReversse
	pathWithName := strings.TrimSuffix(gpgFullPath, gGpgExt)
	gpFilePath := pathWithName + gGpExt
	codeFilePath := pathWithName + gCodeExt
	if err = this.loadCodeFile(codeFilePath); err != nil { //load code file
		return
	}

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
func (this *gpgProcessor) saveGpFile(body, gpFilePath string) (err error) {
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
	if dir == "" { //if not set a dir,use GoPath
		dir = gGoPath
	}
	dir = formatPath(dir)
	var list []string
	if list, err = deepCollectSubFiles(dir, gGpgExt); err == nil {
		if len(list) > 0 {
			fmt.Printf("[gogp]Working at:[%s]\n", relateGoPath(dir))
		}
		for _, gpg := range list {
			nGpg++
			var p gpgProcessor
			if err = p.procGpg(gpg); err != nil {
				return
			}
			nCode += p.nCodeFile
			nSkip += p.nSkipCodeFile
		}
	}
	return
}

//object to process gpg file
type gpgProcessor struct {
	gpgPath    string            //gpg file path
	gpPath     string            //gp file path
	codePath   string            //code file path
	replaceMap map[string]string //cases that need replacing
	//nProcessingLine   int               //line that is processing
	nNoReplaceMathNum int //number of math that has no replace string
	nCodeFile         int
	nSkipCodeFile     int
	gpgContent        *ini.IniFile
	gpContent         string
	codeContent       string
	impName           string
	//newCodeContent    string
}

func (this *gpgProcessor) procGpg(file string) (err error) {
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
func (this *gpgProcessor) loadGpgFile(file string) (err error) {
	file = formatPath(file)
	this.gpPath = ""
	this.gpgPath = formatPath(file)
	this.gpgContent, err = ini.New(this.gpgPath)
	return
}
func (this *gpgProcessor) genCode(impName string) (err error) {
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
func (this *gpgProcessor) getMatch(key string) (match string, ok bool) {
	match, ok = this.replaceMap[key]
	return
}
func (this *gpgProcessor) loadGpFile(file string) (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(file); err == nil {
		this.gpPath = file
		this.gpContent = string(b)
	}
	return
}
func (this *gpgProcessor) loadCodeFile(file string) (err error) {
	var b []byte
	this.codeContent = ""
	this.codePath = file
	if b, err = ioutil.ReadFile(file); err == nil {
		this.codeContent = string(b)
	}
	return
}
func (this *gpgProcessor) saveCodeFile(body string) (err error) {
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

func Version() string {
	return gLibVersion
}

func relateGoPath(full string) string {
	return strings.TrimPrefix(formatPath(full), gGoPath)
}

func formatPath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
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

////main func of gogp
//func Work(dir string) (nGpg, nGp int, err error) {
//	dir = formatPath(dir)

//	files, e := deepCollectSubFiles(dir, gGpgExt)
//	if e != nil {
//		err = e
//		panic(err)
//	}
//	if nGpg = len(files); nGpg > 0 {
//		fmt.Printf("[gogp]Working at:[%s]\n", relateGoPath(dir))
//	}

//	for _, v := range files {
//		name := file_base(v)
//		path_with_name := path.Join(dir, name)
//		n, e := gen_gp_code_by_gpg(path_with_name)
//		if e != nil {
//			err = e
//		}
//		nGp += n
//	}
//	return
//}

//func gen_gp_code_by_gpg(path_with_name string) (nGen int, err error) {
//	fmt.Printf(">[gogp]Processing:%s\n", relateGoPath(path_with_name))
//	gpg_file := path_with_name + gGpgExt
//	if ini, err := ini.New(gpg_file); err == nil {
//		gpg_imps := ini.Sections()
//		for _, gpg_imp := range gpg_imps {
//			gp_reg_srcs := ini.Keys(gpg_imp)
//			g_map_rep = make(map[string]string) //clear map
//			for _, gp_reg_src := range gp_reg_srcs {
//				replace := ini.GetString(gpg_imp, gp_reg_src, "")
//				//				if replace == "" {
//				//					fmt.Println(">>>>[gogp][Warn:]", relateGoPath(gpg_file), gpg_imp, gp_reg_src, "has no replace string")
//				//				}
//				match := fmt.Sprintf(gReplaceKeyFmt, gp_reg_src)
//				g_map_rep[match] = replace
//			}
//			if err = gen_gp_code_by_gp(path_with_name, gpg_imp); err == nil {
//				nGen++
//			} else {
//				panic(err)
//			}
//		}
//	}
//	return
//}

//func gen_gp_code_by_gp(path_with_name string, imp_name string) (err error) {
//	var fin, fout *os.File
//	var gpFilePath = path_with_name
//	//fmt.Println("gen_gp_code_by_gp", relatePath(path_with_name), imp_name)
//	if gp, ok := g_map_rep[gkeyGpFilePath]; ok { //read gp file from another path
//		gpFilePath = formatPath(gGoPath + gp)
//	}
//	gp_file := gpFilePath + gGpExt
//	if fin, err = os.Open(gp_file); err != nil {
//		return
//	}
//	defer fin.Close()

//	code_file := get_code_file(path_with_name, imp_name)

//	if fout, err = os.OpenFile(code_file,
//		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
//		return
//	}
//	defer fout.Close()

//	rd := bufio.NewReader(fin)
//	wt := bufio.NewWriter(fout)
//	if err = write_header(wt, path_with_name+gGpgExt, gp_file, imp_name); err != nil {
//		return
//	}
//	g_proc_line = 0
//	g_match_no_rep = false
//	for {
//		line, e := rd.ReadString('\n')
//		if line != "" {
//			g_proc_line++
//			reped_line, _ := gen_gp_code(line)
//			wt.WriteString(reped_line)
//		}
//		if e != nil {
//			break
//		}
//	}
//	if err = wt.Flush(); err != nil {
//		return
//	}
//	if g_match_no_rep {
//		s := fmt.Sprintf("error:[%s].[%s] not every gp have been replaced\n", relateGoPath(path_with_name), imp_name)
//		fmt.Println(s)
//		err = fmt.Errorf(s)
//	}
//	fmt.Printf(">>[gogp][%s] finish\n", relateGoPath(code_file))
//	return
//}

//func gen_gp_code(src string) (r string, err error) {
//	//	if strings.HasPrefix(src, "//") { //never replace comment line
//	//		return src, nil
//	//	}
//	r = gReplaceExp.ReplaceAllStringFunc(src, match_replace)
//	return
//}

//func write_header(wt *bufio.Writer, gpg_file, gp_file, imp_name string) (err error) {
//	s := fmt.Sprintf(`// This file was auto-generated by [gogp] tool
//// Last update at: [%s]
//// Generate from:
////     [%s]
////     [%s] [%s]
//// !!!!!!!!!NEVER MODIFY IT MANUALLY!!!!!!!!!

//`,
//		time.Now().Format("Mon Jan 02 2006 15:04:05"),
//		relateGoPath(gp_file), relateGoPath(gpg_file), imp_name)
//	wt.WriteString(s)
//	wt.WriteString(gCopyRightCode)
//	wt.WriteString("\n\n")
//	return
//}

//func get_code_file(path_with_name, imp_name string) (r string) {
//	r = fmt.Sprintf("%s_%s_%s%s",
//		path_with_name, gGpFileSuffix, imp_name, gCodeExt)
//	return
//}

//func match_replace(src string) (rep string) {
//	if v, ok := g_map_rep[src]; ok {
//		rep = v
//	} else {
//		fmt.Printf("error: at line %d, %s has no replacing\n", g_proc_line, src)
//		rep = src
//		g_match_no_rep = true
//	}
//	return
//}

//func file_base(file_path string) (file string) {
//	_, full := path.Split(file_path)
//	ext := path.Ext(file_path)
//	file = strings.TrimSuffix(full, ext)
//	return
//}
