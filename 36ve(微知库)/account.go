package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/tidwall/gjson"
)

type (
	info struct {
		itemID       string
		nID          string
		link         string
		fID          string
		resourceType string
		done         bool
	}
	course struct {
		courseID string
		name     string
		teacher  string
		done     bool
		infos    []info
	}
	account struct {
		username string
		password string
		token    string
		ip       string
		courses  []course
		client   *http.Client
	}
)

var wg sync.WaitGroup

func (a account) do(method, url string, body io.Reader) string {
	req, _ := http.NewRequest(method, url, body)
	a.client.Timeout = 3 * time.Second
	req.Header.Set("User-Agent", "okhttp/3.3.0")
	resp, err := a.client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	rBody, _ := ioutil.ReadAll(resp.Body)
	return string(rBody)
}

func (a *account) login() error {
	url := fmt.Sprintf("http://uc.36ve.com/xiaomi/index.php/user/login?device_type=0&username=%s&password=%s", a.username, encryptPass(a.password))
	resp := a.do("GET", url, nil)
	a.token = gjson.Get(resp, "token").String()
	if a.token == "" {
		return errors.New("登陆失败:账号或密码有误")
	}

	return nil
}

func (a *account) qCourses() *account {
	if a.token == "" {
		err := a.login()
		if err != nil {
			log.Panic(err)
			return nil
		}
	}
	url := fmt.Sprintf("%s/?q=get_courses_list_new&student_status=5&username=%s&token=%s", APIServer, a.username, a.token)
	resp := a.do("GET", url, nil)
	courses := gjson.Get(resp, "data")
	for _, v := range courses.Array() {
		a.courses = append(a.courses, course{courseID: v.Get("course_id").String(), name: v.Get("name").String(), teacher: v.Get("author").String()})
	}
	return a
}

func (a *account) cInfoandHandle(courseID string, i, index int) *account {
	url := fmt.Sprintf("%s/?q=get_course_module_items2/%s&username=%s", APIServer, courseID, a.username)
	resp := a.do("GET", url, nil)
	infos := gjson.Parse(resp)
	addProgress(a.courses[index].name, infos)
	pwg.Done()
	for _, k := range infos.Array() {
		for m, j := range k.Get("items").Array() {
			var finished bool
			if j.Get("progress").String() == "10" {
				finished = true
				pwg.Wait()
				cps[i].c <- 1
				//fmt.Printf("course_id=%s item_id=%s is done\n", courseID, j.Get("item_id").String())
			}
			a.courses[index].infos = append(a.courses[index].infos,
				info{itemID: j.Get("item_id").String(),
					nID:          j.Get("nid").String(),
					resourceType: j.Get("type").String(),
					fID:          j.Get("fid").String(),
					link:         a.getLink(j.Get("fid").String(), j.Get("item_id").String()),
					done:         finished,
				})
			if !finished {
				wg.Add(1)
				go func(inf info, i int) {
					defer wg.Done()
					time.Sleep(time.Second * time.Duration(rand.Int63n(5)))
					a.saveRecord(inf)
					a.montiorWatch(courseID, rand.Int63n(20))
					pwg.Wait()
					cps[i].c <- 1
				}(a.courses[index].infos[m], i)
			}
		}
	}
	wg.Wait()
	return a
}

//Spending a long time to get data
//Avoid to using this insted cInfoandHandle
func (a *account) getAllCourseInfos() *account {

	wg.Add(len(a.courses))
	for i, v := range a.courses {
		go func(index int, value course) {
			defer wg.Done()
			url := fmt.Sprintf("%s/?q=get_course_module_items2/%s&username=%s", APIServer, value.courseID, a.username)
			resp := a.do("GET", url, nil)
			infos := gjson.Parse(resp)
			for _, k := range infos.Array() {
				for _, j := range k.Get("items").Array() {
					//Handle resource is watched
					var finished bool
					if j.Get("progress").String() == "10" {
						finished = true
						fmt.Printf("course_id=%s item_id=%s is done\n", value.courseID, j.Get("item_id").String())
					}
					a.courses[index].infos = append(a.courses[index].infos,
						info{itemID: j.Get("item_id").String(),
							nID:          j.Get("nid").String(),
							resourceType: j.Get("type").String(),
							fID:          j.Get("fid").String(),
							link:         a.getLink(j.Get("fid").String(), j.Get("item_id").String()),
							done:         finished,
						})
				}
			}

		}(i, v)
	}
	wg.Wait()
	return a
}

