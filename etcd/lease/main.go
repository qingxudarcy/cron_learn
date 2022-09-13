package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {   // 无限续租约
	var (
		client *clientv3.Client
		err error
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		leaseResp *clientv3.LeaseGrantResponse
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
	)

	config := clientv3.Config{
		Endpoints: []string{"124.222.72.252:12379"},
		DialTimeout: 3 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
   
	lease := clientv3.NewLease(client)

	if leaseResp, err = lease.Grant(context.TODO(), 10); err != nil {   // 设置租约时间
		fmt.Println(err)
		return
	}
	leaseId := leaseResp.ID

	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {   // 获取续租请求chan
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp := <- keepRespChan:
				if keepResp != nil {
					fmt.Println("成功续租", keepResp.ID)
				} else {
					fmt.Println("租约失效了")
					goto END
				}
			case <- time.NewTimer(5 * time.Second).C:
				fmt.Println("续租超时")
				goto END
			}
		}
		END:
	}()

	kv := clientv3.NewKV(client)

	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "first time", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("写入成功", putResp.Header.Revision)
	}
	for {
		if getResp, err = client.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}	
		if getResp.Count != 0 {
			fmt.Println("还没有删除", getResp.Kvs)
		} else {
			fmt.Println("已经被删除了")
			break
		}
		time.Sleep(2 * time.Second)	
	}
}