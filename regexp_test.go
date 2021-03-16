package gogp

import (
	"fmt"
	"regexp"
	"testing"
)

func TestAllRegexpSyntax(t *testing.T) {

	gogpExpAll := regexp.MustCompile(fmt.Sprintf("%s", expTxtTodoReplace))
	txt := `
	<todoReplace>
`

	fmt.Println("gogpExpAll match", gogpExpAll.MatchString(txt))

}
