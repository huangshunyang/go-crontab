package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		e error
	)

	// 客户端的配置
	config = clientv3.Config{
		Endpoints:[]string{"192.168.0.202:2379"},
		DialTimeout:5 * time.Second,
	}

	// 建立连接
	if client, e = clientv3.New(config); e != nil {
		fmt.Println(e)
		return
	}

	client = client
	fmt.Println(client)
}

