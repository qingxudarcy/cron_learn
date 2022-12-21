package main

import (
	"cron/project/worker"
	"flag"
	"fmt"

	"runtime"
)

var (
	confFile string
)

func initArgs() {
	// worker -config ./master.json
	//worker -h
	flag.StringVar(&confFile, "config", "./worker.json", "指定master.json")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	// 初始化命令行参数
	initArgs()

	// 初始化配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 初始化线程
	initEnv()

	// 启动日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	worker.InitExcuter()

	worker.InitScheduler()

	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}


	for {
		select{}
	}

ERR:
	fmt.Println(err)
}