package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	var (
		client *clientv3.Client
		err    error
		getRes *clientv3.GetResponse
	)

	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 3 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv := clientv3.NewKV(client)

	if getRes, err = kv.Get(context.TODO(), "/cron/jobs/job1", clientv3.WithCountOnly()); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getRes.Count)
	}

	if getRes, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getRes.Kvs)
	}
}
