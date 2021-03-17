// adapt apis to std.regexp

package regexp2

import (
	"fmt"
	"io"
)

func CompileStd(expr string) (*RegexpStd, error) {
	re, err := Compile(expr, RE2)
	return re.StdRegexp(), err
}

func MustCompileStd(expr string) *RegexpStd {
	re, err := CompileStd(expr)
	if err != nil {
		panic(err)
	}
	return re
}

type RegexpStd struct {
	p *Regexp
}

func (re *Regexp) StdRegexp() *RegexpStd {
	return &RegexpStd{p: re}
}

func makeRepFunc(f func(string) string) MatchEvaluator {
	return func(m Match) string {
		return f(m.String())
	}
}

func (re *RegexpStd) Regexp2() *Regexp {
	return re.p
}

// String returns the source text used to compile the regular expression.
func (re *RegexpStd) String() string {
	return ""
}

// Copy returns a new StdRegexp object copied from re.
// Calling Longest on one copy does not affect another.
//
// Deprecated: In earlier releases, when using a StdRegexp in multiple goroutines,
// giving each goroutine its own copy helped to avoid lock contention.
// As of Go 1.12, using Copy is no longer necessary to avoid lock contention.
// Copy may still be appropriate if the reason for its use is to make
// two copies with different Longest settings.
func (re *RegexpStd) Copy() *RegexpStd {
	re2 := *re
	return &re2
}

// Longest makes future searches prefer the leftmost-longest match.
// That is, when matching against text, the regexp returns a match that
// begins as early as possible in the input (leftmost), and among those
// it chooses a match that is as long as possible.
// This method modifies the StdRegexp and may not be called concurrently
// with any other methods.
func (re *RegexpStd) Longest() {
}

// SubexpNames returns the names of the parenthesized subexpressions
// in this StdRegexp. The name for the first sub-expression is names[1],
// so that if m is a match slice, the name for m[i] is SubexpNames()[i].
// Since the StdRegexp as a whole cannot be named, names[0] is always
// the empty string. The slice should not be modified.
func (re *RegexpStd) SubexpNames() []string {
	return re.p.GetGroupNames()
}

// SubexpIndex returns the index of the first subexpression with the given name,
// or -1 if there is no subexpression with that name.
//
// Note that multiple subexpressions can be written using the same name, as in
// (?P<bob>a+)(?P<bob>b+), which declares two subexpressions named "bob".
// In this case, SubexpIndex returns the index of the leftmost such subexpression
// in the regular expression.
func (re *RegexpStd) SubexpIndex(name string) int {
	return 0
}

// LiteralPrefix returns a literal string that must begin any match
// of the regular expression re. It returns the boolean true if the
// literal string comprises the entire regular expression.
func (re *RegexpStd) LiteralPrefix() (prefix string, complete bool) {
	panic("")
}

// MatchReader reports whether the text returned by the RuneReader
// contains any match of the regular expression re.
func (re *RegexpStd) MatchReader(r io.RuneReader) bool {
	panic("")
}

// MatchString reports whether the string s
// contains any match of the regular expression re.
func (re *RegexpStd) MatchString(s string) bool {
	if ok, err := re.p.MatchString(s); err == nil {
		return ok
	}
	return false
}

// Match reports whether the byte slice b
// contains any match of the regular expression re.
func (re *RegexpStd) Match(b []byte) bool {
	panic("")
}

