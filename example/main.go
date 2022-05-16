package main

import (
	"fmt"
	"unsafe"

	"github.com/probably-not/q/micro"
	"github.com/probably-not/q/nano"
	"github.com/probably-not/q/pico"
)

func main() {
	picoQ := pico.NewQ()
	nanoQ := nano.NewQ()
	microQ := micro.NewQ(6)
	fmt.Println("=========================== Queue Memory Sizes ===========================")
	fmt.Println("PicoQ:", unsafe.Sizeof(picoQ), "bytes")
	fmt.Println("NanoQ:", unsafe.Sizeof(nanoQ), "bytes")
	fmt.Println("MicroQ:", unsafe.Sizeof(microQ), "bytes")
	fmt.Println("==========================================================================")
}
