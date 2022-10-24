package master

import (
	"cron/project/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	// "time"
)

type apiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *apiServer
)

// POST job={"name": ""job1m, "command": "echo hello", "cronExpr": "* * * *"}
func handleJobSave(w http.ResponseWriter, req *http.Request) {
	var (
		err error
		job common.Job
		oldJob *common.Job
	)
	

	// 解析json
	if err = json.NewDecoder(req.Body).Decode(&job); err != nil {
		fmt.Print(err)
	}
	
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = common.SuccessRes(w, oldJob); err != nil {
		fmt.Print(err)
	}


	// ERR:
	//   _ = common.ErrRes(w, err.Error())
}


// 初始化服务
func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)
	
	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个http服务
	httpServer = &http.Server{
		// ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		// WriteTimeout: time.Duration(G_config.ApiWriteTimeout),
		Handler: mux,
	}

	// 赋值单例
	G_apiServer = &apiServer{
		httpServer: httpServer,
	}

	go G_apiServer.httpServer.Serve(listener)

	return

}