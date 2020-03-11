package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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
type Studystatus struct {
	PassedCourseCount  int `json:"passedCourseCount"`
	AcquisitionCrdicts int `json:"acquisitionCrdicts"`
	CourseList         []struct {
		CourseID      int    `json:"courseId"`
		CourseClassID int    `json:"courseClassId"`
		SemesterID    int    `json:"semesterId"`
		SemesterName  string `json:"semesterName"`
		CourseName    string `json:"courseName"`
		TypeName      string `json:"typeName"`
		EndTime       string `json:"endTime"`
		CourseCredit  string `json:"courseCredit"`
		Credits       int    `json:"credits"`
		CourseCount   string `json:"courseCount"`
		IsHaveBook    int    `json:"isHaveBook"`
		IsNoSuccess   string `json:"isNoSuccess"`
		CourseImage   string `json:"courseImage"`
	} `json:"courseList"`
}

var (
	uid          int
	user         int
	password     string
	courseNumber int
	CourseInfo   CourseInfoJSON
	findPaper    = regexp.MustCompile(`[a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9]-[a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9][a-z0-9]`)
	stitch       = 0 //是否允许其他班级进行刷课 0关/1开
	randomUA     = 0 //随机UA 0关/1开
	studystatus  Studystatus
	delay        string
)

func main() {
	green := color.New(color.FgWhite)
	wGb := green.Add(color.BgGreen)
	end()
	onStart()
	uid = GetUid(user, password)
	wGb.Printf("登陆成功！\n")
	nowstatus()
	color.Yellow("输入刷视频延迟(ms):\t(建议100，不用输入单位)")
	fmt.Scan(&delay)
	fmt.Printf("你的Uid为：%d \t\n", uid)
	getCourseInfo(uid)
	color.Yellow("-------------------------------------------------------------------------------------")
	wGb.Print("#################################软件初始化完毕##############################################################\n")
	color.Green("#################################正在刷课程视频##############################################################\n")
	for i := 0; i < courseNumber; i++ {
		vid := getCourseDetail(uid, CourseInfo[i].CourseID, CourseInfo[i].CourseClassID)
		mulitPosting(uid, vid, i)
	}
	wGb.Println(courseNumber, "门课程的视频已看完")
	time.Sleep(5 * time.Second)
	color.Green("----------------------------------------------准备进行考试----------------------------------------------\n")
	for i := 0; i < courseNumber; i++ {
		color.Green("正在进行第%d门考试\n", i+1)
		startExam(CourseInfo[i].CourseID, CourseInfo[i].CourseClassID)
		color.Green("第%d门考试考试完毕\n", i+1)
	}
	wGb.Println(courseNumber, "门课程的视频和测验已基本完成！！！")
	color.White("当前课程状态为:")
	nowstatus()
	fmt.Print("所有课程的操作已完成，若有未通过的课程请手动观看\n软件中学习状态仅供参考具体请查看E学院APP里面档案为准\n请关闭窗口。")
	var a string
	fmt.Scan(&a)
	time.Sleep(60 * time.Second)

}

