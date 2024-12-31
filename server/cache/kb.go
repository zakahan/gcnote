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
// s:kb:info:{kb_file_id} - 存储单个kb文件的信息
// s:kb:list:{index_id} - 存储index下的kb文件列表
// s:kb:recent:{user_id} - 存储用户最近访问的kb文件列表

func kbInfoKey(kbFileId string) string {
	return fmt.Sprintf("s:kb:info:%v", kbFileId)
}

func kbListKey(indexId string) string {
	return fmt.Sprintf("s:kb:list:%v", indexId)
}

func kbRecentKey(userId string) string {
	return fmt.Sprintf("s:kb:recent:%v", userId)
}

// GetKBInfo 获取单个kb文件信息
func GetKBInfo(ctx context.Context, kbFileId string) (*model.KBFile, error) {
	result, err := config.RedisClient.Get(ctx, kbInfoKey(kbFileId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var kbFile model.KBFile
	err = json.Unmarshal([]byte(result), &kbFile)
	if err != nil {
		return nil, err
	}
	return &kbFile, nil
}

// SetKBInfo 设置单个kb文件信息
func SetKBInfo(ctx context.Context, kbFile model.KBFile) error {
	marshal, err := json.Marshal(kbFile)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, kbInfoKey(kbFile.KBFileId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelKBInfo 删除单个kb文件信息
func DelKBInfo(ctx context.Context, kbFileId string) error {
	_, err := config.RedisClient.Del(ctx, kbInfoKey(kbFileId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetIndexKBList 获取index下的kb文件列表
func GetIndexKBList(ctx context.Context, indexId string) ([]model.KBFile, error) {
	result, err := config.RedisClient.Get(ctx, kbListKey(indexId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var kbFiles []model.KBFile
	err = json.Unmarshal([]byte(result), &kbFiles)
	if err != nil {
		return nil, err
	}
	return kbFiles, nil
}

// SetIndexKBList 设置index下的kb文件列表
func SetIndexKBList(ctx context.Context, indexId string, kbFiles []model.KBFile) error {
	marshal, err := json.Marshal(kbFiles)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, kbListKey(indexId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelIndexKBList 删除index下的kb文件列表
func DelIndexKBList(ctx context.Context, indexId string) error {
	_, err := config.RedisClient.Del(ctx, kbListKey(indexId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetRecentKBList 获取用户最近访问的kb文件列表
func GetRecentKBList(ctx context.Context, userId string) ([]model.KBFile, error) {
	result, err := config.RedisClient.Get(ctx, kbRecentKey(userId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var kbFiles []model.KBFile
	err = json.Unmarshal([]byte(result), &kbFiles)
	if err != nil {
		return nil, err
	}
	return kbFiles, nil
}

// SetRecentKBList 设置用户最近访问的kb文件列表
func SetRecentKBList(ctx context.Context, userId string, kbFiles []model.KBFile) error {
	marshal, err := json.Marshal(kbFiles)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, kbRecentKey(userId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelRecentKBList 删除用户最近访问的kb文件列表
func DelRecentKBList(ctx context.Context, userId string) error {
	_, err := config.RedisClient.Del(ctx, kbRecentKey(userId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// RefreshKBInfo 刷新单个kb文件信息
func RefreshKBInfo(ctx context.Context, kbFileId string) (*model.KBFile, error) {
	resp := model.KBFile{}
	tx := config.DB.Where("kb_file_id = ?", kbFileId).First(&resp)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, tx.Error
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, nil
	}
	err := SetKBInfo(ctx, resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

// RefreshIndexKBList 刷新index下的kb文件列表
func RefreshIndexKBList(ctx context.Context, indexId string) ([]model.KBFile, error) {
	var kbFiles []model.KBFile
	tx := config.DB.Where("index_id = ?", indexId).Find(&kbFiles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	err := SetIndexKBList(ctx, indexId, kbFiles)
	if err != nil {
		return kbFiles, err
	}
	return kbFiles, nil
}

// RefreshRecentKBList 刷新用户最近访问的kb文件列表
func RefreshRecentKBList(ctx context.Context, userId string) ([]model.KBFile, error) {
	var kbFiles []model.KBFile
	tx := config.DB.Where("user_id = ?", userId).Order("updated_at desc").Limit(10).Find(&kbFiles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	err := SetRecentKBList(ctx, userId, kbFiles)
	if err != nil {
		return kbFiles, err
	}
	return kbFiles, nil
}
