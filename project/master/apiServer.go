package master

import (
	"cron/project/common"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

type apiServer struct {
	httpServer *http.Server
}

type JobNameReq struct {
	Name string `json:"name"`
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
		goto ERR
	}
	
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	common.SuccessRes(w, oldJob)

	return


	ERR:
	  common.ErrRes(w, err.Error())
}

func handleJobDelete(w http.ResponseWriter, req *http.Request) {
	var (
		err error
		jobNameReq JobNameReq
	)

	if err = json.NewDecoder(req.Body).Decode(&jobNameReq); err != nil {
		goto ERR
	}

	if err = G_jobMgr.DeleteJob(jobNameReq.Name); err != nil {
		goto ERR
	}

	common.SuccessRes(w, nil)

	return

	ERR:
	  common.ErrRes(w, err.Error())
}

func handleJobList(w http.ResponseWriter, req *http.Request) {
	var (
		err error
		jobNameReq JobNameReq
		jobList []*common.Job
	)
	
	if err = json.NewDecoder(req.Body).Decode(&jobNameReq); err != nil {
		goto ERR
	}

	if jobList, err = G_jobMgr.ListJob(jobNameReq.Name); err != nil {
		goto ERR
	}

	common.SuccessRes(w, jobList)

	return

	ERR:
	  common.ErrRes(w, err.Error())
}

func handleJobKill(w http.ResponseWriter, req *http.Request) {
	var (
		err error
		jobNameReq JobNameReq
	)
	

	if err = json.NewDecoder(req.Body).Decode(&jobNameReq); err != nil {
		goto ERR
	}

	if err = G_jobMgr.KillJob(jobNameReq.Name); err != nil {
		goto ERR
	}

	common.SuccessRes(w, nil)

	return

	ERR:
	   common.ErrRes(w, err.Error())
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
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)


	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个http服务
	httpServer = &http.Server{
		ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler: mux,
	}

	// 赋值单例
	G_apiServer = &apiServer{
		httpServer: httpServer,
	}

	go G_apiServer.httpServer.Serve(listener)

	return

}