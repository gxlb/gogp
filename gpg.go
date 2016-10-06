//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

package gogp

import (
	"fmt"

	"os"
	"path/filepath"
	"regexp"

	"strings"

	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"

	//xhash "github.com/vipally/gx/hash"
)

const (
	gGpgExt         = ".gpg"
	gGpExt          = ".gp"
	gGpFileSuffix   = "gpg"
	gReplaceKeyFmt  = "<%s>"
	gSectionReverse = "GOGP_REVERSE" //gpg section that for gogp reverse only
	gSectionIgnore  = "GOGP_IGNORE"  //gpg section that for gogp never process

	//key that set a gp name in a reverse process, and code suffix in normal work
	grawKeyName       = "GOGP_Name"
	grawKeyGpFilePath = "GOGP_GpFilePath" //read gp file from another path
	gKeyReservePrefix = "<GOGP_"          //reserved key,do not use repalce action

	//generic-programming flag <XXX>
	gReplaceExpTxt = `\<[[:alpha:]][[:word:]]{0,}\>`

	gThisFilePath = "github.com/vipally/gogp/gpg.go"

	gLibVersion = "2.9.0"
)

var (
	gReplaceExp = regexp.MustCompile(gReplaceExpTxt)

	//ignore text format like "//GOGP_IGNORE_BEGIN ... //GOGP_IGNORE_END"
	gGogpIgnoreExp = regexp.MustCompile("(?s)\\s*//GOGP_IGNORE_BEGIN.*?//GOGP_IGNORE_END.*?\\n\\s*")

	gGoPath        = "" //GoPath
	gCopyRightCode = ""
	gCodeExt       = ".go"
	gForceUpdate   = false //force update all products
	gSilence       = true  //work silencely
)

func init() {
	cmdline.Version(gLibVersion)
	gCopyRightCode = cmdline.FormatLineHead(cpright.CopyRight(), "// ")
	gCopyRightCode = cmdline.ReplaceTags(gCopyRightCode)

	//get GoPath
	s := os.Getenv("GOPATH")
	if ss := strings.Split(s, ";"); ss != nil && len(ss) > 0 {
		gGoPath = formatPath(ss[0]) + "/src/"
	}
}

//set silence work mode flag.
func Silence(enable bool) (old bool) {
	old, gSilence = gSilence, enable
	return
}

//set force update product flag.
func ForceUpdate(enable bool) (old bool) {
	old, gForceUpdate = gForceUpdate, enable
	return
}

//set extension of code file, ".go" is default
func CodeExtName(n string) (old string) {
	old = gCodeExt
	if n != "" && gCodeExt != n && n != gGpExt && n != gGpgExt {
		gCodeExt = n
	}
	return
}

//run work process on GoPath
func WorkOnGoPath() (nGpg, nCode, nSkip int, err error) {
	return Work(gGoPath)
}

//run work process on current working path
func WorkOnWorkPath() (nGpg, nCode, nSkip int, err error) {
	return Work(workPath())
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
		if !gSilence && len(list) > 0 {
			fmt.Printf("[gogp]Working at:[%s]\n", relateGoPath(dir))
		}
		for _, gpg := range list { //reverse work
			var p gopgProcessor
			if err = p.procGpg(gpg, true); err != nil {
				return
			}
			nCode += p.nCodeFile
			nSkip += p.nSkipCodeFile
		}
		for _, gpg := range list { //normal work
			nGpg++
			var p gopgProcessor
			if err = p.procGpg(gpg, false); err != nil {
				return
			}
			nCode += p.nCodeFile
			nSkip += p.nSkipCodeFile
		}
	}

	if true || !gSilence { //always show this message
		fmt.Printf("[gogp][%s] %d/%d product(s) updated from %d gpg file(s).\n", relateGoPath(dir), nCode, nCode+nSkip, nGpg)
	}

	return
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

//get version of this gogp lib
func Version() string {
	return gLibVersion
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
