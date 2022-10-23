package main

import (
        "fmt"
        "time"

        "github.com/gorhill/cronexpr"
)

func main() {
        var expr *cronexpr.Expression
        var err error

        // 秒 分 时 天 月 星期几 年

        if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
                fmt.Println(err)
        }
        now := time.Now()

        nextTime := expr.Next(now)
        time.AfterFunc(nextTime.Sub(now), func() {
                fmt.Println("被调度了", nextTime)
        })
        time.Sleep(7 * time.Second)
}
