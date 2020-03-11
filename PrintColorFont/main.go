package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	for i := 0; i <= 10; i++ {
		green := color.New(color.FgWhite)   //定义前景色
		greenbg := green.Add(color.BgGreen) //di定义背景色
		greenbg.Printf("This is %d", i)
	}
	var t string
	fmt.Print("TEST")
	fmt.Scan(&t)
	os.Exit(0)
}
