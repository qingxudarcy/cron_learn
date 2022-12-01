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
		)

		result = &common.JobExcuteResult{
			JobExcuteInfo: info,
			StartTime: time.Now(),
		}
		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
		output, err = cmd.CombinedOutput()
		result.EndTime = time.Now()
		result.Output = string(output)
		result.Err = err

		G_scheduler.PushJobResult(result)
	}()

}

func InitExcuter() (err error) {
	G_excuter = &Excuter{}
	return
}