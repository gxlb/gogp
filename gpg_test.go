package gogp_test

import (
	"fmt"
	"regexp"
	"testing"
)

//run gogp tool to auto-generate go file(s) in test process
func TestCallGogp(t *testing.T) {
	r := regexp.MustCompile("hello|he")
	s := r.FindAllString("I think he is say hello to HE hehe", -1)
	fmt.Printf("%#v", s)
}