// MatchString reports whether the string s
// contains any match of the regular expression pattern.
// More complicated queries need to use Compile and the full StdRegexp interface.
func MatchString(pattern string, s string) (matched bool, err error) {
	re, err := CompileStd(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

// Match reports whether the byte slice b
// contains any match of the regular expression pattern.
// More complicated queries need to use Compile and the full StdRegexp interface.
func MatchStd(pattern string, b []byte) (matched bool, err error) {
	re, err := CompileStd(pattern)
	if err != nil {
		return false, err
	}
	return re.Match(b), nil
}

// ReplaceAllString returns a copy of src, replacing matches of the StdRegexp
// with the replacement string repl. Inside repl, $ signs are interpreted as
// in Expand, so for instance $1 represents the text of the first submatch.
func (re *RegexpStd) ReplaceAllString(src, repl string) string {
	rep, err := re.p.Replace(src, repl, 0, -1)
	if err != nil {
		fmt.Println(err)
		return src
	}
	return rep
}

// ReplaceAllLiteralString returns a copy of src, replacing matches of the StdRegexp
// with the replacement string repl. The replacement repl is substituted directly,
// without using Expand.
func (re *RegexpStd) ReplaceAllLiteralString(src, repl string) string {
	// return string(re.ReplaceAll(nil, src, 2, func(dst []byte, match []int) []byte {
	// 	return append(dst, repl...)
	// }))
	panic("")
}

// ReplaceAllStringFunc returns a copy of src in which all matches of the
// StdRegexp have been replaced by the return value of function repl applied
// to the matched substring. The replacement returned by repl is substituted
// directly, without using Expand.
func (re *RegexpStd) ReplaceAllStringFunc(src string, repl func(string) string) string {
	rep, err := re.p.ReplaceFunc(src, makeRepFunc(repl), 0, -1)
	if err != nil {
		fmt.Println(err)
		return src
	}
	return rep
}

// ReplaceAll returns a copy of src, replacing matches of the StdRegexp
// with the replacement text repl. Inside repl, $ signs are interpreted as
// in Expand, so for instance $1 represents the text of the first submatch.
func (re *RegexpStd) ReplaceAll(src, repl []byte) []byte {
	// n := 2
	// if bytes.IndexByte(repl, '$') >= 0 {
	// 	n = 2 * (re.numSubexp + 1)
	// }
	// srepl := ""
	// b := re.replaceAll(src, "", n, func(dst []byte, match []int) []byte {
	// 	if len(srepl) != len(repl) {
	// 		srepl = string(repl)
	// 	}
	// 	return re.expand(dst, srepl, src, "", match)
	// })
	// return b
	panic("")
}

// ReplaceAllLiteral returns a copy of src, replacing matches of the StdRegexp
// with the replacement bytes repl. The replacement repl is substituted directly,
// without using Expand.
func (re *RegexpStd) ReplaceAllLiteral(src, repl []byte) []byte {
	// return re.replaceAll(src, "", 2, func(dst []byte, match []int) []byte {
	// 	return append(dst, repl...)
	// })
	panic("")
}

// ReplaceAllFunc returns a copy of src in which all matches of the
// StdRegexp have been replaced by the return value of function repl applied
// to the matched byte slice. The replacement returned by repl is substituted
// directly, without using Expand.
func (re *RegexpStd) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte {
	// return re.replaceAll(src, "", 2, func(dst []byte, match []int) []byte {
	// 	return append(dst, repl(src[match[0]:match[1]])...)
	// })
	panic("")
}

// QuoteMeta returns a string that escapes all regular expression metacharacters
// inside the argument text; the returned string is a regular expression matching
// the literal text.
func QuoteMeta(s string) string {
	// // A byte loop is correct because all metacharacters are ASCII.
	// var i int
	// for i = 0; i < len(s); i++ {
	// 	if special(s[i]) {
	// 		break
	// 	}
	// }
	// // No meta characters found, so return original string.
	// if i >= len(s) {
	// 	return s
	// }

	// b := make([]byte, 2*len(s)-i)
	// copy(b, s[:i])
	// j := i
	// for ; i < len(s); i++ {
	// 	if special(s[i]) {
	// 		b[j] = '\\'
	// 		j++
	// 	}
	// 	b[j] = s[i]
	// 	j++
	// }
	// return string(b[:j])
	panic("")
}

// Find returns a slice holding the text of the leftmost match in b of the regular expression.
// A return value of nil indicates no match.
func (re *RegexpStd) Find(b []byte) []byte {
	// var dstCap [2]int
	// a := re.doExecute(nil, b, "", 0, 2, dstCap[:0])
	// if a == nil {
	// 	return nil
	// }
	// return b[a[0]:a[1]:a[1]]
	panic("")
}

// FindIndex returns a two-element slice of integers defining the location of
// the leftmost match in b of the regular expression. The match itself is at
// b[loc[0]:loc[1]].
// A return value of nil indicates no match.
func (re *RegexpStd) FindIndex(b []byte) (loc []int) {
	// a := re.doExecute(nil, b, "", 0, 2, nil)
	// if a == nil {
	// 	return nil
	// }
	// return a[0:2]
	panic("")
}

// FindString returns a string holding the text of the leftmost match in s of the regular
// expression. If there is no match, the return value is an empty string,
// but it will also be empty if the regular expression successfully matches
// an empty string. Use FindStringIndex or FindStringSubmatch if it is
// necessary to distinguish these cases.
func (re *RegexpStd) FindString(s string) string {
	// var dstCap [2]int
	// a := re.doExecute(nil, nil, s, 0, 2, dstCap[:0])
	// if a == nil {
	// 	return ""
	// }
	// return s[a[0]:a[1]]
	panic("")
}

// FindStringIndex returns a two-element slice of integers defining the
// location of the leftmost match in s of the regular expression. The match
// itself is at s[loc[0]:loc[1]].
// A return value of nil indicates no match.
func (re *RegexpStd) FindStringIndex(s string) (loc []int) {
	// a := re.doExecute(nil, nil, s, 0, 2, nil)
	// if a == nil {
	// 	return nil
	// }
	// return a[0:2]
	panic("")
}

// FindReaderIndex returns a two-element slice of integers defining the
// location of the leftmost match of the regular expression in text read from
// the RuneReader. The match text was found in the input stream at
// byte offset loc[0] through loc[1]-1.
// A return value of nil indicates no match.
func (re *RegexpStd) FindReaderIndex(r io.RuneReader) (loc []int) {
	// a := re.doExecute(r, nil, "", 0, 2, nil)
	// if a == nil {
	// 	return nil
	// }
	// return a[0:2]
	panic("")
}

// FindSubmatch returns a slice of slices holding the text of the leftmost
// match of the regular expression in b and the matches, if any, of its
// subexpressions, as defined by the 'Submatch' descriptions in the package
// comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindSubmatch(b []byte) [][]byte {
	// var dstCap [4]int
	// a := re.doExecute(nil, b, "", 0, re.prog.NumCap, dstCap[:0])
	// if a == nil {
	// 	return nil
	// }
	// ret := make([][]byte, 1+re.numSubexp)
	// for i := range ret {
	// 	if 2*i < len(a) && a[2*i] >= 0 {
	// 		ret[i] = b[a[2*i]:a[2*i+1]:a[2*i+1]]
	// 	}
	// }
	// return ret
	panic("")
}

// Expand appends template to dst and returns the result; during the
// append, Expand replaces variables in the template with corresponding
// matches drawn from src. The match slice should have been returned by
// FindSubmatchIndex.
//
// In the template, a variable is denoted by a substring of the form
// $name or ${name}, where name is a non-empty sequence of letters,
// digits, and underscores. A purely numeric name like $1 refers to
// the submatch with the corresponding index; other names refer to
// capturing parentheses named with the (?P<name>...) syntax. A
// reference to an out of range or unmatched index or a name that is not
// present in the regular expression is replaced with an empty slice.
//
// In the $name form, name is taken to be as long as possible: $1x is
// equivalent to ${1x}, not ${1}x, and, $10 is equivalent to ${10}, not ${1}0.
//
// To insert a literal $ in the output, use $$ in the template.
func (re *RegexpStd) Expand(dst []byte, template []byte, src []byte, match []int) []byte {
	//return re.expand(dst, string(template), src, "", match)
	panic("")
}

// ExpandString is like Expand but the template and source are strings.
// It appends to and returns a byte slice in order to give the calling
// code control over allocation.
func (re *RegexpStd) ExpandString(dst []byte, template string, src string, match []int) []byte {
	//return re.expand(dst, template, nil, src, match)
	panic("")
}

// FindSubmatchIndex returns a slice holding the index pairs identifying the
// leftmost match of the regular expression in b and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' and 'Index' descriptions
// in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindSubmatchIndex(b []byte) []int {
	//return re.pad(re.doExecute(nil, b, "", 0, re.prog.NumCap, nil))
	panic("")
}

