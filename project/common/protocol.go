package common

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

func UnpackJob(value []byte) (*Job, error) {
	var (
		job *Job
		err error
	)

	job = &Job{}
	err = json.Unmarshal(value, job)
	return job, err
}

func ExtractJobName(key string) (string) {
	return strings.TrimPrefix(key, JobKeyPrefix)
}

type JobEvent struct {
	JonType int // PUT DELETE
	Job *Job
}

func InitJobEvent(jobType int, job *Job) (*JobEvent) {
	return &JobEvent{
		JonType: jobType,
		Job: job,
	}
} 

type Response struct {
	ErrNo int  `json:"errNo"`
	Msg string `json:"msg"`
	Data interface{}  `json:"data"`
}

func SuccessRes(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := &Response{
		ErrNo: 0,
		Msg: "success",
		Data: data,
	}
	json.NewEncoder(w).Encode(resp)
}

func ErrRes(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := &Response{
		ErrNo: -1,
		Msg: errMsg,
		Data: nil,
	}

	jsonRes, _ :=  json.Marshal(resp)
	w.Write(jsonRes)
}