package common

var (
	JobKeyPrefix string = "/cron/job/"
	JobKillerPrefix string = "/cron/killer/"
	JobLockDir string = "/cron/lock/"
)

var (
	JobDeleteEvent int = 0
	JobPutEvent int = 1
)