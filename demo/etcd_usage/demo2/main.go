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
		putResponse *clientv3.PutResponse
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

	fmt.Println(client)

	// KV用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	fmt.Println(kv)
	if putResponse, e = kv.Put(context.TODO(), "/cron/jobs", "aaa", clientv3.WithPrevKV()); e != nil {
		fmt.Println(e)
		return
	}

	if fmt.Println("Revision:", putResponse.Header.Revision); putResponse.PrevKv != nil {
		fmt.Println("prevValue:", string(putResponse.PrevKv.Value))
	}
}

