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

func SetAccountInfoCache(rdb *redis.Client, info *model.AccountInfo) error {
	dataBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%v%v", AccountInfoCachePrefix, info.UserId)
	return rdb.Set(ctx, key, string(dataBytes), 0).Err()
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
