package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoUrl string = "mongodb://admin:admin@127.0.0.1:27017/?authSource=admin"

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

type FindJobByName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client *mongo.Client
		cond   *FindJobByName
		cur    *mongo.Cursor
		err    error
	)
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl).SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	collection := client.Database("cron").Collection("log")

	cond = &FindJobByName{JobName: "job10"}

	if cur, err = collection.Find(context.TODO(), cond, options.Find().SetSkip(0), options.Find().SetLimit(2)); err != nil {
		fmt.Println(err)
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.TODO()) {
		record := &LogRecord{}
		if err = cur.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}
}
