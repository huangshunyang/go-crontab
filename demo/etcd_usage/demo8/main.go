package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		e error
		kv clientv3.KV
		putOp clientv3.Op
		getOp clientv3.Op
		opResponse clientv3.OpResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"192.168.0.202:2379"},
		DialTimeout:5 * time.Second,
	}

	// 建立一个客户端
	if client, e = clientv3.New(config); e != nil {
		fmt.Println(e)
		return
	}

	// KV用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	putOp = clientv3.OpPut("/cron/jobs/job8", "2")

	// 执行op
	if opResponse, e = kv.Do(context.TODO(), putOp); e != nil {
		fmt.Println(e)
		return
	}

	getOp = clientv3.OpGet("/cron/jobs/job8")
	if opResponse, e = kv.Do(context.TODO(), getOp); e != nil {
		fmt.Println(e)
		return
	}

	// 打印
	fmt.Println("数据Revision：", opResponse.Get().Kvs[0].ModRevision)
	fmt.Println("数据Value：", string(opResponse.Get().Kvs[0].Value))
}
