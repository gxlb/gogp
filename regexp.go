package gogp

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	txtReplaceKeyFmt  = "<%s>"
	txtSectionReverse = "GOGP_REVERSE" //gpg section prefix that for gogp reverse only
	txtSectionIgnore  = "GOGP_IGNORE"  //gpg section prefix that for gogp never process

	keyReservePrefix  = "<GOGP_"            //reserved key, who will not use repalce action
	rawKeyIgnore      = "GOGP_Ignore"       //ignore this section
	rawKeyProductName = "GOGP_CodeFileName" //code file name part
	rawKeySrcPathName = "GOGP_GpFilePath"   //gp file path and name
	rawKeyDontSave    = "GOGP_DontSave"     //do not save
	rawKeyKeyType     = "KEY_TYPE"          //key_type
	rawKeyValueType   = "VALUE_TYPE"        //value_type

	// // #GOGP_COMMENT
	expTxtGogpComment = `(?sm:(?P<COMMENT>/{2,}[ |\t]*#GOGP_COMMENT))`

	//generic-programming flag <XXX>
	expTxtTodoReplace = `(?P<P>.?)(?P<W>\<[[:alpha:]_][[:word:]]*\>)(?P<S>.?)`

	// ignore all text format:
	// //#GOGP_IGNORE_BEGIN <content> //#GOGP_IGNORE_END
	expTxtIgnore = `(?sm:\s*//#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\r|\n]*)`
	expTxtGPOnly = `(?sm:\s*//#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\r|\n]*)`

	// select by condition <cd> defines in gpg file:
	// //#GOGP_IFDEF <cd> <true_content> //#GOGP_ELSE <false_content> //#GOGP_ENDIF
	// "<key> || ! <key> || <key> == xxx || <key> != xxx"
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
)

var (
	gogpExpTodoReplace      = regexp.MustCompile(expTxtTodoReplace)
	gogpExpPretreatAll      = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtRequire, expTxtGetGpgCfg, expTxtOnce, expTxtReplaceKey, expTxtMapKey))
	gogpExpIgnore           = regexp.MustCompile(expTxtIgnore)
	gogpExpCodeSelector     = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtGPOnly, expTxtIf, expTxtIf2, expTxtMapKey, expTxtSwitch))
	gogpExpCases            = regexp.MustCompile(expTxtCase)
	gogpExpEmptyLine        = regexp.MustCompile(expTxtEmptyLine)
	gogpExpTrimEmptyLine    = regexp.MustCompile(expTxtTrimEmptyLine)
	gogpExpRequire          = regexp.MustCompile(expTxtRequire)
	gogpExpRequireAll       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtRequire, expTxtFileBegin, expTxtFileEnd))
	gogpExpReverseIgnoreAll = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtFileBegin, expTxtFileEnd, expTxtIgnore))
	gogpExpCondition        = regexp.MustCompile(expTxtRequire)
	gogpExpComment          = regexp.MustCompile(expTxtGogpComment)

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
)
