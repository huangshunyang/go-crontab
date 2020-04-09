package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var (
		client *mongo.Client
		e error
		db *mongo.Database
		collection *mongo.Collection
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
	db = client.Database("my_db")

	// 3.选择表
	collection = db.Collection("my_collection")

	collection = collection
}

