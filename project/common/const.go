package common

var (
	JobKeyPrefix string = "/cron/job/"
	JobKillerPrefix string = "/cron/killer/"
)

var (
	JobDeleteEvent int = 0
	JobPutEvent int = 1
)