// FindStringSubmatch returns a slice of strings holding the text of the
// leftmost match of the regular expression in s and the matches, if any, of
// its subexpressions, as defined by the 'Submatch' description in the
// package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindStringSubmatch(s string) []string {
	// var dstCap [4]int
	// a := re.doExecute(nil, nil, s, 0, re.prog.NumCap, dstCap[:0])
	// if a == nil {
	// 	return nil
	// }
	// ret := make([]string, 1+re.numSubexp)
	// for i := range ret {
	// 	if 2*i < len(a) && a[2*i] >= 0 {
	// 		ret[i] = s[a[2*i]:a[2*i+1]]
	// 	}
	// }
	// return ret
	panic("")
}

// FindStringSubmatchIndex returns a slice holding the index pairs
// identifying the leftmost match of the regular expression in s and the
// matches, if any, of its subexpressions, as defined by the 'Submatch' and
// 'Index' descriptions in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindStringSubmatchIndex(s string) []int {
	//return re.pad(re.doExecute(nil, nil, s, 0, re.prog.NumCap, nil))
	panic("")
}

// FindReaderSubmatchIndex returns a slice holding the index pairs
// identifying the leftmost match of the regular expression of text read by
// the RuneReader, and the matches, if any, of its subexpressions, as defined
// by the 'Submatch' and 'Index' descriptions in the package comment. A
// return value of nil indicates no match.
func (re *RegexpStd) FindReaderSubmatchIndex(r io.RuneReader) []int {
	//return re.pad(re.doExecute(r, nil, "", 0, re.prog.NumCap, nil))
	panic("")
}

