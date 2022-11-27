package common

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
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
	JobType int // PUT DELETE
	Job *Job
}

func InitJobEvent(jobType int, job *Job) (*JobEvent) {
	return &JobEvent{
		JobType: jobType,
		Job: job,
	}
}

// 任务调度计划
type JobSchedulerPlan struct {
	Job *Job
	Expr *cronexpr.Expression
	NextTime time.Time
}

func BuildJobSchedulerPlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulerPlan = &JobSchedulerPlan{
		Job: job,
		Expr: expr,
		NextTime: expr.Next(time.Now()),
	}

	return
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