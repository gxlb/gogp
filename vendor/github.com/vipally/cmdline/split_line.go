// CopyRight 2016 @Ally Dale. All rights reserved.
// Author  : Ally Dale(vipally@gmail.com)
// Blog    : http://blog.csdn.net/vipally
// Site    : https://github.com/vipally

package cmdline

import (
	"regexp"
)

var ( //new line with any \t
	lineHeadExpr = regexp.MustCompile(`(?m)^\t*`)
)

func isSpace(c byte) bool {
	return (c == ' ' || c == '\t')
}

//SplitLine splits a command-line text separated with any ' ' or '\t'
func SplitLine(s string) []string {
	n := (len(s) + 1) / 2
	lenSep := 1
	start := 0
	a := make([]string, n)
	na := 0
	inString := 0
	escape := 0
	lastQuot := byte(0)
	for i := 0; i+lenSep <= len(s) && na+1 < n; i++ {
		// consider " xxx 'yyy' zzz" as a single string
		// " xxxx yyyy " case, do not include \"
		if (s[i] == '\'' || s[i] == '"') && (inString%2 == 0 || lastQuot == s[i]) {
			inString++
			escape = 0
			lastQuot = s[i]
		} else {
			if !isSpace(s[i]) {
				escape = 0
			}
		}
		if inString%2 == 0 && isSpace(s[i]) {
			if start == i { //escape continuous space
				start += lenSep
			} else {
				a[na] = s[start+escape : i-escape]
				na++
				start = i + lenSep
				i += lenSep - 1
			}
		}
	}
	if start < len(s) {
		a[na] = s[start+escape : len(s)-escape]
	} else {
		na--
	}

	return a[0 : na+1]
}

//FormatLineHead ensure all lines of s are lead with linehead string
func FormatLineHead(s, lineHead string) (r string) {
	r = lineHeadExpr.ReplaceAllString(s, lineHead)
	return
}
