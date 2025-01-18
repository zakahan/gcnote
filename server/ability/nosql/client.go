// -------------------------------------------------
// Package nosql
// Author: hanzhi
// Date: 2025/1/18
// -------------------------------------------------

package nosql

import (
	"context"
	"fmt"
	"gcnote/server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// 插入文档到指定集合
func InsertDocument(collectionName string, document interface{}) error {
	// 获取 MongoDB 客户端和数据库
	db := config.MongoDBConf.Database(config.ServerCfg.MongodbConf.DBName) // 数据库名为 testdb
	collection := db.Collection(collectionName)

	// 执行插入操作
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		zap.S().Errorf("插入文档失败: %v", err)
		return fmt.Errorf("插入文档失败: %v", err)
	}
	zap.S().Infof("文档成功插入到集合: %s", collectionName)
	return nil
}

// 查询指定集合中的所有文档
func FindDocuments(collectionName string) ([]bson.M, error) {
	// 获取 MongoDB 客户端和数据库
	db := config.MongoDBConf.Database(config.ServerCfg.MongodbConf.DBName) // 数据库名为 testdb
	collection := db.Collection(collectionName)

	// 查询所有文档
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		zap.S().Errorf("查询文档失败: %v", err)
		return nil, fmt.Errorf("查询文档失败: %v", err)
	}
	defer cursor.Close(context.Background())

	// 读取所有文档到切片
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		zap.S().Errorf("读取文档失败: %v", err)
		return nil, fmt.Errorf("读取文档失败: %v", err)
	}

	zap.S().Infof("查询到 %d 个文档", len(results))
	// 转化格式

	return results, nil
}

// 删除指定集合中的某个文档，按文档名称（假设文档中有 "name" 字段）
func DeleteDocument(collectionName string, documentName string) error {
	// 获取 MongoDB 客户端和数据库
	db := config.MongoDBConf.Database(config.ServerCfg.MongodbConf.DBName) // 数据库名为 config.ServerCfg.MongodbConf.DBName
	collection := db.Collection(collectionName)

	// 定义查询条件，假设文档中有 "name" 字段
	filter := bson.M{"name": documentName}

	// 执行删除操作
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		zap.S().Errorf("删除文档失败: %v", err)
		return fmt.Errorf("删除文档失败: %v", err)
	}
	if result.DeletedCount == 0 {
		zap.S().Infof("未找到要删除的文档: %s", documentName)
		return fmt.Errorf("未找到要删除的文档: %s", documentName)
	}

	zap.S().Infof("成功删除文档: %s", documentName)
	return nil
}
