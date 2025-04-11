package main

import (
	"art/functions"
	"fmt"
)

func main() {
	input := "###----\n___"
	fmt.Println(functions.EncodeString(input, true))
}