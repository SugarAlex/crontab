package worker

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	// 单例
	G_jobMgr *JobMgr
)

func (jobMgr *JobMgr)watchJobs() (err error) {
	
	return
}

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeOut) * time.Millisecond,
	}
	if client, err = clientv3.New(config); err != nil {
		return
	}
	// 得到kv和lease的api子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 给单例赋值
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}