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
		lease clientv3.Lease
		leaseGrantResponse *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		kv clientv3.KV
		keepResponseChan <-chan *clientv3.LeaseKeepAliveResponse
		keepResponse *clientv3.LeaseKeepAliveResponse
		putResponse *clientv3.PutResponse
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


	// 申请一个lease （租约）
	lease = clientv3.NewLease(client)

	// 申请一个10秒钟的租约
	if leaseGrantResponse, e = lease.Grant(context.TODO(), 10); e != nil {
		fmt.Println(e)
		return
	}

	// 拿到租约的ID
	leaseId = leaseGrantResponse.ID
	if keepResponseChan, e = lease.KeepAlive(context.TODO(), leaseId); e != nil {
		fmt.Println(e)
		return
	}


	// 处理续租应答的协程
	go func() {
		for  {
			select {
			case keepResponse = <-keepResponseChan:
				if keepResponseChan == nil {
					fmt.Println("租约已经失效了！")
					goto END
				} else {
					// 每秒续租一次，所以会受到一次应答
					fmt.Println("收取自动续租应答：", keepResponse.ID)
				}
			}
		}
		END:
	}()

	// 获取kvApi的子集
	kv = clientv3.NewKV(client)

	// put一个kv，让它与租约关联起来，实现10秒后自动过期
	if putResponse, e = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); e != nil {
		fmt.Println(e)
		return
	}

	fmt.Println("写入成功:", putResponse.Header.Revision)

	for  {
		if getResponse, e = kv.Get(context.TODO(), "/cron/lock/job1"); e != nil {
			fmt.Println(e)
			return
		}

		if getResponse.Count == 0 {
			fmt.Println("过期了")
			break
		}

		fmt.Println("还没有过期:", getResponse.Kvs)
		time.Sleep(time.Second * 2)
	}
}

