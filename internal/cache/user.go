package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/redis/go-redis/v9"
)

var (
	ctx                    = context.Background()
	cachePrefix            = "im:server:"
	AccountInfoCachePrefix = cachePrefix + "user:info:"
)

func SetAccountInfoCache(rdb *redis.Client, info ...model.AccountInfo) error {

	list := make(map[string]interface{})

	for _, v := range info {
		dataBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%v%v", AccountInfoCachePrefix, v.UserId)

		list[key] = string(dataBytes)
	}

	return rdb.MSet(ctx, list).Err()
}

func GetAccountInfoCache(rdb *redis.Client, userId int64) (*model.AccountInfo, error) {
	key := fmt.Sprintf("%v%v", AccountInfoCachePrefix, userId)
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if result != "" {
		info := &model.AccountInfo{}
		if err := json.Unmarshal([]byte(result), info); err != nil {
			return nil, err
		}
		return info, nil
	}

	return nil, errors.New("not found")
}

func GetUserInfoListCache(rdb *redis.Client, uids ...string) ([]model.UserInfo, error) {
	var (
		length = len(uids)
		list   = make([]model.UserInfo, 0, length)
	)

	if length > 0 {
		keys := make([]string, 0, length)
		for _, v := range uids {
			keys = append(keys, fmt.Sprintf("%v%v", AccountInfoCachePrefix, v))
		}

		result, err := rdb.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		for _, v := range result {
			item := model.UserInfo{}
			if err = json.Unmarshal([]byte(v.(string)), &item); err != nil {
				return nil, err
			}
			list = append(list, item)
		}

	}
	return list, nil
}
