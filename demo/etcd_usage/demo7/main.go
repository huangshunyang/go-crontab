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
		client *clientv3.Client
		e error
		kv clientv3.KV
		getResponse *clientv3.GetResponse
		watchStartRevision int64
		watcher clientv3.Watcher
		watchResponseChan clientv3.WatchChan
		watchResponse clientv3.WatchResponse
		event *clientv3.Event
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

	//	模拟kv的变化
	go func() {
		for  {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")
			kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续的变化
	if getResponse, e = kv.Get(context.TODO(), "/cron/jobs/job7"); e != nil {
		fmt.Println(e)
		return
	}

	// 现在key是存在的
	if len(getResponse.Kvs) != 0 {
		fmt.Println("当前的值:", string(getResponse.Kvs[0].Value))
	}

	// 当前etcd集群事务ID，单调递增的
	watchStartRevision = getResponse.Header.Revision + 1

	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本向后监听：", watchStartRevision)

	// 开启监听，每次监听到变化后会向watchResponseChan这个管道中写入一个watchResponse
	watchResponseChan = watcher.Watch(context.TODO(), "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	// 处理kv的变化事件，watchResponse中的Event可以拿变化后的信息
	for watchResponse = range watchResponseChan {
		for _, event = range watchResponse.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}
