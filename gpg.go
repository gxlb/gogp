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
	"go/format"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
)

const (
	gpgExt           = ".gpg"
	gpExt            = ".gp"
	gpCodeFileSuffix = "gp"

	thisFilePath = "github.com/gxlb/gogp/gpg.go"
	libVersion   = "v3.1.1"
)

var (
	goPath                = "" //GoPath
	copyRightCode         = ""
	codeExt               = ".go"
	optForceUpdate        = false //force update all products
	optSilence            = true  //work silencely
	optRemoveProductsOnly = false //remove products only

	onceMap       map[string]bool //record once processed files
	savedCodeFile map[string]bool //record saved code files
	debug         = false         //debug switch
)

func init() {
	cmdline.Version(libVersion)
	copyRightCode = cmdline.FormatLineHead(cpright.CopyRight(), "// ")
	copyRightCode = cmdline.ReplaceTags(copyRightCode)

	//get GoPath
	s := os.Getenv("GOPATH")
	if ss := strings.Split(s, ";"); ss != nil && len(ss) > 0 {
		goPath = formatPath(ss[0]) + "/src/"
	}
	onceMap = make(map[string]bool)
	savedCodeFile = make(map[string]bool)
}

// enable/disable work mode RemoveProductsOnly.
func RemoveProductsOnly(enable bool) (old bool) {
	old, optRemoveProductsOnly = optRemoveProductsOnly, enable
	return
}

//set debug mode flag.
func Debug(enable bool) (old bool) {
	old, debug = debug, enable
	return
}

//set silence work mode flag.
func Silence(enable bool) (old bool) {
	old, optSilence = optSilence, enable
	return
}

//set force update product flag.
func ForceUpdate(enable bool) (old bool) {
	old, optForceUpdate = optForceUpdate, enable
	return
}

//set extension of code file, ".go" is default
func CodeExtName(n string) (old string) {
	old = codeExt
	if n != "" && codeExt != n && n != gpExt && n != gpgExt {
		codeExt = n
	}
	return
}

//run work process on GoPath
func WorkOnGoPath() (nGpg, nCode, nSkip int, err error) {
	return Work(goPath)
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

	start := time.Now()

	if dir == "" || strings.ToLower(dir) == "gopath" { //if not set a dir,use GoPath
		dir = goPath
	} else if dir == "." || strings.ToLower(dir) == "workpath" {
		dir = workPath()
	}
	dir = formatPath(dir)
	//println(dir)

	var list []string
	if list, err = deepCollectSubFiles(dir, gpgExt); err == nil {
		//fmt.Println("list", list)
		if !optSilence && len(list) > 0 {
			fmt.Printf("[gogp]Working at:[%s]\n", relateGoPath(dir))
		}

		steps := getProcessingSteps(optRemoveProductsOnly)
		nGpg = len(list)
		for _, step := range steps {
			for _, gpg := range list {
				var p gopgProcessor
				if err = p.procGpg(gpg, step); err != nil {
					return
				}
				nCode += p.nCodeFile
				nSkip += p.nSkipCodeFile
			}
		}
	}

	if true || !optSilence { //always show this message
		cost := time.Now().Sub(start)
		fmt.Printf("[gogp][%s] %d/%d product(s) updated from %d gpg file(s) in %s.\n", relateGoPath(dir), nCode, nCode+nSkip, nGpg, cost)
	}

	return
}

//get version of this gogp lib
func Version() string {
	return libVersion
}

func getTxtFileBeginContent(open bool) (r string) {
	if open {
		r = txtFileBeginContentOpen
	} else {
		r = txtFileBeginContent
	}
	return
}

func relateGoPath(full string) string {
	fp := filepath.ToSlash(filepath.Clean(full))
	fg := formatPath(goPath)
	//println("relateGoPath", fp, fg)
	if !filepath.HasPrefix(fp, fg) {
		return fp
	}
	return strings.TrimPrefix(fp, goPath)
}
func expadGoPath(path string) (r string) {
	r = path
	if filepath.VolumeName(path) == "" {
		r = filepath.Join(goPath, path)
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

func getHash(s string) string {
	h := crc32.NewIEEE()
	h.Write([]byte(s))
	r := fmt.Sprintf("%04x", (h.Sum32() & 0xFFFF))
	return r
}

func goFmt(s, file string) (r string) {
	if b, e := format.Source([]byte(s)); e != nil {
		fmt.Println(relateGoPath(file), e)
		r = s
	} else {
		r = string(b)
	}
	return
}