func (a account) getLink(fID, itemID string) string {
	url := fmt.Sprintf("%s/?q=get_file_url/%s&resource=1&item_id=%s&username=%s", APIServer, fID, itemID, a.username)
	resp := a.do("GET", url, nil)
	return gjson.Get(resp, "url").String()
}

func (a *account) saveRecord(info info) {
	if info.resourceType != "video_material" {
		body := url.Values{
			"username": {a.username},
			"itemid":   {info.itemID},
			"nid":      {info.nID},
		}
		http.PostForm(APIServer+"/?q=save_user_item_progress_by_app", body)
		//fmt.Println("PPT Sumbit=", info.itemID)
		return
	}
	//fmt.Println("Got:", info)
	totalTime := getDuration(info.link)
	body := url.Values{}
	data := map[string]string{
		"time":      totalTime,
		"totaltime": totalTime,
		"nid":       info.nID,
		"fid":       info.fID,
		"item_id":   info.itemID,
		"resource":  "1",
		"username":  a.username,
		"token":     a.token,
	}
	for k, v := range data {
		body.Set(k, v)
	}
	http.PostForm(APIServer+"/?q=items/study/current/save", body)
	//rbody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("Save responsed:", string(rbody))
	cps[1].c <- 1
}

func (a *account) handleACourse(courseID string) {
	for indexC, v := range a.courses {
		if v.courseID == courseID {
			if len(v.infos) == 0 {
				fmt.Printf("courseID:%s 's infos ungotten\n", courseID)
				return
			}
			for _, inf := range v.infos {
				if inf.done {
					return
				}
				t := info{
					nID:          inf.nID,
					fID:          inf.fID,
					itemID:       inf.itemID,
					link:         inf.link,
					resourceType: inf.resourceType,
				}
				a.saveRecord(t)
			}
		}
		a.courses[indexC].done = true
	}
}

func (a *account) handleCourses(cIDs ...string) {
	for _, cID := range cIDs {
		a.handleACourse(cID)
	}
}

func (a *account) montiorWatch(courseID string, delayTime int64) {
	url := fmt.Sprintf("%s/api.php?action=monitor&username=%s&t=%d&class=M4_CoursesDes&ip=%s&objectid=&courseid=%s&equipment=0", APIServer, a.username, delayTime, a.ip, courseID)
	a.do("POST", url, nil)
}

func (a *account) montiorShow(courseID string, delayTime int64) {
	url := fmt.Sprintf("%s/api.php?action=monitor&username=%s&t=%d&class=ShowResAct&ip=%s&objectid=&courseid=%s&equipment=0", APIServer, a.username, delayTime, a.ip, courseID)
	a.do("POST", url, nil)
}

func (a *account) displayCourse() {
	t := table.NewWriter()
	t.SetStyle(table.StyleDouble)
	t.AppendHeader(table.Row{"#", "课程名", "课程创建教师", "ID"})
	for index, v := range a.courses {
		t.AppendRow(table.Row{index, v.name, v.teacher, v.courseID})
	}
	fmt.Println(t.Render())
}

func newUser(name, pass string) (*account, error) {
	if name == "" || pass == "" {
		return nil, errors.New("name or password is null")
	}
	rand.Seed(time.Now().Unix())
	return &account{
		username: name,
		password: pass,
		courses:  []course{},
		client:   new(http.Client),
		ip:       randomLocalIP(),
	}, nil
}
