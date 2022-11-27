package worker

import (
	"cron/project/common"
	"time"
)


type Scheduler struct {
	jobEventChan chan *common.JobEvent
	jobPlanTable map[string]*common.JobSchedulerPlan
}


var (
	G_scheduler *Scheduler
)


func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulerPlan *common.JobSchedulerPlan
		isExisted bool	
		err error
	)

	switch jobEvent.JobType {
	case common.JobPutEvent:
		if jobSchedulerPlan, err = common.BuildJobSchedulerPlan(jobEvent.Job); err != nil{
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan

	case common.JobDeleteEvent:
		if jobSchedulerPlan, isExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; isExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}


func (scheduler *Scheduler) tryScheduler() (schedulerAfter time.Duration) {
	var (
		jobPlan *common.JobSchedulerPlan
		now time.Time
		nearTime *time.Time
	)

	if len(scheduler.jobPlanTable) == 0 {
		schedulerAfter = time.Second
	}

	now = time.Now()

	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			// TODO: 尝试执行任务 前一次任务没结束 则不执行
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		// 统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.After(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	// 下次调度时间 = 最近要执行的任务时间 - 现在的时间
	schedulerAfter = nearTime.Sub(now)

	return
}


// 调度协程
func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent *common.JobEvent
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
	)

	schedulerAfter = scheduler.tryScheduler()
	schedulerTimer = time.NewTimer(schedulerAfter)

	for {
		select {
		case jobEvent = <- scheduler.jobEventChan:
			scheduler.handleJobEvent(jobEvent)
		case <- schedulerTimer.C:  // 最近的任务要执行了
		}

		schedulerAfter = scheduler.tryScheduler()
		schedulerTimer = time.NewTimer(schedulerAfter)
	}
}


func (scheduler *Scheduler) PushJobevent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}


func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulerPlan),
	}
	
	go G_scheduler.schedulerLoop()

	return
}