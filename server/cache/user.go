// -------------------------------------------------
// Package cache
// Author: hanzhi
// Date: 2024/12/10
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
func userInfoKey(userId string) string {
	return fmt.Sprintf("s:userinfo:%v", userId)
}

func GetUserInfo(ctx context.Context, userId string) (*model.User, error) {
	result, err := config.RedisClient.Get(ctx, userInfoKey(userId)).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, redis.Nil
	}

	var user model.User
	err = json.Unmarshal([]byte(result), &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SetUserInfo(ctx context.Context, user model.User) error {
	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, userInfoKey(user.UserId), marshal, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func DelUserInfo(ctx context.Context, userId string) error {
	_, err := config.RedisClient.Del(ctx, userInfoKey(userId)).Result()
	if err != nil {
		return err
	}
	return nil
}

// RefreshUserInfo 刷新用户信息缓存
func RefreshUserInfo(ctx context.Context, userId string) (*model.User, error) {
	resp := model.User{}
	tx := config.DB.Where("user_id = ?", userId).First(&resp)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, tx.Error
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &resp, nil
	}
	err := SetUserInfo(ctx, resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}
