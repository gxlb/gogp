// adapt apis to std.regexp

package regexp2

import (
	"fmt"
)

func CompileStd(expr string) (*StdRegexp, error) {
	re, err := Compile(expr, RE2)
	return re.StdRegexp(), err
}

func MustCompileStd(expr string) *StdRegexp {
	re, err := CompileStd(expr)
	if err != nil {
		panic(err)
	}
	return re
}

type StdRegexp struct {
	p *Regexp
}

func (re *Regexp) StdRegexp() *StdRegexp {
	return &StdRegexp{p: re}
}

func makeRepFunc(f func(string) string) MatchEvaluator {
	return func(m Match) string {
		return f(m.String())
	}
}

func (re *StdRegexp) ReplaceAllString(input string, repace string) string {
	rep, err := re.p.Replace(input, repace, 0, -1)
	if err != nil {
		fmt.Println(err)
		return input
	}
	return rep
}

func (re *StdRegexp) Regexp2() *Regexp {
	return re.p
}

func (re *StdRegexp) ReplaceAllStringFunc(input string, f func(string) string) string {
	rep, err := re.p.ReplaceFunc(input, makeRepFunc(f), 0, -1)
	if err != nil {
		fmt.Println(err)
		return input
	}
	return rep
}

func (re *StdRegexp) FindAllStringSubmatch(input string, num int) [][]string {
	m, err := re.p.FindStringMatch(input)
	if err != nil {
		return nil
	}
	r := make([][]string, 0, len(m.Captures))
	for m != nil && err == nil {
		mm := make([]string, 0, len(m.otherGroups)+1)
		mm = append(mm, m.String())
		for i := 0; i < len(m.otherGroups); i++ {
			mm = append(mm, (&m.otherGroups[i]).String())
		}
		r = append(r, mm)
		m, err = re.p.FindNextMatch(m)
	}
	return r
}

func (re *StdRegexp) MatchString(input string) bool {
	if ok, err := re.p.MatchString(input); err == nil {
		return ok
	}
	return false
}

func (re *StdRegexp) SubexpNames() []string {
	return re.p.GetGroupNames()
}
