package gogp

import (
	"fmt"
	"testing"
)

func TestAllRegexpSyntax(t *testing.T) {
	expAll := compileMultiRegexps(res...)
	fmt.Printf("%#v\n", expAll.SubexpNames())
}
