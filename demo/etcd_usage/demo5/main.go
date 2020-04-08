package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		e error
		v *mvccpb.KeyValue
		client *clientv3.Client
		kv clientv3.KV
		delResponse *clientv3.DeleteResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"192.168.45.202:2379"},
		DialTimeout:5 * time.Second,
	}

	// 建立一个客户端
	if client, e = clientv3.New(config); e != nil {
		fmt.Println(e)
		return
	}

	fmt.Println(client)

	// KV用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	if delResponse, e = kv.Delete(context.TODO(),"/cron/jobs/job2", clientv3.WithPrevKV()); e != nil {
		fmt.Println(e)
		return
	}

	if len(delResponse.PrevKvs) > 0 {
		for _, v = range delResponse.PrevKvs {
			fmt.Println("被删除了:", string(v.Key), string(v.Value))
		}
	}

}

