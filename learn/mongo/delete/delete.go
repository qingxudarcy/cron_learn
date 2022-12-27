package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoUrl string = "mongodb://admin:admin123@127.0.0.1:27017/?authSource=admin"

type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client    *mongo.Client
		deleteRes *mongo.DeleteResult
		err       error
	)
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl).SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		// 延迟释放连接
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	clollection := client.Database("cron").Collection("log")

	deleteCond := &DeleteCond{BeforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	if deleteRes, err = clollection.DeleteMany(context.TODO(), deleteCond); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("删除了%v条记录", deleteRes.DeletedCount)
}