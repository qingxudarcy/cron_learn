package worker

import (
	"context"
	"cron/project/common"
	"os/exec"
	"time"
)


var (
	G_excuter *Excuter
)

type Excuter struct {

}

func (excuter *Excuter) ExcuteJob(info *common.JobExcuteInfo) {
	go func ()  {
		var (
			cmd *exec.Cmd
			output []byte
			err error
			result *common.JobExcuteResult
			jobLock *JobLock
		)

		result = &common.JobExcuteResult{
			JobExcuteInfo: info,
		}

		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		err = jobLock.TryLock()
		defer jobLock.UnLock()

		result.StartTime = time.Now()

		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
			output, err = cmd.CombinedOutput()
			result.EndTime = time.Now()
			result.Output = string(output)
			result.Err = err
		}

		G_scheduler.PushJobResult(result)
	}()

}

func InitExcuter() (err error) {
	G_excuter = &Excuter{}
	return
}