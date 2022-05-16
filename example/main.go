package main

import (
	"fmt"
	"unsafe"

	"github.com/probably-not/q/pico"
)

func main() {
	picoQ := pico.NewQ()
	fmt.Println("=========================== Queue Memory Sizes ===========================")
	fmt.Println("PicoQ:", unsafe.Sizeof(picoQ), "bytes")
	fmt.Println("==========================================================================")
}
