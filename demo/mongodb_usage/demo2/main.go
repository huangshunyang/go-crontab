package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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

func main() {

	var (
		client *mongo.Client
		e error
		db *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertOneResult
		docId primitive.ObjectID
		clientOptions *options.ClientOptions
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

	// 4。插入记录
	record = &LogRecord{
		JobName:"job1",
		Command: "echo hello",
		Error:"",
		Content:"hello",
		TimePoint:TimePoint{StartTime:time.Now().Unix(), EndTime:time.Now().Unix() + 10},
	}

	if result, e = collection.InsertOne(context.TODO(), record); e != nil {
		fmt.Println(e)
		return
	}

	// _id默认生成一个全局唯一ID，objectID，12字节的二进制
	docId = result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID:", docId.Hex())
}

