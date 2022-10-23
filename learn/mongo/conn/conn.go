package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var (
		client *mongo.Client
		err    error
	)

	uriCO := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	timeoutCO := options.Client().SetConnectTimeout(5 * time.Second)
	client, err = mongo.Connect(context.TODO(), uriCO, timeoutCO)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = client.Disconnect(nil); err != nil {
			panic(err)
		}
	}()

	database := client.Database("my_db")
	database.Collection("my_collection")

}
