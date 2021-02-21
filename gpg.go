//    CopyRight @Ally Dale 2016
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

package gogp

import (
	"fmt"
	"go/format"
	"hash/crc32"
	"time"

	"os"
	"path/filepath"
	"regexp"

	"strings"

	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
)

const (
	gGpgExt           = ".gpg"
	gGpExt            = ".gp"
	gGpCodeFileSuffix = "gp"
	gReplaceKeyFmt    = "<%s>"
	gSectionReverse   = "GOGP_REVERSE" //gpg section prefix that for gogp reverse only
	gSectionIgnore    = "GOGP_IGNORE"  //gpg section prefix that for gogp never process

	gKeyReservePrefix  = "<GOGP_"            //reserved key, who will not use repalce action
	grawKeyIgnore      = "GOGP_Ignore"       //ignore this section
	grawKeyProductName = "GOGP_CodeFileName" //code file name part
	grawKeySrcPathName = "GOGP_GpFilePath"   //gp file path and name
	grawKeyDontSave    = "GOGP_DontSave"     //do not save
	grawKeyKeyType     = "KEY_TYPE"          //key_type
	grawKeyValueType   = "VALUE_TYPE"        //value_type

	grawStringNotComment = "//#GOGP_COMMENT"

	//generic-programming flag <XXX>
	gsExpTxtReplace = `(?P<P>.?)(?P<W>\<[[:alpha:]][[:word:]]*\>)(?P<S>.?)`

	// ignore all text format:
	// //#GOGP_IGNORE_BEGIN <content> //#GOGP_IGNORE_END
	gsExpTxtIgnore = "(?sm:\\s*//#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\\r|\\n]*)"
	gsExpTxtGPOnly = "(?sm:\\s*//#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\\r|\\n]*)"

	// select by condition <cd> defines in gpg file:
	// //#GOGP_IFDEF <cd> <true_content> //#GOGP_ELSE <false_content> //#GOGP_ENDIF
	gsExpTxtChoice = "(?sm:\\s*//#GOGP_IFDEF[ |\\t]+(?P<CONDK>[[:word:]]+)(?:[ |\\t]*?//.*?$)?[\\r|\\n]*(?P<T>.*?)[\\r|\\n]*(?:[ |\\t]*?(?://)??#GOGP_ELSE(?:[ |\\t]*?//.*?$)?[\\r|\\n]*(?P<F>.*?)[\\r|\\n]*)?[ |\\t]*?(?://)??#GOGP_ENDIF.*?$[\\r|\\n]*)"

	// require another gp file:
	// //#GOGP_REQUIRE(<gpPath> [, <gpgSection>])
	gsExpTxtRequire       = "(?sm:\\s*(?P<REQ>^[ |\\t]*(?://)?#GOGP_REQUIRE\\((?P<REQP>[^\\n\\r,]*?)(?:[ |\\t]*?,[ |\\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\\((?P<REQGPG>[[:word:]]+)\\)))??(?:[ |\\t]*?\\))).*?$[\\r|\\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\\([^\\n\\r,]*?\\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\\([^\\n\\r,]*?\\))?[\\r|\\n]*)"
	gsExpTxtEmptyLine     = "(?sm:(?P<EMPTY_LINE>[\\r|\\n]{3,}))"
	gsExpTxtTrimEmptyLine = "(?s:^[\\r|\\n]*(?P<CONTENT>.*?)[\\r|\\n]*$)"

	// get gpg config string:
	// #GOGP_GPGCFG(<cfgName>)
	gsExpTxtGetGpgCfg = "(?-sm:(?://)?#GOGP_GPGCFG\\((?P<GPGCFG>[[:word:]]+)\\))"

	// #GOGP_REPLACE(<src>,<dst>)
	gsExpTxtReplaceKey = "(?-sm:(?://)?#GOGP_REPLACE\\((?P<REPSRC>\\S+)[ |\\t]*,[ |\\t]*?(?P<REPDST>\\S+)\\))"

	//remove "*" from value type such as "*string -> string"
	// #GOGP_RAWNAME(<strValueType>)
	//gsExpTxtRawName = "(?-sm:(?://)?#GOGP_RAWNAME\\((?P<RAWNAME>\\S+)\\))"

	// only generate <content> once from a gp file:
	// //#GOGP_ONCE <content> //#GOGP_END_ONCE
	gsExpTxtOnce = "(?sm:\\s*//#GOGP_ONCE(?:[ |\\t]*?//.*?$)?[\\r|\\n]*(?P<ONCE>.*?)[\\r|\\n]*[ |\\t]*?(?://)??#GOGP_END_ONCE.*?$[\\r|\\n]*)"

	gsExpTxtFileBegin = "(?sm:\\s*(?P<FILEB>//#GOGP_FILE_BEGIN(?:[ |\\t]+(?P<OPEN>[[:word:]]+))?).*?$[\\r|\\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_begin.*?(?://)?#GOGP_IGNORE_END ///gogp_file_begin.*?$)?[\\r|\\n]*)"
	gsExpTxtFileEnd   = "(?sm:\\s*(?P<FILEE>//#GOGP_FILE_END).*?$[\\r|\\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_end.*?(?://)?#GOGP_IGNORE_END ///gogp_file_end.*?$)?[\\r|\\n]*)"

	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	gsTxtRequireResultFmt   = "//#GOGP_IGNORE_BEGIN ///require begin from(%s)\n%s\n//#GOGP_IGNORE_END ///require end from(%s)"
	gsTxtRequireAtResultFmt = "///require begin from(%s)\n%s\n///require end from(%s)"
	gsTxtGogpIgnoreFmt      = "//#GOGP_IGNORE_BEGIN%s%s//#GOGP_IGNORE_END%s"

	gThisFilePath = "github.com/gxlb/gogp/gpg.go"
	gLibVersion   = "v3.0.0.final"
)

