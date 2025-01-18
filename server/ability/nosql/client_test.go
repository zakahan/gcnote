// -------------------------------------------------
// Package nosql
// Author: hanzhi
// Date: 2025/1/18
// -------------------------------------------------

package nosql

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

func TestDemo(t *testing.T) {
	// 设置 MongoDB 连接 URI
	uri := "mongodb://localhost:27017"

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建客户端并连接
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = client.Disconnect(disconnectCtx); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	// 选择数据库和集合
	collection := client.Database("testdb").Collection("users")

	// 插入文档
	insertCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	insertResult, err := collection.InsertOne(insertCtx, bson.M{
		"name":  "Alice",
		"age":   25,
		"email": "alice@example.com",
	})
	if err != nil {
		log.Fatal("Insert failed:", err)
	}
	fmt.Println("Inserted document ID:", insertResult.InsertedID)

	// 查询文档
	findCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result bson.M
	err = collection.FindOne(findCtx, bson.M{"name": "Alice"}).Decode(&result)
	if err != nil {
		log.Fatal("Find failed:", err)
	}
	fmt.Println("Found document:", result)

	// 删除文档
	deleteCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	deleteResult, err := collection.DeleteOne(deleteCtx, bson.M{"name": "Alice"})
	if err != nil {
		log.Fatal("Delete failed:", err)
	}
	fmt.Println("Deleted document count:", deleteResult.DeletedCount)
}