// FindAll is the 'All' version of Find; it returns a slice of all successive
// matches of the expression, as defined by the 'All' description in the
// package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAll(b []byte, n int) [][]byte {
	// if n < 0 {
	// 	n = len(b) + 1
	// }
	// var result [][]byte
	// re.allMatches("", b, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([][]byte, 0, startSize)
	// 	}
	// 	result = append(result, b[match[0]:match[1]:match[1]])
	// })
	// return result
	panic("")
}

// FindAllIndex is the 'All' version of FindIndex; it returns a slice of all
// successive matches of the expression, as defined by the 'All' description
// in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllIndex(b []byte, n int) [][]int {
	// if n < 0 {
	// 	n = len(b) + 1
	// }
	// var result [][]int
	// re.allMatches("", b, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([][]int, 0, startSize)
	// 	}
	// 	result = append(result, match[0:2])
	// })
	// return result
	panic("")
}

// FindAllString is the 'All' version of FindString; it returns a slice of all
// successive matches of the expression, as defined by the 'All' description
// in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllString(s string, n int) []string {
	// if n < 0 {
	// 	n = len(s) + 1
	// }
	// var result []string
	// re.allMatches(s, nil, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([]string, 0, startSize)
	// 	}
	// 	result = append(result, s[match[0]:match[1]])
	// })
	// return result
	panic("")
}

// FindAllStringIndex is the 'All' version of FindStringIndex; it returns a
// slice of all successive matches of the expression, as defined by the 'All'
// description in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllStringIndex(s string, n int) [][]int {
	// if n < 0 {
	// 	n = len(s) + 1
	// }
	// var result [][]int
	// re.allMatches(s, nil, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([][]int, 0, startSize)
	// 	}
	// 	result = append(result, match[0:2])
	// })
	// return result
	panic("")
}

