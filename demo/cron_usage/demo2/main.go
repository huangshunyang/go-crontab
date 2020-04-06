package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//代表一个任务
type CronJob struct {
	expr *cronexpr.Expression
	nextTime time.Time // expr.Next(now)
}

func main() {
	// 需要有一个调度的协程，它定时检查所有的Cron任务，谁过期了就执行谁
	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		scheduleTable map[string]*CronJob // 调度表 key:任务的名字
	)

	scheduleTable = make(map[string]*CronJob)
	now = time.Now()

	// 1.我们定义两个CronJob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{expr:expr, nextTime:expr.Next(now)}
	scheduleTable["job1"] = cronJob // 任务注册到调度表

	// 定义第二个任务
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{expr:expr, nextTime:expr.Next(now)}
	scheduleTable["job2"] = cronJob

	// 启动一个调度的协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now time.Time
		)
		// 定时检查一下任务调度表
		for  {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程, 执行这个任务
					go func(jobName string) {
						fmt.Println("执行:", jobName)
					}(jobName)

					// 计算下一次的调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次的调度的时间:", cronJob.nextTime)
				}
			}

			// 睡眠100毫秒
			select {
				case <-time.NewTimer(100 * time.Millisecond).C:
			}
		}
	}()

	time.Sleep(100 * time.Second)
}