var (
	gGogpExpReplace          = regexp.MustCompile(gsExpTxtReplace)
	gGogpExpPretreatAll      = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", gsExpTxtIgnore, gsExpTxtRequire, gsExpTxtGetGpgCfg, gsExpTxtOnce, gsExpTxtReplaceKey))
	gGogpExpIgnore           = regexp.MustCompile(gsExpTxtIgnore)
	gGogpExpCodeIgnore       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", gsExpTxtIgnore, gsExpTxtGPOnly, gsExpTxtChoice))
	gGogpExpEmptyLine        = regexp.MustCompile(gsExpTxtEmptyLine)
	gGogpExpTrimEmptyLine    = regexp.MustCompile(gsExpTxtTrimEmptyLine)
	gGogpExpRequire          = regexp.MustCompile(gsExpTxtRequire)
	gGogpExpRequireAll       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", gsExpTxtRequire, gsExpTxtFileBegin, gsExpTxtFileEnd))
	gGogpExpReverseIgnoreAll = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", gsExpTxtFileBegin, gsExpTxtFileEnd, gsExpTxtIgnore))
	//gGogpExpRawName          = regexp.MustCompile(gsExpTxtRawName)
	//gGogpExpChoice = regexp.MustCompile(gsExpTxtChoice)

	gGoPath             = "" //GoPath
	gCopyRightCode      = ""
	gCodeExt            = ".go"
	gForceUpdate        = false //force update all products
	gSilence            = true  //work silencely
	gRemoveProductsOnly = false //remove products only

	gsTxtFileBeginContent = `//
/*   //This line can be uncommented to disable all this file, and it doesn't effect to the .gp file
//	 //If test or change .gp file required, comment it to modify and compile as normal go file
//
// This is a fake go code file
// It is used to generate .gp file by gogp tool
// Real go code file will be generated from .gp file
//
`
	gsTxtFileBeginContentOpen = strings.Replace(gsTxtFileBeginContent, "/*", "///*", 1)
	gsTxtFileEndContent       = "//*/\n"

	gOnceMap       map[string]bool //record once processed files
	gSavedCodeFile map[string]bool //record saved code files
	gDebug         = false         //debug switch
)

func gGetTxtFileBeginContent(open bool) (r string) {
	if open {
		r = gsTxtFileBeginContentOpen
	} else {
		r = gsTxtFileBeginContent
	}
	return
}

type gogp_proc_step int

func (me gogp_proc_step) IsReverse() bool {
	return me >= gogp_step_REQUIRE && me <= gogp_step_REVERSE
}

func (me gogp_proc_step) String() (s string) {
	switch me {
	case gogp_step_REQUIRE:
		s = "Step=[1RequireReplace]"
	case gogp_step_REVERSE:
		s = "Step=[2ReverseWork]"
	case gogp_step_PRODUCE:
		s = "Step=[3NormalProduce]"
	default:
		s = "Step=Unknown"
	}
	return
}

const (
	gogp_step_REQUIRE gogp_proc_step = iota + 1 //require replace in fake go file
	gogp_step_REVERSE                           //gen gp file from fake go file
	gogp_step_PRODUCE                           //gen go file from gp file
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
	gOnceMap = make(map[string]bool)
	gSavedCodeFile = make(map[string]bool)
}

// enable/disable work mode RemoveProductsOnly.
func RemoveProductsOnly(enable bool) (old bool) {
	old, gRemoveProductsOnly = gRemoveProductsOnly, enable
	return
}

//set debug mode flag.
func Debug(enable bool) (old bool) {
	old, gDebug = gDebug, enable
	return
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

	start := time.Now()

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

		steps := []gogp_proc_step{gogp_step_REVERSE, gogp_step_REQUIRE, gogp_step_REVERSE, gogp_step_PRODUCE} //reverse work first
		if gRemoveProductsOnly {
			steps = []gogp_proc_step{gogp_step_PRODUCE, gogp_step_REQUIRE, gogp_step_REVERSE} //normal work first
		}
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

	if true || !gSilence { //always show this message
		cost := time.Now().Sub(start)
		fmt.Printf("[gogp][%s] %d/%d product(s) updated from %d gpg file(s) in %s.\n", relateGoPath(dir), nCode, nCode+nSkip, nGpg, cost)
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

func get_hash(s string) string {
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

//remove "*" from src
func gGetRawName(src string) (r string) {
	r = strings.Replace(src, "*", "", -1)
	return
}
