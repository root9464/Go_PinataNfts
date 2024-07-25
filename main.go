package main

import (
	"root/src/utils"
)

func main() {
	_, printPinataResponse := utils.GetPinataResponseFuncs()
	printPinataResponse()
}
