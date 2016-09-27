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
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp/ini"
)

type replaceCase struct {
	src, dst string
}

//object to process gpg file
type gpgProcessor struct {
	gpgPath         string            //gpg file path
	replaceMap      map[string]string //cases that need replacing
	nProcessingLine int               //line that is processing
	nNoRepMath      int               //number of math that has no replace string
}

const (
	gGpExt         = ".gp"
	gGpgExt        = ".gpg"
	gCodeExt       = ".go"
	gGpFileSuffix  = "gpg"
	gReplaceKeyFmt = "<%s>"

	//generic-programming flag <XXX>
	gReplaceExpTxt = `\<[[:alpha:]][[:word:]]{0,}\>`

	gkeyGpFilePath = "<GpFilePath>" //read gp file from another path
	gThisFilePath  = "github.com/vipally/gogp/gpg.go"

	gLibVersion = "2.1.0"
)

var (
	gReplaceExp    = regexp.MustCompile(gReplaceExpTxt)
	g_map_rep      = make(map[string]string)
	g_match_no_rep = false
	g_proc_line    = 0
	goPath         = "" //GoPath

	copyRightCode = "//    " + strings.Replace(cpright.CopyRight(), "\n", "\n//", strings.Count(cpright.CopyRight(), "\n")-1)
)

func init() {
	cmdline.Version(gLibVersion)
	copyRightCode = cmdline.ReplaceTags(copyRightCode)

	//get GoPath
	if _, __file, _, __ok := runtime.Caller(0); __ok { //0 means init func itself
		__file = filepath.ToSlash(__file)
		goPath = strings.TrimSuffix(__file, gThisFilePath)
	}

	Work(workPath()) //auto work at working path
}

func Version() string {
	return gLibVersion
}

func relateGoPath(full string) string {
	return strings.TrimPrefix(full, goPath)
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

//main func of gogp
func Work(dir string) (nGpg, nGp int, err error) {
	dir = formatPath(dir)

	files, e := collect_sub_files(dir, gGpgExt)
	if e != nil {
		err = e
		panic(err)
	}
	if nGpg = len(files); nGpg > 0 {
		fmt.Printf("[gogp]Working at:[%s]\n", relateGoPath(dir))
	}

	for _, v := range files {
		name := file_base(v)
		path_with_name := path.Join(dir, name)
		n, e := gen_gp_code_by_gpg(path_with_name)
		if e != nil {
			err = e
		}
		nGp += n
	}
	return
}

func gen_gp_code_by_gpg(path_with_name string) (nGen int, err error) {
	fmt.Printf(">[gogp]Processing:%s\n", relateGoPath(path_with_name))
	gpg_file := path_with_name + gGpgExt
	if ini, err := ini.New(gpg_file); err == nil {
		gpg_imps := ini.Sections()
		for _, gpg_imp := range gpg_imps {
			gp_reg_srcs := ini.Keys(gpg_imp)
			g_map_rep = make(map[string]string) //clear map
			for _, gp_reg_src := range gp_reg_srcs {
				replace := ini.GetString(gpg_imp, gp_reg_src, "")
				//				if replace == "" {
				//					fmt.Println(">>>>[gogp][Warn:]", relateGoPath(gpg_file), gpg_imp, gp_reg_src, "has no replace string")
				//				}
				match := fmt.Sprintf(gReplaceKeyFmt, gp_reg_src)
				g_map_rep[match] = replace
			}
			if err = gen_gp_code_by_gp(path_with_name, gpg_imp); err == nil {
				nGen++
			} else {
				panic(err)
			}
		}
	}
	return
}

func gen_gp_code_by_gp(path_with_name string, imp_name string) (err error) {
	var fin, fout *os.File
	var gpFilePath = path_with_name
	//fmt.Println("gen_gp_code_by_gp", relatePath(path_with_name), imp_name)
	if gp, ok := g_map_rep[gkeyGpFilePath]; ok { //read gp file from another path
		gpFilePath = formatPath(goPath + gp)
	}
	gp_file := gpFilePath + gGpExt
	if fin, err = os.Open(gp_file); err != nil {
		return
	}
	defer fin.Close()

	code_file := get_code_file(path_with_name, imp_name)

	if fout, err = os.OpenFile(code_file,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
		return
	}
	defer fout.Close()

	rd := bufio.NewReader(fin)
	wt := bufio.NewWriter(fout)
	if err = write_header(wt, path_with_name+gGpgExt, gp_file, imp_name); err != nil {
		return
	}
	g_proc_line = 0
	g_match_no_rep = false
	for {
		line, e := rd.ReadString('\n')
		if line != "" {
			g_proc_line++
			reped_line, _ := gen_gp_code(line)
			wt.WriteString(reped_line)
		}
		if e != nil {
			break
		}
	}
	if err = wt.Flush(); err != nil {
		return
	}
	if g_match_no_rep {
		s := fmt.Sprintf("error:[%s].[%s] not every gp have been replaced\n", relateGoPath(path_with_name), imp_name)
		fmt.Println(s)
		err = fmt.Errorf(s)
	}
	fmt.Printf(">>[gogp][%s] finish\n", relateGoPath(code_file))
	return
}

func gen_gp_code(src string) (r string, err error) {
	//	if strings.HasPrefix(src, "//") { //never replace comment line
	//		return src, nil
	//	}
	r = gReplaceExp.ReplaceAllStringFunc(src, match_replace)
	return
}

func write_header(wt *bufio.Writer, gpg_file, gp_file, imp_name string) (err error) {
	s := fmt.Sprintf(`// This file was auto-generated by [gogp] tool
// Last update at: [%s]
// Generate from:
//     [%s]
//     [%s] [%s]
// !!!!!!!!!NEVER MODIFY IT MANUALLY!!!!!!!!!

`,
		time.Now().Format("Mon Jan 02 2006 15:04:05"),
		relateGoPath(gp_file), relateGoPath(gpg_file), imp_name)
	wt.WriteString(s)
	wt.WriteString(copyRightCode)
	wt.WriteString("\n\n")
	return
}

func get_code_file(path_with_name, imp_name string) (r string) {
	r = fmt.Sprintf("%s_%s_%s%s",
		path_with_name, gGpFileSuffix, imp_name, gCodeExt)
	return
}

func match_replace(src string) (rep string) {
	if v, ok := g_map_rep[src]; ok {
		rep = v
	} else {
		fmt.Printf("error: at line %d, %s has no replacing\n", g_proc_line, src)
		rep = src
		g_match_no_rep = true
	}
	return
}

func collect_sub_files(_dir string,
	ext string) (subfiles []string, err error) {
	f, err := os.Open(_dir)
	if err != nil {
		return
	}
	defer f.Close()

	dirs, err := f.Readdir(0)
	if err != nil {
		return
	}
	//subfiles = make([]string, 0)
	for _, v := range dirs {
		if !v.IsDir() {
			filename := v.Name()
			if ext == "" || path.Ext(filename) == ext {
				subfiles = append(subfiles, filename)
			}
		}
	}
	return
}

func file_base(file_path string) (file string) {
	_, full := path.Split(file_path)
	ext := path.Ext(file_path)
	file = strings.TrimSuffix(full, ext)
	return
}
