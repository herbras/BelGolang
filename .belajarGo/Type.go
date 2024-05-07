package main

import (
	"fmt"
)

var positiveNumber uint8 = 89
var negativeNumber = -1243423644
var exist bool = true

const (
	it             = "ibrah"
	tof            = "him"
	er       bool  = true
	numeric  uint8 = 1
	floatNum       = 2.2
)

func typeExample() {
	fmt.Printf("bilangan positif: %d\n", positiveNumber)
	fmt.Printf("bilangan negatif: %d\n", negativeNumber)
	fmt.Printf("exist? %t \n", exist)
	fmt.Print(it, tof, er, numeric, floatNum)
}
