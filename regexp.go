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

	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	txtRequireResultFmt   = "//#GOGP_IGNORE_BEGIN ///require begin from(%s)\n%s\n//#GOGP_IGNORE_END ///require end from(%s)"
	txtRequireAtResultFmt = "///require begin from(%s)\n%s\n///require end from(%s)"
	txtGogpIgnoreFmt      = "//#GOGP_IGNORE_BEGIN%s%s//#GOGP_IGNORE_END%s"
)

var (
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
