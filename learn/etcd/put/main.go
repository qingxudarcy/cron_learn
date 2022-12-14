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
		putRes *clientv3.PutResponse
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

	if putRes, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hi chang", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("version is %d\n", putRes.Header.Revision)
		if putRes.PrevKv != nil {
			fmt.Printf("prekv is %v\n", string(putRes.PrevKv.Value))
		}
	}

}
