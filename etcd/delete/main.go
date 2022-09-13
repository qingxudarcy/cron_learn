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
		err error
		delResp *clientv3.DeleteResponse
	)

	config := clientv3.Config{
		Endpoints: []string{"124.222.72.252:12379"},
		DialTimeout: 3 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv := clientv3.NewKV(client)

	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job2", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	} else {
		for _, res := range delResp.PrevKvs {
			fmt.Println("删除了", string(res.Key), string(res.Value))
		}
	}

	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(delResp.Header.String())
	}
}