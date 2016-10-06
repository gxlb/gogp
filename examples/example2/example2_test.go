package example2

import (
	"fmt"
	"testing"

	_ "github.com/vipally/gogp/auto" //auto run gogp tool at GoPath in test process
)

func TestPersonStack(t *testing.T) {
	stack := NewPersonStack()
	stack.Push(&Person{"tom", 10})
	stack.Push(&Person{"jim", 12})
	stack.Push(&Person{"jeep", 6})
	fmt.Println(stack.Show())
}