// FindAllSubmatch is the 'All' version of FindSubmatch; it returns a slice
// of all successive matches of the expression, as defined by the 'All'
// description in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllSubmatch(b []byte, n int) [][][]byte {
	// if n < 0 {
	// 	n = len(b) + 1
	// }
	// var result [][][]byte
	// re.allMatches("", b, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([][][]byte, 0, startSize)
	// 	}
	// 	slice := make([][]byte, len(match)/2)
	// 	for j := range slice {
	// 		if match[2*j] >= 0 {
	// 			slice[j] = b[match[2*j]:match[2*j+1]:match[2*j+1]]
	// 		}
	// 	}
	// 	result = append(result, slice)
	// })
	// return result
	panic("")
}

// FindAllSubmatchIndex is the 'All' version of FindSubmatchIndex; it returns
// a slice of all successive matches of the expression, as defined by the
// 'All' description in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllSubmatchIndex(b []byte, n int) [][]int {
	// if n < 0 {
	// 	n = len(b) + 1
	// }
	// var result [][]int
	// re.allMatches("", b, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([][]int, 0, startSize)
	// 	}
	// 	result = append(result, match)
	// })
	// return result
	panic("")
}

// FindAllStringSubmatch is the 'All' version of FindStringSubmatch; it
// returns a slice of all successive matches of the expression, as defined by
// the 'All' description in the package comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllStringSubmatch(s string, n int) [][]string {
	m, err := re.p.FindStringMatch(s)
	if err != nil {
		println(err.Error())
		return nil
	}

	groups := m.Groups()
	mm := make([]string, 0, len(groups))
	for i := 0; i < len(groups); i++ {
		mm = append(mm, (&groups[i]).String())
	}

	r := [][]string{mm}
	return r
}

// FindAllStringSubmatchIndex is the 'All' version of
// FindStringSubmatchIndex; it returns a slice of all successive matches of
// the expression, as defined by the 'All' description in the package
// comment.
// A return value of nil indicates no match.
func (re *RegexpStd) FindAllStringSubmatchIndex(s string, n int) [][]int {
	// if n < 0 {
	// 	n = len(s) + 1
	// }
	// var result [][]int
	// re.allMatches(s, nil, n, func(match []int) {
	// 	if result == nil {
	// 		result = make([][]int, 0, startSize)
	// 	}
	// 	result = append(result, match)
	// })
	// return result
	panic("")
}

// Split slices s into substrings separated by the expression and returns a slice of
// the substrings between those expression matches.
//
// The slice returned by this method consists of all the substrings of s
// not contained in the slice returned by FindAllString. When called on an expression
// that contains no metacharacters, it is equivalent to strings.SplitN.
//
// Example:
//   s := regexp.MustCompile("a*").Split("abaabaccadaaae", 5)
//   // s: ["", "b", "b", "c", "cadaaae"]
//
// The count determines the number of substrings to return:
//   n > 0: at most n substrings; the last substring will be the unsplit remainder.
//   n == 0: the result is nil (zero substrings)
//   n < 0: all substrings
func (re *RegexpStd) Split(s string, n int) []string {

	// if n == 0 {
	// 	return nil
	// }

	// if len(re.expr) > 0 && len(s) == 0 {
	// 	return []string{""}
	// }

	// matches := re.FindAllStringIndex(s, n)
	// strings := make([]string, 0, len(matches))

	// beg := 0
	// end := 0
	// for _, match := range matches {
	// 	if n > 0 && len(strings) >= n-1 {
	// 		break
	// 	}

	// 	end = match[0]
	// 	if match[1] != 0 {
	// 		strings = append(strings, s[beg:end])
	// 	}
	// 	beg = match[1]
	// }

	// if end != len(s) {
	// 	strings = append(strings, s[beg:])
	// }

	// return strings
	panic("")
}
