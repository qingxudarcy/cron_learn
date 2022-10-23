package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

type apiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *apiServer
)


func handleJobSave(w http.ResponseWriter, r *http.Request) {

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
		ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout),
		Handler: mux,
	}

	// 赋值单例
	G_apiServer = &apiServer{
		httpServer: httpServer,
	}

	go G_apiServer.httpServer.Serve(listener)

	return

}