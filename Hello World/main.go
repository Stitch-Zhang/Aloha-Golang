package main

//导入fmt"包
import "fmt"

//入口Main函数
//在此导入其他的函数（helloworld函数已导入）
func main() {
	fmt.Print("我的GO语言之旅开始啦\n")
	helloWorld()
}

//HelloWorld函数
func helloWorld() {
	fmt.Print("Hello World!!")
}