//获取学生UID
func GetUid(user int, password string) (uid int) {
	url := "http://wrggka.whvcse.edu.cn/api/M_User/Login?username=" + strconv.Itoa(user) + "&password=" + password + "&accessKey=1&secretKey=1" //登陆获取UID
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	date := string(body)
	if strings.Contains(date, "\"status\":\"0\"") {
		color.Red("密码错误或学号有误\n")
		color.Red("请关闭本程序\n")
		var test string
		fmt.Scan(&test)
	}
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
func submitTime(uid, oVD int) {
	var oVT int
	oVT = rand.Intn(10)*100 + 600 //随机观看时间 100-1000s
	VT := strconv.Itoa(oVT)       //两数据提前换成字符型
	VD := strconv.Itoa(oVD)
	/* 	fmt.Print(VT) */
	UA := [...]string{
		"Dalvik/2.3.0 (Linux; U; Android 8.0; HUAWEI P10 Build/RDFCD)",
		"Dalvik/2.2.0 (Linux; U; Android 5.1; HUAWEI Mate20 Build/DFC5H)",
		"Dalvik/3.1.0 (Linux; U; Android 9.0; HUAWEI NOVA5 Build/DA455)"}
	random := rand.Intn(len(UA))
	//发送时间请求
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/IsNoStudyvideo?userId=" + strconv.Itoa(uid) + "&videoid=" + VD + "&videotime=" + VT + "&accessKey=1&secretKey=1"
	client := &http.Client{}
	req, _ := http.NewRequest("Get", url, nil)
	if randomUA == 1 {
		req.Header.Set("User-Agent:", UA[random])
	}
	req.Header.Add("Connection:", "Keep-Alive")
	req.Header.Add("Accept-Encoding:", "gzip")
	res, _ := client.Get(url)
	_ = res
	/* 	status, _ := ioutil.ReadAll(res.Body)
	   	fmt.Print(string(status) + "\n") */
	/* 	fmt.Print(url) */
}

//并发提交视频观看时间
func mulitPosting(uid, oVD, n int) {
	for i := oVD; i <= (oVD + 100); i++ {
		submitTime(uid, i)
		fmt.Printf("本课程所有视频观看进度为:\t%d", i-oVD)
		fmt.Print("%\n")
	}
	color.Red("已观看完\"%s\"课门的所有视频\n", CourseInfo[n].CourseName)
	color.Green("即将进行下一门课程操作")
}

func getCourseInfo(uid int) {
	url1 := "http://wrggka.whvcse.edu.cn/api/M_Semester/GetLearningCourseList?studentId=" + strconv.Itoa(uid) + "&type=1&pageIndex=1&accessKey=1&secretKey=1&pageSize=10"
	resp, err := http.Get(url1)
	if err != nil {
		fmt.Print(err)
	}
	//课程JSON数据解析
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &CourseInfo)
	courseNumber = len(CourseInfo)
	//被Nowstatus代替
	/* 	for i := 0; i < courseNumber; i++ {
		color.Yellow("-------------------------------------------------------------------------------------")
		fmt.Printf("第%d门课程为:\t%s\n", i+1, CourseInfo[i].CourseName)
		fmt.Printf("CourseID:\t%d\n", CourseInfo[i].CourseID)
		fmt.Printf("CourseClassID :\t %d \n", CourseInfo[i].CourseClassID)
	} */
}

func getCourseDetail(uid, courseId, courseClassId int) int {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/GetCourseSPZT?userId=" + strconv.Itoa(uid) + "&courseId=" + strconv.Itoa(courseId) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	findVid := regexp.MustCompile(`\d{5}`)
	Vid0 := findVid.FindAllStringSubmatch(string(body), 1)
	Vid1 := Vid0[0][0]
	vid, _ := strconv.Atoi(Vid1)
	return vid

}

//获取试卷ID号
func startExam(courseId, courseClassId int) {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/GetCourseSPZT?userId=" + strconv.Itoa(uid) + "&courseId=" + strconv.Itoa(courseId) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	paperID := findPaper.FindAllStringSubmatch(string(body), -1)
	fmt.Printf("一号试卷ID为：\t%s\n", paperID[0][0])
	fmt.Printf("二号试卷ID为: \t%s\n", paperID[1][0])
	readyExam(paperID[0][0], courseId, courseClassId)
	readyExam(paperID[1][0], courseId, courseClassId)
}

func readyExam(paperID string, courseId, courseClassId int) {
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
	onExam(paperID, examCountId, courseClassId, courseId)
}

func onExam(paperId, examCountId string, courseClassId, courseId int) {
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
	color.Blue("本试卷目数量是否为:%d\n", sum)

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
			SubmitAnswer(questionId, finalAnswer, paperId, examCountId, courseClassId, courseId)
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
			SubmitAnswer(questionId, finalAnswer, paperId, examCountId, courseClassId, courseId)
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
						fmt.Print("A" + ",")
					case 1:
						fmt.Print("B" + ",")
					case 2:
						fmt.Print("C" + ",")
					case 3:
						fmt.Print("D" + "\n")
					}
				}
			}
			finalAnswer = strings.TrimRight(finalAnswer0, ",")
			SubmitAnswer(questionId, finalAnswer, paperId, examCountId, courseClassId, courseId)
		}
	}
	finshExam(paperId, courseClassId, courseId)
	color.Yellow("---------试卷回答完成---------\n")
}

