package main

import (
	"cron/project/master"
	"flag"
	"fmt"

	"runtime"
)

var (
	confFile string
)

func initArgs() {
	// master -config ./master.json
	//master -h
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 初始化线程
	initEnv()
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}


	return

ERR:
	fmt.Println(err)
}