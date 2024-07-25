package main

import (
	"root/src/utils"
)

func main() {
	printPinataResponse, _ := utils.GetPinataResponseFuncs()
	printPinataResponse()
}
