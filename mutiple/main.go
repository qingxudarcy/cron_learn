package main

import (
    "fmt"
    "time"

    "github.com/gorhill/cronexpr"
)

type CronJon struct {
    expr *cronexpr.Expression
    nextExcuteTime time.Time
}

func main() {  // 定义一个调度器 执行到期的CronJob
    // 调度表
    scheduleTable := make(map[string]*CronJon)
    
    now := time.Now()
    
    expr1 := cronexpr.MustParse("*/5 * * * * * *")
    cronJob1 := &CronJon{expr: expr1, nextExcuteTime: expr1.Next(now)}
    scheduleTable["job1"] = cronJob1

    expr2 := cronexpr.MustParse("*/40 * * * * * *")
    cronJob2 := &CronJon{expr: expr2, nextExcuteTime: expr2.Next(now)}
    scheduleTable["job2"] = cronJob2

    for {
    now = time.Now()
    for jobName, cronJob := range scheduleTable {
        if cronJob.nextExcuteTime.Before(now) || cronJob.nextExcuteTime.Equal(now) {
            go func (jobName string)  {
                fmt.Printf("%s 任务调度了", jobName)
            }(jobName)
        }
        cronJob.nextExcuteTime = cronJob.expr.Next(now)
    }
    select {
        case <- time.After(500 * time.Millisecond):
      }
  }
}
