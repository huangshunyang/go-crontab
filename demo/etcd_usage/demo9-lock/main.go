package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 基于ETCD实现分布式乐观锁
func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		e error
		lease clientv3.Lease
		leaseGrantResponse *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		keepaliveResponseChan <-chan *clientv3.LeaseKeepAliveResponse
		keepaliveResponse *clientv3.LeaseKeepAliveResponse
		ctx context.Context
		cancelFunc context.CancelFunc
		kv clientv3.KV
		txn clientv3.Txn
		lockKey string
		txnResponse *clientv3.TxnResponse
	)

	lockKey = "/cron/lock/job9"

	config = clientv3.Config{
		Endpoints:[]string{"192.168.0.202:2379"},
		DialTimeout:5 * time.Second,
	}

	// 建立一个客户端
	if client, e = clientv3.New(config); e != nil {
		fmt.Println(e)
		return
	}

	// lease实现锁的自动过期
	// op操作
	// txn事务： if else then


	// 1.上锁（创建租约，自动续租，拿着租约去抢占一个key）
	lease = clientv3.NewLease(client)

	// 申请一个5秒钟的租约
	if leaseGrantResponse, e = lease.Grant(context.TODO(), 5); e != nil {
		fmt.Println(e)
		return
	}

	// 拿到租约的ID
	leaseId = leaseGrantResponse.ID

	// 准备一个用于取消的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	// 确保函数退出后，自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	// 5秒钟后会取消自动续租
	if keepaliveResponseChan, e = lease.KeepAlive(ctx, leaseId); e != nil {
		fmt.Println(e)
		return
	}

	// 处理续约应答的协程
	go func() {
		for  {
			select {
			case keepaliveResponse = <- keepaliveResponseChan:
				if keepaliveResponseChan == nil {
					fmt.Println("租约已经失败了")
					goto END
				} else {
					fmt.Println(keepaliveResponse)
				}
			}
		}
		END:
	}()

	// if不存在key，then设置它，else抢锁失败
	kv = clientv3.KV(client)

	// 创建一个事务
	txn = kv.Txn(context.TODO())

	// 如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey)) // 否则抢锁失败

	// 提交事务
	if txnResponse, e = txn.Commit(); e != nil {
		fmt.Println(e)
		return
	}

	// 判断是否抢到了锁
	if !txnResponse.Succeeded {
		fmt.Println("锁被占用:", string(txnResponse.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}


	// 2.处理业务 (在锁内，很安全)
	// 3.释放锁（取消自动续租，释放租约）

	fmt.Println("处理任务！")
	time.Sleep(time.Second * 10)
}
