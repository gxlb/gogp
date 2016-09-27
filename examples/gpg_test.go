package example

import (
	"testing"

	"github.com/vipally/gogp" //auto run gogp tool at current path in test process
)

func init() {
	gogp.Work("./example2") //run gogp at another path
}

//run gogp tool to auto-generate go file(s) in test process
func TestCallGogp(t *testing.T) {

}
