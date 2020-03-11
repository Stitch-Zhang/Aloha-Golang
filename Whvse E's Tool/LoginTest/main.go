package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var (
		user     int
		password string
	)
	fmt.Print("user")
	fmt.Scan(&user)
	fmt.Print("PAssword")
	fmt.Scan(&password)
	url := "http://wrggka.whvcse.edu.cn/api/M_User/Login?username=" + strconv.Itoa(user) + "&password=" + password + "&accessKey=1&secretKey=1" //登陆获取UID
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	date := string(body)
	if strings.Contains(date, "status\":\"0\"") {
		fmt.Print("密码错误")
		os.Exit(0)
	}
	fmt.Print("登陆成功")
	resultuid := Validate(date)
	uid := resultuid
	/* 	date = Validate(date) */
	defer res.Body.Close()
	fmt.Printf("UID:%d", uid)
}

func Validate(date string) int {
	findUid := regexp.MustCompile(`\d{4}`)
	result0 := findUid.FindAllStringSubmatch(date, -1)
	result1 := (result0[0][0]) //取二位数组的00号元素，即UID
	result2 := strings.TrimLeft(result1, "0")
	fResult, _ := strconv.Atoi(result2) //字符串类转化为整数型
	return fResult
}
