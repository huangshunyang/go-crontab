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
		e error
		client *clientv3.Client
		kv clientv3.KV
		getResponse *clientv3.GetResponse
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

	if getResponse, e = kv.Get(context.TODO(),"/cron/jobs/job1"); e != nil {
		fmt.Println(e)
		return
	}

	fmt.Println(getResponse.Kvs)
}

