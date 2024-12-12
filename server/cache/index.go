// -------------------------------------------------
// Package cache
// Author: hanzhi
// Date: 2024/12/11
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

// Redis key 规范 s:service:xxx  s代表key的类型为string service为服务名 xxx为自定义值
// s代表string
// hs代表hashmap
// se代表set
// zs代表zset
// l代表list
// bf代表布隆过滤器
// hy代表hyperloglog
// b代表bitmap
func indexInfoKey(indexId string) string {
	return fmt.Sprintf("s:indexinfo:%v", indexId)
}

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

func DelIndexInfo(ctx context.Context, indexId string) error {
	_, err := config.RedisClient.Del(ctx, indexInfoKey(indexId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// RefreshIndexInfo 刷新用户信息缓存
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
