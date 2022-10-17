package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func main() {
	var (
		client     *mongo.Client
		result     *mongo.InsertOneResult
		manyResult *mongo.InsertManyResult
		err        error
	)
	clientOP := options.Client().ApplyURI(mongoUrl)
	clientOP = clientOP.SetTimeout(10 * time.Second)
	if client, err = mongo.Connect(context.TODO(), clientOP); err != nil {
		fmt.Printf("new client err is %v", err)
		return
	}

	defer func() {
		// 延迟释放连接
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("cron").Collection("log")

	record := LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Printf("insert err is %v", err)
		return
	}

	fmt.Println(result.InsertedID.(primitive.ObjectID).Hex())

	logArr := []interface{}{record, record, record}

	if manyResult, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(manyResult.InsertedIDs...)
}
