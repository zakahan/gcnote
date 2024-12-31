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
// s:recycle:info:{kb_file_id} - 存储单个回收站文件的信息
// s:recycle:list:{user_id} - 存储用户的回收站文件列表

func recycleInfoKey(kbFileId string) string {
	return fmt.Sprintf("s:recycle:info:%v", kbFileId)
}

func recycleListKey(userId string) string {
	return fmt.Sprintf("s:recycle:list:%v", userId)
}

// GetRecycleInfo 获取单个回收站文件信息
func GetRecycleInfo(ctx context.Context, kbFileId string) (*model.Recycle, error) {
	result, err := config.RedisClient.Get(ctx, recycleInfoKey(kbFileId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var recycle model.Recycle
	err = json.Unmarshal([]byte(result), &recycle)
	if err != nil {
		return nil, err
	}
	return &recycle, nil
}

// SetRecycleInfo 设置单个回收站文件信息
func SetRecycleInfo(ctx context.Context, recycle model.Recycle) error {
	marshal, err := json.Marshal(recycle)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, recycleInfoKey(recycle.KBFileId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelRecycleInfo 删除单个回收站文件信息
func DelRecycleInfo(ctx context.Context, kbFileId string) error {
	_, err := config.RedisClient.Del(ctx, recycleInfoKey(kbFileId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetUserRecycleList 获取用户的回收站文件列表
func GetUserRecycleList(ctx context.Context, userId string) ([]model.Recycle, error) {
	result, err := config.RedisClient.Get(ctx, recycleListKey(userId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var recycles []model.Recycle
	err = json.Unmarshal([]byte(result), &recycles)
	if err != nil {
		return nil, err
	}
	return recycles, nil
}

// SetUserRecycleList 设置用户的回收站文件列表
func SetUserRecycleList(ctx context.Context, userId string, recycles []model.Recycle) error {
	marshal, err := json.Marshal(recycles)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, recycleListKey(userId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelUserRecycleList 删除用户的回收站文件列表
func DelUserRecycleList(ctx context.Context, userId string) error {
	_, err := config.RedisClient.Del(ctx, recycleListKey(userId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// RefreshRecycleInfo 刷新单个回收站文件信息
func RefreshRecycleInfo(ctx context.Context, kbFileId string) (*model.Recycle, error) {
	resp := model.Recycle{}
	tx := config.DB.Where("kb_file_id = ?", kbFileId).First(&resp)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, tx.Error
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, nil
	}
	err := SetRecycleInfo(ctx, resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

// RefreshUserRecycleList 刷新用户的回收站文件列表
func RefreshUserRecycleList(ctx context.Context, userId string) ([]model.Recycle, error) {
	var recycles []model.Recycle
	tx := config.DB.Where("user_id = ?", userId).Find(&recycles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	err := SetUserRecycleList(ctx, userId, recycles)
	if err != nil {
		return recycles, err
	}
	return recycles, nil
}