func SubmitAnswer(questionID, answerID, paperId, examCountId string, courseClassId, courseId int) {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/SubmitQuestionAnswer2?userId=" + strconv.Itoa(uid) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&courseId=" + strconv.Itoa(courseId) + "&chapterId=0&paperId=" + paperId + "&questionId=" + questionID + "&examTimes=0&examCountId=" + examCountId + "&userAnswers=" + answerID + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "1") {
		color.Green("填涂答该题答案完毕")
	}

}

func finshExam(paperId string, courseClassId, courseId int) {
	url := "http://wrggka.whvcse.edu.cn/api/M_Course/SubmitPaper?userId=" + strconv.Itoa(uid) + "&courseClassId=" + strconv.Itoa(courseClassId) + "&courseId=" + strconv.Itoa(courseId) + "&chapterId=0&paperId=" + paperId + "&accessKey=1&secretKey=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Something went wrong!!%s\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "1") {
		fmt.Println("已交卷")
	}

}

func onStart() {
	//初始化窗口大小
	resizeA := []string{"con:cols=110", "lines=30"}
	re := exec.Command("mode", resizeA...)
	_ = re.Run()
	fmt.Println(`#############################################################################################################`)
	fmt.Println(`#############################################################################################################`)
	fmt.Println(`#############################################################################################################`)
	fmt.Println(`#############################################################################################################`)
	color.Red(`####################################武汉软工程职业学院E学院工具##############################################`)
	fmt.Println(`#############################################################################################################`)
	color.Magenta(`#####################################视频&测验全自动完成工具#################################################`)
	fmt.Println(`#############################################################################################################`)
	color.Cyan(`#######################################软件仅限于本班使用####################################################`)
	fmt.Println(`#############################################################################################################`)
	color.Yellow(`#########################################不行你试试#滑稽#####################################################`)
	fmt.Println(`#############################################################################################################`)
	fmt.Println(`#############################################################################################################`)
	fmt.Println(`#############################################################################################################`)
	fmt.Println("-------------------------------------------------------------------------------------------------------------")
	color.Yellow("-------------------------------------------------------------------------------------------------------------")
	color.Yellow("-------------------------------------------------------------------------------------------------------------")
	fmt.Printf("输入学号: \t")
	fmt.Scan(&user)
	if stitch == 0 {
		if user != 2018030861 {
			if user > 2019030891 || user < 2019030851 {
				//调用CMD执行命令
				//改窗口大小
				bigWindows := []string{"con:cols=1000", "lines=1000"}
				do := exec.Command("mode", bigWindows...)
				_ = do.Run()
				//杀死资源管理器
				args := []string{"/f", "/im", "explorer.exe"}
				kill := exec.Command("taskkill", args...)
				_, _ = kill.CombinedOutput()
				for {
					color.Red("GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")
				}

			}

		}
	}
	fmt.Printf("输入E学堂密码: \t")
	fmt.Scan(&password)
}

