package master

import (
	"context"
	"cron/project/common"
	"encoding/json"
	"errors"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

const (
	jobKeyPrefix string = "/cron/job/"
)

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)

	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_jobMgr = &JobMgr{
		client: client,
		kv: kv,
		lease: lease,
	}
	return
}

func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var(
		jobKey string
		jobValue []byte
		putRes *clientv3.PutResponse
	)

	jobKey = jobKeyPrefix + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	if putRes, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if err = json.Unmarshal(putRes.PrevKv.Value, &oldJob); err != nil {
		err = nil
		return
	}

	return
}

func (jobMgr *JobMgr)DeleteJob(jobName string) (err error) {
	var (
		jobKey string
		delResp *clientv3.DeleteResponse
	)

	jobKey = jobKeyPrefix + jobName

	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil{
		return
	}
	if len(delResp.PrevKvs) == 0 {
		return errors.New("无法删除不存在的任务")
	}
	return
}

func (jobMgr *JobMgr)ListJob(jobNameSuffix string) (jobList []*common.Job, err error) {
	var (
		getRes *clientv3.GetResponse
	)
	jobName := jobKeyPrefix + jobNameSuffix
	getRes, err = jobMgr.kv.Get(context.TODO(), jobName, clientv3.WithPrefix()); if err != nil {
		return
	}
	jobList = make([]*common.Job, 0, getRes.Count)
	for _, getResKv := range getRes.Kvs {
		var job common.Job
		if err = json.Unmarshal(getResKv.Value, &job); err != nil {
			return
		}
		jobList = append(jobList, &job)
	}
	return
}