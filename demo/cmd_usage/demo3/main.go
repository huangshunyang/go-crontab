package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// 使用一个结构体保存子协程的输出
type result struct {
	e error
	output []byte
}

func main() {
	var (
		cxt context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan *result
		r *result
	)

	// 创建了一个结果队列
	resultChan = make(chan *result, 1000)

	// 执行1个cmd，让它在一个协和里执行2秒，sleep2，echo aaa
	cxt, cancelFunc = context.WithCancel(context.TODO())
	go func() {
		var (
			output []byte
			e error
		)
		cmd = exec.CommandContext(cxt, "F:\\Git\\bin\\bash.exe", "-c", "sleep 2; echo aaa;")
		// 执行任务，捕获输出
		output, e = cmd.CombinedOutput()
		// 将输出的结果，传给main协程
		resultChan <- &result{e: e,output:output}
	}()

	// 继续往下走
	time.Sleep(1 * time.Second)

	// 取消上下文
	cancelFunc()

	// 在main协程里，等待子协程的退出，并打印任务执行的结果
	r = <-resultChan

	//	打印任务执行结果
	fmt.Println(r.e, string(r.output))
}
