package main

import (
	"fmt"
	"strconv"
)

var result string

func main() {
	number := 200                        //要转换的数字
	Hex := 3                             //转换的进制
	result = (transformHex(number, Hex)) //将 数字和进制参数带入transformHex函数中
	fmt.Print(result)                    //显示结果
}

func transformHex(a int, Hex int) string {
	result := ""
	for ; a > 0; a /= Hex {
		last := a % Hex
		result = strconv.Itoa(last) + result //strconv.Itoa将"last"整形转换为字符串型
	}
	return result //返回计算后结果
}
