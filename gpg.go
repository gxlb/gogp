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
	gpgExt           = ".gpg"
	gpExt            = ".gp"
	gpCodeFileSuffix = "gp"
	txtReplaceKeyFmt    = "<%s>"
	txtSectionReverse   = "GOGP_REVERSE" //gpg section prefix that for gogp reverse only
	txtSectionIgnore    = "GOGP_IGNORE"  //gpg section prefix that for gogp never process

	keyReservePrefix  = "<GOGP_"            //reserved key, who will not use repalce action
	rawKeyIgnore      = "GOGP_Ignore"       //ignore this section
	rawKeyProductName = "GOGP_CodeFileName" //code file name part
	rawKeySrcPathName = "GOGP_GpFilePath"   //gp file path and name
	rawKeyDontSave    = "GOGP_DontSave"     //do not save
	rawKeyKeyType     = "KEY_TYPE"          //key_type
	rawKeyValueType   = "VALUE_TYPE"        //value_type

	rawStringNotComment = "//#GOGP_COMMENT"

	//generic-programming flag <XXX>
	expTxtReplace = `(?P<P>.?)(?P<W>\<[[:alpha:]][[:word:]]*\>)(?P<S>.?)`

	// ignore all text format:
	// //#GOGP_IGNORE_BEGIN <content> //#GOGP_IGNORE_END
	expTxtIgnore = `(?sm:\s*//#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\r|\n]*)`
	expTxtGPOnly = `(?sm:\s*//#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\r|\n]*)`

	// select by condition <cd> defines in gpg file:
	// //#GOGP_IFDEF <cd> <true_content> //#GOGP_ELSE <false_content> //#GOGP_ENDIF
	// "<key> || ! <key> || <key> == xxx || <key> != xxx"
	//gsExpTxtIf = "(?sm:\s*//#GOGP_IFDEF[ |\t]+(?P<CONDK>[[:word:]<>\|!]+)(?:[ |\t]*?//.*?$)?[\r|\n]*(?P<T>.*?)[\r|\n]*(?:[ |\t]*?(?://)??#GOGP_ELSE(?:[ |\t]*?//.*?$)?[\r|\n]*(?P<F>.*?)[\r|\n]*)?[ |\t]*?(?://)??#GOGP_ENDIF.*?$[\r|\n]?)"
	expTxtIf  = `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF[ |\t]+(?P<CONDK>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE(?:.*?$[\r|\n]?)[\r|\n]*(?P<F>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF.*?$[\r|\n]?)`
	expTxtIf2 = `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF2[ |\t]+(?P<CONDK2>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T2>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE2(?:.*?$[\r|\n]?)[\r|\n]*(?P<F2>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF2.*?$[\r|\n]?)`

	// " <key> || !<key> || <key> == xxx || <key> != xxx "
	// [<NOT>] <KEY> [<OP><VALUE>]
	expCondition = `(?sm:^[ |\t]*(?P<NOT>!)?[ |\t]*(?P<KEY>[[:word:]<>]+)[ |\t]*(?:(?P<OP>==|!=)[ |\t]*(?P<VALUE>[[:word:]]+))?[ |\t]*)`

	//#GOGP_SWITCH [<SWITCHKEY>] <CASES> #GOGP_GOGP_ENDSWITCH
	expTxtSwitch = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:#GOGP_SWITCH)(?:[ |\t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASES>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDSWITCH.*?$[\r|\n]?)`

	//#GOGP_CASE <COND> <CASE> #GOGP_ENDCASE
	//#GOGP_DEFAULT <CASE> #GOGP_ENDCASE
	expTxtCase = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:(?:#GOGP_CASE[ |\t]+(?P<COND>[[:word:]<>\|!]+))|(?:#GOGP_DEFAULT))(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASE>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDCASE.*?$[\r|\n]*)`

	// require another gp file:
	// //#GOGP_REQUIRE(<gpPath> [, <gpgSection>])
	expTxtRequire   = `(?sm:\s*(?P<REQ>^[ |\t]*(?://)?#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ |\t]*?,[ |\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ |\t]*?\))).*?$[\r|\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r|\n]*)`
	expTxtEmptyLine = `(?sm:(?P<EMPTY_LINE>[\r|\n]{3,}))`

	//must be "-sm", otherwise it with will repeat every line
	expTxtTrimEmptyLine = `(?-sm:^[\r|\n]*(?P<CONTENT>.*?)[\r|\n]*$)`

	// get gpg config string:
	// #GOGP_GPGCFG(<cfgName>)
	expTxtGetGpgCfg = `(?sm:(?://)?#GOGP_GPGCFG\((?P<GPGCFG>[[:word:]]+)\))`

	// #GOGP_REPLACE(<src>,<dst>)
	expTxtReplaceKey = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ |\t]*,[ |\t]*(?P<REPDST>\S+)\))`
	expTxtMapKey     = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ |\t]*,[ |\t]*(?P<MAPDST>\S+)\))`

	//remove "*" from value type such as "*string -> string"
	// #GOGP_RAWNAME(<strValueType>)
	//gsExpTxtRawName = "(?-sm:(?://)?#GOGP_RAWNAME\((?P<RAWNAME>\S+)\))"

	// only generate <content> once from a gp file:
	// //#GOGP_ONCE <content> //#GOGP_END_ONCE
	expTxtOnce = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)//#GOGP_ONCE(?:[ |\t]*?//.*?$)?[\r|\n]*(?P<ONCE>.*?)[\r|\n]*[ |\t]*?(?://)??#GOGP_END_ONCE.*?$[\r|\n]*)`

	expTxtFileBegin = `(?sm:\s*(?P<FILEB>//#GOGP_FILE_BEGIN(?:[ |\t]+(?P<OPEN>[[:word:]]+))?).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_begin.*?(?://)?#GOGP_IGNORE_END ///gogp_file_begin.*?$)?[\r|\n]*)`
	expTxtFileEnd   = `(?sm:\s*(?P<FILEE>//#GOGP_FILE_END).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_end.*?(?://)?#GOGP_IGNORE_END ///gogp_file_end.*?$)?[\r|\n]*)`

	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	txtRequireResultFmt   = "//#GOGP_IGNORE_BEGIN ///require begin from(%s)\n%s\n//#GOGP_IGNORE_END ///require end from(%s)"
	txtRequireAtResultFmt = "///require begin from(%s)\n%s\n///require end from(%s)"
	txtGogpIgnoreFmt      = "//#GOGP_IGNORE_BEGIN%s%s//#GOGP_IGNORE_END%s"

	thisFilePath = "github.com/gxlb/gogp/gpg.go"
	libVersion   = "v3.1.0"
)

