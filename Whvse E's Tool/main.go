package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		uid      int
		user     int
		password string
		oVD      int
	)
	fmt.Printf("输入学号: \t")
	fmt.Scan(&user)
	fmt.Printf("输入E学堂密码: \t")
	fmt.Scan(&password)
	uid = GetUid(user, password)
	fmt.Printf("你的Uid为：%d \t\n", uid)
	fmt.Printf("输入初始VD: \t")
	fmt.Scan(&oVD)
	fmt.Print("正在提交视频时间")
	mulitPosting(uid, oVD)
}

//获取学生UID
func GetUid(user int, password string) (uid int) {
	url := "http://wrggka.whvcse.edu.cn/api/M_User/Login?username=" + strconv.Itoa(user) + "&password=" + password + "&accessKey=1&secretKey=1" //登陆获取UID
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	date := string(body)
	resultuid := Validate(date)
	uid = resultuid
	/* 	date = Validate(date) */
	defer res.Body.Close()
	return uid
}

//数据转换
func Validate(date string) int {
	findUid := regexp.MustCompile(`\d{4}`)
	result0 := findUid.FindAllStringSubmatch(date, -1)
	result1 := (result0[0][0]) //取二位数组的00号元素，即UID
	result2 := strings.TrimLeft(result1, "0")
	fResult, _ := strconv.Atoi(result2) //字符串类转化为整数型
	return fResult
}

//提=提交观看视频时间给服务器
func submitTime(uid int, oVD int) {
	var oVT int
	oVT = rand.Intn(10) * 100 //随机观看时间 100-1000s
	VT := strconv.Itoa(oVT)   //两数据提前换成字符型
	VD := strconv.Itoa(oVD)
	fmt.Print(VT)
	//发送时间请求
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/IsNoStudyvideo?userId=" + strconv.Itoa(uid) + "&videoid=" + VD + "&videotime=" + VT + "&accessKey=1&secretKey=1"
	client := &http.Client{}
	req, _ := http.NewRequest("Get", url, nil)
	req.Header.Set("User-Agent:", "Dalvik/2.1.0 (Linux; U; Android 9.0; HUAWEI NOVA5 Build/RDD0M)")
	req.Header.Add("Connection:", "Keep-Alive")
	req.Header.Add("Accept-Encoding:", "gzip")
	res, _ := client.Get(url)
	fmt.Print(res)
	/* 	fmt.Print(url) */
}

//并发提交视频观看时间
func mulitPosting(uid, oVD int) {
	for i := oVD; i <= (oVD + 100); i++ {
		submitTime(uid, i)
		time.Sleep(30000)
		fmt.Printf("已提交次数%d\n", i-oVD)

	}
	fmt.Print("已完成改课程的作业提交")
}
