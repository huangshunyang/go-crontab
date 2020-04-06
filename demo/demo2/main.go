package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		output []byte
		e error
	)

	cmd = exec.Command("F:\\Git\\bin\\bash.exe", "-c", "sleep 3; ls -l")

	// 执行命令，捕获子进程的输出
	if output, e = cmd.CombinedOutput(); e != nil {
		fmt.Println(e)
		return
	}

	// 子进程的输出
	fmt.Println(string(output))
}