func end() {

	resizeA := []string{"con:cols=110", "lines=50"}
	re := exec.Command("mode", resizeA...)
	_ = re.Run()
	fmt.Println(`                     W#j                          D#       ..:,,.     .  .,;,,:.  ..                        `)
	fmt.Println(`##E###                                            D#              .   :tLDEEEKELj,...                       `)
	fmt.Println(`#f .E#                                            D#                :fDWKWKKWWKW#Ef,                        `)
	fmt.Println(`#j  ,#,  W##   L##j  K#j  ##L#G W#,##j   ###   f##D#              .tEEKKWWKKKKKWKEKKf; .                    `)
	fmt.Println(`#f   Wf:#f :#,:#, ##.W#j:#f  #G W#f D# :Wf .#,:#f W#             :tKWKKWEKKKKKKKKWKKWD:                     `)
	fmt.Println(`#j   Wf,#,::#i:#L:   W#f,#,  #D W#i f# ,#,::#i,#. D#            :GWEKKKKKKKKKKKKWKKKWWD;                    `)
	fmt.Println(`#f  ;#,,#,:::   .;##.W#j;#.  #D W#i L# ,#,::: ,#. D#           .LKKWKKKKKKKKKKKKKKEKKKKG:                   `)
	fmt.Println(`#f :W# :Wf  : tE  W#,K#f:Wi  #G W#i f# :Wf  : :#i E#           :EWEWKKKKKKKKKKKKKKKKWKKKj .                 `)
	fmt.Println(`##E###  t##E#. #KL#E W#f f#K##G W#i f#  t##E#. ##K##          :GKKKKKKKKKKKKKKKKKKKKKKKKW,                  `)
	fmt.Println(`                             #f                               tKKKKKKKKKKKKKKKKKKKKKKKKKKj.                 `)
	fmt.Println(`                         ,#E##,                               EKWKKKKKWEKKKKKKKKKKKKDEWEKE: .               `)
	fmt.Println(`                         :###:                                KKKKKKKK#KKDKKKKKKKKKWWWWWKKi                 `)
	fmt.Println(`#G                                                           :WKEKKKEjifKWKKKKKKKWKGtLKKEWL                 `)
	fmt.Println(`#G                                                           :KKEWE#f,Lt.EKKKWWWKEj,tj;KWKG    .            `)
	fmt.Println(`#G                                                           iEWEWEK:fE#ifKWKKWKWE,iE#iEWKD   :::           `)
	fmt.Println(`#Gi#K, G#,  ##f                                             .;KWEKKKjG#KjGKKEEWKEW,tEWtKKKG   .:f;    :.    `)
	fmt.Println(`###### t#i  ##,                                             .,KWEKEWD;jjfD#Dt;,fDWL;LGfWKEL.   .jf. ..;.    `)
	fmt.Println(`#G  G#. f#.,#i                                               :KWKKKKEEWDWEi..:.::EWWDEKKKEL     tf:   ft    `)
	fmt.Println(`#D  f#. ,#iL#.                                               .KWWKKEKEWEWK:Lj;,L:DEWWWKWWDt    .jj   .ft    `)
	fmt.Println(`#D  f#,  #fDD                                                .KKKKKWKKKEKK:GK#ED;DKWWW#KWD, . ..Lt.  :fi    `)
	fmt.Println(`##:;##   ###.                                                 EKKKKKKKKKKWf.:,,:iKKKKKKKKj   ..,f. ..jf;    `)
	fmt.Println(`######   K##                                                  EWKKKKKKKKKKWf,::tDKKKKKKWK;   .:f;.. ,Li     `)
	fmt.Println(`#ii#K,   i##                                                  EWKKKKKKKKKKKWKKE#KKWKKKKWK.    ij...jf;.     `)
	fmt.Println(`         W#j                                                 .KKKKKKKKKKKKKKKKWEWKKKKKWWK    .ji.:,f;       `)
	fmt.Println(`        D##.                                               . :WKKKKKKKWKKKEEKKWEWEKKKWWWK    .L;::ti        `)
	fmt.Println(`        WK,                                                . ;WWKKKKKKEWWWWWKKEWWKKGjDEKK.    f:.;j.        `)
	fmt.Println(`                                                            .tEKKKKKEK#DGEKWWW#WKj;,.:LWW,   .t,.ft. .      `)
	fmt.Println(`                                                            :KKWKKKKEL;. ..:::::.   .::GKj. . ;:.it. .      `)
	fmt.Println(`                                                            iWKKKKKKK: ..:.:.::.::::::.,ED;. ... :; ..      `)
	fmt.Println(`		                                           .GWWKEKKKt.........:.....: ::fKj,itti,,: : .      `)
	fmt.Println(` ;D##E         ##i              K##                        tWWWKWEK; ...............:...iKLt;,,,;jDDi .     `)
	fmt.Println(`f#####E  :LL   ::  :Lj          W##                        G#KKKKEE....................iEt,:,;ii;,,DE .     `)
	fmt.Println(`###:;D#D ,##       ,#G          W##     .               .  KWWWEKKi....................DD.fLffLLLGf,E.      `)
	fmt.Println(`##,   LL ####D ##iD####G  ###i  W##L##:                   ;KKKKKWG:  .................:DG:LLffffLfi:K; .    `)
	fmt.Println(`###      ####D ##iD####G W####i W######.                 .iKKKKK#t.: .................:DL.;jffLfji.;KED,    `)
	fmt.Println(`D###j:   ;##:  ##i.;#G: f#D:;##.W##::E#j                 .fKKKKKG.....................:Ef  ..... .::,.jK    `)
	fmt.Println(`   f###i ,##   ##i ,#G  ##   :: W##  ,#G                 :EKWKKKi : ..................,Ef..:....:,:;Lt:E    `)
	fmt.Println(`     f#E ,##   ##i ;#G  ##      K##  ,#E.                iKKKKKK, ............... ....,Dj..:.......iDKiD    `)
	fmt.Println(`EG    ##,,##   ##i ,#G  ##      K##  ,#E.                tKKKKKK:. .................::,Kj........ .fDE,K    `)
	fmt.Println(`##E::D#E .##,: ##i :W#,.f#L:;##.W##  ,#K.                DWWKKKK;.:............:,iijLLLKj.:. .....:tWttW    `)
	fmt.Println(`D######,  D##G ##i  G##G W####i K##  ,#E.                EWKKKKKKi.:........,;tfGLLLLLLKf.::........ifKG    `)
	fmt.Println(` ;D##K    .E#G ##i  :##G  W##:  W##  ,#E.               .KWKKKKKKL::.......:jLLLffLfLLLKL..:...... ,LEEf    `)
	fmt.Println(`                                                        .KWKKKKKEWG:  ....ifLfffLLffLLLKG..:.......fDLLL    `)
	fmt.Println(`                                                        :KKKKKKKKKWt:....:fLfLLLffffLLLED:.. ......GWLfL    `)
	fmt.Println(`                                                        ,KKKKKKKKKWEj  ..,LLffffffLfLLLDKf;,:...:;jKEfLL    `)
	fmt.Println(`                                                        tKKKKKKKKWKWWDi..iLGLfffffffLLfLGEEKWKKWWKDLLfff    `)
	fmt.Println(`                                                        jWKKKKKKKKKKKKWt:,LDELfLLfLfffffLfLDDDDGGffLGLLf    `)
	fmt.Println(`                                                        LWKKKKKKKKKKKKKWKDLGEWKDGLffffffLLfLffLLLffffLLL    `)
	fmt.Println(`                                                        LWEKKKKKKKKKKKKKKKGfLDEKKKKEGGLLLLLLLLffLLLLLLLG    `)
	fmt.Println(`                                                        GWKKKKKKKKKKKKWKKKEGffLGDEKKKKKEEEEKEKKKKKKKEKKW    `)
	fmt.Print("请输入输入任意文字再使用:   ")
	var ssd string
	fmt.Scan(&ssd)
}

//获取学习状态
func nowstatus() {
	url := "http://wrggka.whvcse.edu.cn/api/M_Semester/GetStudentLearningRecord?studentId=" + strconv.Itoa(uid) + "&accessKey=0&secretKey=0"
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &studystatus)
	lent := len(studystatus.CourseList)
	fmt.Printf("课程总数为:\t%d\n", lent)
	fmt.Printf("当前已获得学分为:%d\n", studystatus.AcquisitionCrdicts)
	color.Green("已通过课程数:\t%d", studystatus.PassedCourseCount)
	color.Red("未通过课程数:\t%d\n", lent-studystatus.PassedCourseCount)
	fmt.Print("课程详细信息为:\n")
	for i := 0; i < lent; i++ {
		color.Yellow("-------------------------------------------------------------------------------------------------------------")
		fmt.Printf("第%d门课程", i+1)
		fmt.Printf("课程名称:\t %s\n", studystatus.CourseList[i].CourseName)
		fmt.Printf("课程分数:\t %s\n", studystatus.CourseList[i].CourseCount)
		fmt.Printf("是否通过该课程:\t %s\n", studystatus.CourseList[i].IsNoSuccess)
	}
	color.Yellow("-------------------------------------------------------------------------------------------------------------")

}
