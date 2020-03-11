package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 关闭计算机
	fmt.Println("关闭主机")
	arg := []string{"-s", "-t", "200000"}
	cmd := exec.Command("shutdown", arg...)
	_, _ = cmd.CombinedOutput()
	return
}
