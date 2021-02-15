package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Stitch-Zhang/gmp4"
)

//APIServer is the address of study plantform which based on 36ve
const APIServer = "http://125.221.38.2"

// encrypt to 32 size lowcase md5
func encryptPass(raw string) string {
	m := md5.New()
	m.Write([]byte(raw))
	return hex.EncodeToString(m.Sum(nil))
}

func getDuration(url string) string {
	video, err := gmp4.NewRemote(url)
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(gmp4.GetDuration(video)))
}

func autoGrow(u, p string, courseID ...string) {
	cMap := map[string]bool{}
	for _, v := range courseID {
		cMap[v] = true
	}
	user, err := newUser(u, p)
	if err != nil {
		panic(err)
	}
	user.login()
	user.qCourses()
	var (
		times      = 0
		loginTimes = 1
		tTime      = 0
	)
	for {
		user.login()
		loginTimes++
		for _, v := range user.courses {
			if !cMap[v.courseID] {
				continue
			}
			delayTime := rand.Int63n(20)
			time.Sleep(time.Second * time.Duration(delayTime))
			user.montiorWatch(v.courseID, delayTime)
			user.montiorShow(v.courseID, delayTime)
			times++
			tTime += int(delayTime)
		}
		log.Printf("登录次数:+%d \t课程访问量:+%d\t课程总时长:+%ds", loginTimes, times, tTime)
	}
}

func cmd(arg string) {
	cmd := exec.Command("cmd.exe", "/c", arg)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func randomLocalIP() string {
	return fmt.Sprintf("%s%s", "192.168.1.", strconv.Itoa(rand.Intn(255)))
}
