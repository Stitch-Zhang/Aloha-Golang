package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CourseInfoJSON []struct {
	ID             int    `json:"id"`
	CourseClassID  int    `json:"CourseClassId"`
	CourseID       int    `json:"courseId"`
	ParentCourseID int    `json:"ParentCourseId"`
	CourseName     string `json:"courseName"`
	CourseImage    string `json:"courseImage"`
	IsCharge       int    `json:"IsCharge"`
	Type           int    `json:"type"`
	Credits        int    `json:"Credits"`
	StudyTime      int    `json:"studyTime"`
	OperDate       string `json:"operDate"`
	TypeName       string `json:"typeName"`
}

func main() {
	var (
		uid = 8100
	)
	url := "http://wrggka.whvcse.edu.cn/api/M_Semester/GetLearningCourseList?studentId=" + strconv.Itoa(uid) + "&type=1&pageIndex=1&accessKey=1&secretKey=1&pageSize=10"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var courseInfo CourseInfoJSON
	json.Unmarshal(body, &courseInfo)
	lenth := len(courseInfo)
	for i := 0; i < lenth; i++ {
		fmt.Printf("第%d门课程名称为:%s\t\n", i+1, courseInfo[i].CourseName)
		fmt.Printf("UID :\t%d", courseInfo[i].CourseID)
	}
}
