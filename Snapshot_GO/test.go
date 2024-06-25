package main

import (
	"fmt"
	"sync/atomic"
)

var a atomic.Bool

func main() {
	a.CompareAndSwap(true, true)
	fmt.Println(a)
}
