package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	var (
		u, p, r, o string
		c          []string
	)
	fmt.Print("账号:")
	fmt.Scan(&u)
	fmt.Print("密码:")
	fmt.Scan(&p)
	user, err := newUser(u, p)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		return
	}
	err = user.login()
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		return
	}
	fmt.Println("登录成功")
	time.Sleep(time.Second)
	cmd("cls")
	user.qCourses().displayCourse()
	fmt.Print("输入课程ID,多个以英文的逗号区分(如:2222,3333):")
	fmt.Scan(&r)
	c = strings.Split(r, ",")
	fmt.Printf("A:课件 \t B:挂机\t:")
	fmt.Scan(&o)
	o = strings.ToUpper(o)
	switch o {
	case "A":
		courses := []course{}
		cmd("cls")
		fmt.Println("正在观看课件...")
		csmap := map[string]int{}
		for i, v := range user.courses {
			csmap[v.courseID] = i
		}
		pwg.Add(len(c))
		for index, v := range c {
			courses = append(courses, course{name: user.courses[csmap[v]].name})
			go user.cInfoandHandle(v, index, csmap[v])
		}
		pwg.Wait()
		displayProgress(cps)
		wg.Wait()
		fmt.Println("Everything is done")
		time.Sleep(time.Hour)
	case "B":
		cmd("cls")
		fmt.Println("正在进行挂机...")
		autoGrow(u, p, c...)
	}

}
