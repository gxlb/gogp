package examples

import (
	"fmt"
	"testing"

	//_ "github.com/vipally/gogp/auto" //auto run gogp tool at current path in test process
)

func TestCallGogp(t *testing.T) {
	l := NewIntList()
	a := []int{9, 7, 3, 5, 1, 7, 8}
	//fmt.Printf("%p\n", l)
	for _, v := range a {
		l.PushBack(v)
	}
	showList(l, "before sort")
	//showListReverse(l, "Rbefore sort")
	l.Sort()
	showList(l, "after  sort")
}

func showList(l *IntList, head string) {
	fmt.Printf("%s: ", head)
	for v, n := l.Visitor(), 0; v.Next(); {
		n++
		if n > 10 {
			break
		}
		fmt.Printf("%d ", v.Get().Get())
	}
	fmt.Printf("\n")
}

func showListReverse(l *IntList, head string) {
	fmt.Printf("%s: ", head)
	for v := l.Visitor(); v.Prev(); {
		fmt.Printf("%d ", v.Get().Get())
	}
	fmt.Printf("\n")
}
