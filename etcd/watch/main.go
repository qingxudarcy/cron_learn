package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	var (
		clinet    *clientv3.Client
		err       error
		kvGetResp *clientv3.GetResponse
	)

	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 3 * time.Second,
	}

	if clinet, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv := clientv3.NewKV(clinet)

	go func() { // 模仿etcd日常操作
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "job7")
			time.Sleep(1 * time.Second)
			kv.Delete(context.TODO(), "/cron/jobs/job7")
		}
	}()

	if kvGetResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	if len(kvGetResp.Kvs) != 0 {
		fmt.Println("当前值为", string(kvGetResp.Kvs[0].Value))
	}

	watchStartRevision := kvGetResp.Header.Revision + 1
	fmt.Printf("从%d版本开始监听\n", watchStartRevision)

	watcher := clientv3.NewWatcher(clinet)
	watchChan := watcher.Watch(context.TODO(), "/cron/jobs/job7", clientv3.WithRev(watchStartRevision)) // 指定监听版本

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为", string(event.Kv.Value), "Revision", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", event.Kv.ModRevision)
			}
		}
	}

}
