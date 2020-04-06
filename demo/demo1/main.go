package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		e error
	)

	cmd = exec.Command("F:\\Git\\bin\\bash.exe", "-c", "echo 1")
	e = cmd.Run()
	fmt.Println(e)
}
