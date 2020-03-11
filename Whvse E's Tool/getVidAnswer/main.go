package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PaperStruct []struct {
	ItemID        string `json:"ItemID"`
	ItemTitle     string `json:"ItemTitle"`
	ItemType      string `json:"ItemType"`
	ItemCategory  string `json:"ItemCategory"`
	ItemResolving string `json:"ItemResolving"`
	ItemOptions   []struct {
		ItemID        string `json:"ItemID"`
		ItemTitle     string `json:"ItemTitle"`
		ItemIsCorrect string `json:"ItemIsCorrect"`
	} `json:"ItemOptions"`
}

var (
	findPaper     = regexp.MustCompile(`[a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9]`)
	uid           = 8128
	courseId      = 751
	courseClassId = 522
)

func main() {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/GetCourseSPZT?userId=" + strconv.Itoa(uid) + "&courseId=" + strconv.Itoa(courseId) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	paperID := findPaper.FindAllStringSubmatch(string(body), -1)
	fmt.Printf("一号试卷ID为：\t%s\n", paperID[0][0])
	fmt.Printf("二号试卷ID为: \t%s\n", paperID[1][0])
	readyExam(paperID[0][0])
	readyExam(paperID[1][0])
	fmt.Print("Press Enter to exit")
	var A int
	fmt.Scan(&A)
	_ = A
	os.Exit(0)

}

//做第一张试卷
func readyExam(paperID string) {
	url1 := "http://wrggka.whvcse.edu.cn/api/M_Course/GetChapterTestInfo?userId=" + strconv.Itoa(uid) + "&courseId=" + strconv.Itoa(courseId) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&chapterId=0&paperId=" + paperID + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url1)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("已进入考场\n")
	body, _ := ioutil.ReadAll(resp.Body)
	examCountId0 := findPaper.FindAllStringSubmatch(string(body), -1)
	examCountId := examCountId0[1][0]

	fmt.Printf("试卷的密封ID为: %s \n", examCountId)
	onExam(paperID, examCountId)
}

func onExam(paperId, examCountId string) {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/GetPaperQuestions3?paperId=" + paperId + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	//请求过快
	if err != nil {
		fmt.Printf("考试作弊被抓！！！！\t%s\n", err)
	}
	//判断返回服务器状态码是否正常
	if !strings.Contains(resp.Status, "200") {
		fmt.Print("本场考试服务器go die 了！！\n")
		fmt.Print("本场考试服务器go die 了！！\n")
		fmt.Print("本场考试服务器go die 了！！\n")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var paperStruct PaperStruct
	json.Unmarshal(body, &paperStruct)
	//计算题目数量
	sum := strings.Count(string(body), "ItemType") //string(body)中包含一次"ItemType"计数一次
	fmt.Printf("确认题目数量是否为:%d\n", sum)

	//判断题目类型
	for i := 0; i < sum; i++ {
		questionType, _ := strconv.Atoi(paperStruct[i].ItemType)
		fmt.Printf("第%d题类型为:", i+1)
		switch questionType {
		case 6:
			var (
				finalAnswer string
				questionId  string
			)
			fmt.Print("判断题" + "\n")
			for k := 0; k < 2; k++ {
				correct := paperStruct[i].ItemOptions[k].ItemIsCorrect
				if strings.Contains(correct, "1") {
					finalAnswer = paperStruct[i].ItemOptions[k].ItemID
					questionId = paperStruct[i].ItemID
					/* 					fmt.Printf("第%d题的Question ID Is %s \n", i+1, questionId)
					   					fmt.Printf("第%d题的Answer ID is %s \n", i+1, finalAnswer) */
					switch k {
					case 0:
						fmt.Printf("第%d题答案为\t \"对\"\n", i+1)
					case 1:
						fmt.Printf("第%d题答案为:\t \"错\"\n", i+1)
					}
				}
			}
			SubmitAnswer(questionId, finalAnswer, paperId, examCountId)
		case 1:
			fmt.Print("单选题" + "\n")
			var (
				finalAnswer string
				questionId  string
			)
			lenth := len(paperStruct[i].ItemOptions) //len来判断选择题选项的个数实在是太香！！
			for k := 0; k < lenth; k++ {
				correct := paperStruct[i].ItemOptions[k].ItemIsCorrect
				if strings.Contains(correct, "1") {
					finalAnswer = paperStruct[i].ItemOptions[k].ItemID
					questionId = paperStruct[i].ItemID
					switch k {
					case 0:
						fmt.Printf("第%d题答案为:\tA\n", i+1)
					case 1:
						fmt.Printf("第%d题答案为:\tB\n", i+1)
					case 2:
						fmt.Printf("第%d题答案为:\tC\n", i+1)
					case 3:
						fmt.Printf("第%d题答案为:\tD\n", i+1)

					}
				}

			}
			SubmitAnswer(questionId, finalAnswer, paperId, examCountId)
		case 2:
			fmt.Print("多选题" + "\n")
			fmt.Printf("第%d题答案为:\t", i+1)
			var (
				finalAnswer0 string
				finalAnswer  string
				questionId   string
			)
			lenth := len(paperStruct[i].ItemOptions) //len来判断选择题选项的个数实在是太香！！
			for k := 0; k < lenth; k++ {
				correct := paperStruct[i].ItemOptions[k].ItemIsCorrect
				if strings.Contains(correct, "1") {
					answer := paperStruct[i].ItemOptions[k].ItemID
					questionId = paperStruct[i].ItemID
					/* 					fmt.Printf("第%d题的Question ID Is %s \n", i+1, questionId)
					   					fmt.Printf("第%d题的Answer ID is %s \n", i+1, answer) */
					finalAnswer0 = answer + "," + finalAnswer0
					switch k {
					case 0:
						fmt.Print("A")
					case 1:
						fmt.Print("B")
					case 2:
						fmt.Print("C")
					case 3:
						fmt.Print("D")
					}
				}
			}
			finalAnswer = strings.TrimRight(finalAnswer0, ",")
			SubmitAnswer(questionId, finalAnswer, paperId, examCountId)
		}
	}
	finshExam(paperId)
	fmt.Print("---------试卷回答完成---------")
}

func SubmitAnswer(questionID, answerID, paperId, examCountId string) {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/SubmitQuestionAnswer2?userId=" + strconv.Itoa(uid) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&courseId=" + strconv.Itoa(courseId) + "&chapterId=0&paperId=" + paperId + "&questionId=" + questionID + "&examTimes=0&examCountId=" + examCountId + "&userAnswers=" + answerID + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(body))

}

func finshExam(paperId string) {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/SubmitPaper?userId=" + strconv.Itoa(uid) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&courseId=" + strconv.Itoa(courseId) + "&chapterId=0&paperId=" + paperId + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Something went wrong!!%s\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print("已完成试卷的提交" + "\n")
	fmt.Print(string(body))

}
