package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr *cronexpr.Expression
		e error
		now time.Time
		nextTime time.Time
	)

	// linux crontab 哪分钟(0-59) 哪小时(0-23) 哪天(1-31) 哪月(1-12) 星期几（0-6）

	// 每分钟执行一次 (支持秒/年 2018-2099)
	if expr, e = cronexpr.Parse("*/5 * * * * * * * *"); e != nil {
		fmt.Println(e)
		return
	}

	// 获取到当前的时间
	now = time.Now()
	// 传入当前的时间可以获取到下次调度的时间
	nextTime = expr.Next(now)

	// 等待这个定时器超时
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了:", nextTime)
	})

	time.Sleep(5 * time.Second)
	fmt.Println(now, nextTime)
}