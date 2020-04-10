package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
}


//一条日志
type LogRecord struct {
	JobName string `bson:"jobName"` // 任务名
	Command string `bson:"command"`// 脚本命令
	Error string `bson:"error"`// 脚本错误输出
	Content string `bson:"content"`// 脚本输出
	TimePoint TimePoint `bson:"timePoint"`// 任务的执行时间点
}

// jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"` // 将jobName赋值为job2
}

func main() {

	var (
		client *mongo.Client
		e error
		db *mongo.Database
		collection *mongo.Collection
		clientOptions *options.ClientOptions
		condition *FindByJobName
		cursor *mongo.Cursor
		findOpt *options.FindOptions
		record *LogRecord
	)

	// 设置客户端连接配置
	clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	if client, e = mongo.Connect(context.TODO(), clientOptions); e != nil {
		fmt.Println(e)
		return
	}

	// 检查连接
	if e = client.Ping(context.TODO(), nil); e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("Connected to MongoDB!")

	// 2.选择数据库
	db = client.Database("cron")

	// 3.选择表
	collection = db.Collection("log")

	// 4.按照jobName字段过滤，想找出jobName=job2，找出3条
	condition = &FindByJobName{JobName:"job2"} // {"jobName": "job10"}

	findOpt = options.Find()
	findOpt.SetSkip(0)
	findOpt.SetLimit(3)

	// 5.查询过滤
	if cursor, e = collection.Find(context.TODO(), condition, findOpt); e != nil {
		fmt.Println(e)
		return
	}

	// 函数退出时关闭游标
	defer cursor.Close(context.TODO())

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		// 定义一个日志的对象
		record = &LogRecord{}

		// 反序列化bson对象
		if e = cursor.Decode(record); e != nil {
			fmt.Println(e)
			return
		}
		// 把日志打印出来
		fmt.Println(*record)
	}
}

