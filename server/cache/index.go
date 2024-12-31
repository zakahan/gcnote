// -------------------------------------------------
// Package cache
// Author: hanzhi
// Date: 2024/12/30
// -------------------------------------------------

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gcnote/server/config"
	"gcnote/server/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Redis key规范
// s:index:info:{index_id} - 存储单个index的信息
// s:index:list:{user_id} - 存储用户的index列表

func indexInfoKey(indexId string) string {
	return fmt.Sprintf("s:index:info:%v", indexId)
}

func indexListKey(userId string) string {
	return fmt.Sprintf("s:index:list:%v", userId)
}

// GetIndexInfo 获取单个index信息
func GetIndexInfo(ctx context.Context, indexId string) (*model.Index, error) {
	result, err := config.RedisClient.Get(ctx, indexInfoKey(indexId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var index model.Index
	err = json.Unmarshal([]byte(result), &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// SetIndexInfo 设置单个index信息
func SetIndexInfo(ctx context.Context, index model.Index) error {
	marshal, err := json.Marshal(index)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, indexInfoKey(index.IndexId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelIndexInfo 删除单个index信息
func DelIndexInfo(ctx context.Context, indexId string) error {
	_, err := config.RedisClient.Del(ctx, indexInfoKey(indexId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetUserIndexList 获取用户的index列表
func GetUserIndexList(ctx context.Context, userId string) ([]model.Index, error) {
	result, err := config.RedisClient.Get(ctx, indexListKey(userId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var indexes []model.Index
	err = json.Unmarshal([]byte(result), &indexes)
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

// SetUserIndexList 设置用户的index列表
func SetUserIndexList(ctx context.Context, userId string, indexes []model.Index) error {
	marshal, err := json.Marshal(indexes)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, indexListKey(userId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelUserIndexList 删除用户的index列表
func DelUserIndexList(ctx context.Context, userId string) error {
	_, err := config.RedisClient.Del(ctx, indexListKey(userId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// RefreshIndexInfo 刷新单个index信息
func RefreshIndexInfo(ctx context.Context, indexId string) (*model.Index, error) {
	resp := model.Index{}
	tx := config.DB.Where("index_id = ?", indexId).First(&resp)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, tx.Error
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, nil
	}
	err := SetIndexInfo(ctx, resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

// RefreshUserIndexList 刷新用户的index列表
func RefreshUserIndexList(ctx context.Context, userId string) ([]model.Index, error) {
	var indexes []model.Index
	tx := config.DB.Where("user_id = ?", userId).Find(&indexes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	err := SetUserIndexList(ctx, userId, indexes)
	if err != nil {
		return indexes, err
	}
	return indexes, nil
}
