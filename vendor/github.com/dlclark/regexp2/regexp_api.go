// adapt apis to std.regexp

package regexp2

import (
	"fmt"
)

func makeRepFunc(f func(string) string) MatchEvaluator {
	return func(m Match) string {
		return f(m.String())
	}
}

func (re *Regexp) ReplaceAllString(input string, repace string) string {
	rep, err := re.Replace(input, repace, 0, -1)
	if err != nil {
		fmt.Println(err)
		return input
	}
	return rep
}

func (re *Regexp) ReplaceAllStringFunc(input string, f func(string) string) string {
	rep, err := re.ReplaceFunc(input, makeRepFunc(f), 0, -1)
	if err != nil {
		fmt.Println(err)
		return input
	}
	return rep
}

func (re *Regexp) FindAllStringSubmatch(input string, num int) [][]string {
	m, err := re.FindStringMatch(input)
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
		m, err = re.FindNextMatch(m)
	}
	return r
}

func (re *Regexp) IsMatchString(input string) bool {
	if ok, err := re.MatchString(input); err == nil {
		return ok
	}
	return false
}

func (re *Regexp) SubexpNames() []string {
	return re.GetGroupNames()
}