var (
	gogpExpReplace          = regexp.MustCompile(expTxtReplace)
	gogpExpPretreatAll      = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtRequire, expTxtGetGpgCfg, expTxtOnce, expTxtReplaceKey, expTxtMapKey))
	gogpExpIgnore           = regexp.MustCompile(expTxtIgnore)
	gogpExpCodeSelector       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtGPOnly, expTxtIf, expTxtIf2, expTxtMapKey, expTxtSwitch))
	gogpExpCodeCases        = regexp.MustCompile(expTxtCase)
	gogpExpEmptyLine        = regexp.MustCompile(expTxtEmptyLine)
	gogpExpTrimEmptyLine    = regexp.MustCompile(expTxtTrimEmptyLine)
	gogpExpRequire          = regexp.MustCompile(expTxtRequire)
	gogpExpRequireAll       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtRequire, expTxtFileBegin, expTxtFileEnd))
	gogpExpReverseIgnoreAll = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtFileBegin, expTxtFileEnd, expTxtIgnore))
	gogpExpCondition        = regexp.MustCompile(expTxtRequire)
	//gGogpExpRawName          = regexp.MustCompile(gsExpTxtRawName)
	//gGogpExpChoice = regexp.MustCompile(gsExpTxtChoice)

	goPath             = "" //GoPath
	copyRightCode      = ""
	codeExt            = ".go"
	optForceUpdate        = false //force update all products
	optSilence            = true  //work silencely
	optRemoveProductsOnly = false //remove products only

	txtFileBeginContent = `//
/*   //This line can be uncommented to disable all this file, and it doesn't effect to the .gp file
//	 //If test or change .gp file required, comment it to modify and compile as normal go file
//
// This is a fake go code file
// It is used to generate .gp file by gogp tool
// Real go code file will be generated from .gp file
//
`
	txtFileBeginContentOpen = strings.Replace(txtFileBeginContent, "/*", "///*", 1)
	txtFileEndContent       = "//*/\n"

	onceMap       map[string]bool //record once processed files
	savedCodeFile map[string]bool //record saved code files
	debug         = false         //debug switch
)

func gGetTxtFileBeginContent(open bool) (r string) {
	if open {
		r = txtFileBeginContentOpen
	} else {
		r = txtFileBeginContent
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

		steps := []gogp_proc_step{gogp_step_REVERSE, gogp_step_REQUIRE, gogp_step_REVERSE, gogp_step_PRODUCE} //reverse work first
		if optRemoveProductsOnly {
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
