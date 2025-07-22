package main

import (
	"fmt"
	"strconv"
)

func main1() {
	fmt.Printf("%b\n", '€')
	for _, v := range []byte("€") {
		fmt.Printf("%b", v)
	}
	// fmt.Printf("%v\n", []byte("€"))

	str_base10 := "11"
	bridge, _ := strconv.Atoi(str_base10)
	str_base16 := strconv.FormatInt(int64(bridge), 16)
	fmt.Println(str_base16)
}
