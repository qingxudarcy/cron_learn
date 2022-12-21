package worker

import (
	"context"
	"cron/project/common"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan *common.JobLog
}

var (
	G_logSink *LogSink
)

func (logSink *LogSink) writeLoop() {
	var jobLog *common.JobLog

	for {
		select {
		case jobLog = <- logSink.logChan:
			logSink.logCollection.InsertOne(context.TODO(), jobLog)
		}
	}
}


func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(G_config.MongoUri).SetConnectTimeout(5 * time.Millisecond)); err != nil {
		return
	}


	G_logSink = &LogSink{
		client: client,
		logCollection: client.Database("cron").Collection("log"),
		logChan: make(chan *common.JobLog, 10000),
	}

	go G_logSink.writeLoop()

	return

}