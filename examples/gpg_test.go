package example

import (
	"os"
	"testing"

	"github.com/vipally/gogp"
)

//run gogp tool to auto-generate go file(s) in test process
func TestCallGogp(t *testing.T) {
	if dir, err := os.Getwd(); err == nil {
		gogp.Work(dir)
		gogp.Work(dir + "/example2")
	} else {
		panic(err)
	}
}
