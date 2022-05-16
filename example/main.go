package main

import (
	"fmt"
	"unsafe"

	"github.com/probably-not/q/nano"
	"github.com/probably-not/q/pico"
)

func main() {
	picoQ := pico.NewQ()
	nanoQ := nano.NewQ()
	fmt.Println("=========================== Queue Memory Sizes ===========================")
	fmt.Println("PicoQ:", unsafe.Sizeof(picoQ), "bytes")
	fmt.Println("NanoQ:", unsafe.Sizeof(nanoQ), "bytes")
	fmt.Println("==========================================================================")
}
