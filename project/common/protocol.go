package common

import (
	"encoding/json"
	"net/http"
)

type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	ErrNo int  `json:"errNo"`
	Msg string `json:"msg"`
	Data interface{}  `json:"data"`
}

func SuccessRes(w http.ResponseWriter, data interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	resp := &Response{
		ErrNo: 0,
		Msg: "success",
		Data: data,
	}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		return
	}

	return
}

func ErrRes(w http.ResponseWriter, errMsg string) (err error) {
	resp := &Response{
		ErrNo: -1,
		Msg: errMsg,
		Data: nil,
	}

	jsonRes, _ :=  json.Marshal(resp)
	w.Write(jsonRes)
	w.Header().Set("Content-Type", "application/json")

	return
}