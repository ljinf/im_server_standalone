package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/redis/go-redis/v9"
	"math"
	"math/rand"
	"time"
)

var (
	randTime = 172800 //2 day

	//会话中的消息序列号
	IncrConversationMsgPrefix = cachePrefix + "conversationmsg:seq:"

	// 会话
	ConversationInfoPrefix = cachePrefix + "conversation:info:"
	ConversationExpire     = 259200 //72 hour
	//会话下的用户列表
	ConversationUserListPrefix = cachePrefix + "conversation:userlist:"
	//会话消息链
	ConversationMsgListPrefix = cachePrefix + "conversation:msglist:"

	//用户会话
	UserConversationInfoPrefix = cachePrefix + "user:conversation:info:"
	userConversationExpire     = 259200

	//消息
	MsgInfoCachePrefix = cachePrefix + "msg:info:"
	MsgExpire          = 604800 //7天
)

// 加1
func IncrConversationMsg(rdb *redis.Client, conversationId int64) int64 {
	key := fmt.Sprintf("%v%v", IncrConversationMsgPrefix, conversationId)
	return rdb.Incr(ctx, key).Val()
}

// 减1
func DecrConversationMsg(rdb *redis.Client, conversationId int64) {
	key := fmt.Sprintf("%v%v", IncrConversationMsgPrefix, conversationId)
	rdb.IncrBy(ctx, key, -1)
}

// 最近消息缓存  String类型
func SetMsgCache(rdb *redis.Client, msg *model.MsgResp) error {
	msgData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v%v", MsgInfoCachePrefix, msg.MsgId)
	return rdb.Set(ctx, key, string(msgData), time.Duration(MsgExpire)*time.Second).Err()
}

func GetMsgCache(rdb *redis.Client, msgIds ...interface{}) ([]model.MsgResp, error) {

	var (
		length  = len(msgIds)
		msgList = make([]model.MsgResp, 0, length)
	)

	if length > 0 {
		keys := make([]string, 0, length)
		for _, v := range msgIds {
			keys = append(keys, fmt.Sprintf("%v%v", MsgInfoCachePrefix, v))
		}

		result, err := rdb.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		for _, v := range result {
			item := model.MsgResp{}
			if err = json.Unmarshal([]byte(v.(string)), &item); err != nil {
				return nil, err
			}
			msgList = append(msgList, item)
		}
	}

	return msgList, nil
}

// 会话  String类型
func SetConversationCache(rdb *redis.Client, conv *model.ConversationList) error {
	convData, err := json.Marshal(conv)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v%v", ConversationInfoPrefix, conv.ConversationId)
	return rdb.Set(ctx, key, string(convData), time.Duration(rand.Intn(randTime)+ConversationExpire)*time.Second).Err()
}

func GetConversationCache(rdb *redis.Client, convIds ...int64) ([]model.ConversationList, error) {
	var (
		length   = len(convIds)
		convList = make([]model.ConversationList, 0, length)
	)

	if length > 0 {
		keys := make([]string, 0, length)
		for _, v := range convIds {
			keys = append(keys, fmt.Sprintf("%v%v", ConversationInfoPrefix, v))
		}

		result, err := rdb.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		for _, v := range result {
			item := model.ConversationList{}
			if err = json.Unmarshal([]byte(v.(string)), &item); err != nil {
				return nil, err
			}
			convList = append(convList, item)
		}
	}

	return convList, nil
}

// 用户的会话设置  zset类型
func SetUserConversationCache(rdb *redis.Client, info ...model.UserConversationList) error {
	if len(info) > 0 {
		list := make([]redis.Z, 0, len(info))
		key := fmt.Sprintf("%v%v", UserConversationInfoPrefix, info[0].UserId)
		for _, v := range info {
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			list = append(list, redis.Z{
				Score:  float64(v.ConversationId),
				Member: string(data),
			})
		}
		if err := rdb.ZAdd(ctx, key, list...).Err(); err != nil {
			return err
		}
		return rdb.Expire(ctx, key, time.Duration(rand.Intn(randTime)+userConversationExpire)).Err()
	}
	return nil
}

// 列表
func GetUserConversationListCache(rdb *redis.Client, userId, pageNum, pageSize int64) ([]model.UserConversationList, error) {
	var (
		list  = make([]model.UserConversationList, 0, pageSize)
		key   = fmt.Sprintf("%v%v", UserConversationInfoPrefix, userId)
		start = (pageNum - 1) * pageSize
		end   = start + pageSize - 1
	)

	result, err := rdb.ZRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		item := model.UserConversationList{}
		if err = json.Unmarshal([]byte(v), &item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

// 单个
func GetUserConversationCache(rdb *redis.Client, userId, convId int64) (*model.UserConversationList, error) {
	key := fmt.Sprintf("%v%v", UserConversationInfoPrefix, userId)
	min := fmt.Sprintf("%v", convId)
	result, err := rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: min,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		item := &model.UserConversationList{}
		if err = json.Unmarshal([]byte(result[0]), item); err != nil {
			return nil, err
		}
		return item, nil
	}
	return nil, errors.New("not Found")
}

func DelUserConversationCache(rdb *redis.Client, userId int64) error {
	return rdb.Del(ctx, fmt.Sprintf("%v%v", UserConversationInfoPrefix, userId)).Err()
}

// 会话下的用户列表(群聊)  set类型
func AddConversationUserListCache(rdb *redis.Client, convId int64, uids ...int64) error {
	key := fmt.Sprintf("%v%v", ConversationUserListPrefix, convId)
	members := make([]interface{}, 0, len(uids))
	for _, v := range uids {
		members = append(members, v)
	}
	if err := rdb.SAdd(ctx, key, members...).Err(); err != nil {
		return err
	}
	return rdb.Expire(ctx, key, time.Duration(rand.Intn(randTime)+ConversationExpire)*time.Second).Err()
}

func RemConversationUserListCache(rdb *redis.Client, convId int64, uids ...int64) error {
	key := fmt.Sprintf("%v%v", ConversationUserListPrefix, convId)
	members := make([]interface{}, 0, len(uids))
	for _, v := range uids {
		members = append(members, v)
	}
	return rdb.SRem(ctx, key, members...).Err()
}

// 获取会话的用户列表
func GetConversationUserListCache(rdb *redis.Client, convId int64) ([]model.UserInfo, error) {
	key := fmt.Sprintf("%v%v", ConversationUserListPrefix, convId)
	userIds, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return GetUserInfoListCache(rdb, userIds...)
}

// 会话下的用户数
func GetConversationUserCount(rdb *redis.Client, convId int64) (int64, error) {
	key := fmt.Sprintf("%v%v", ConversationUserListPrefix, convId)
	return rdb.SCard(ctx, key).Result()
}

// 会话消息链 ,最近消息
func AddConversationMsgCache(rdb *redis.Client, msgs ...model.MsgResp) error {
	if len(msgs) > 0 {
		key := fmt.Sprintf("%v%v", ConversationMsgListPrefix, msgs[0].ConversationId)
		cacheList := make([]redis.Z, 0, len(msgs))
		for _, v := range msgs {
			cacheList = append(cacheList, redis.Z{
				Score:  float64(v.Seq),
				Member: v.MsgId,
			})
		}
		return rdb.ZAdd(ctx, key, cacheList...).Err()
	}
	return nil
}

// 会话下的最近消息列表
func GetConversationMsgList(rdb *redis.Client, convId, seq, pageNum, pageSize int64) ([]model.MsgResp, error) {
	var (
		key   = fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
		start = (pageNum - 1) * pageSize
		end   = start + pageNum - 1
	)

	msgIds, err := rdb.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%v", seq),
		Max:    fmt.Sprintf("%v", math.MaxInt64),
		Offset: start,
		Count:  end,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(msgIds) < 1 {
		return nil, errors.New("msgIds is nil")
	}

	ids := make([]interface{}, 0, len(msgIds))
	for _, v := range msgIds {
		ids = append(ids, v)
	}

	msgList, err := GetMsgCache(rdb, ids...)
	if err != nil {
		return nil, err
	}
	return msgList, nil
}

// 最新一条
func GetConversationNewestMsg(rdb *redis.Client, convId int64) (*model.MsgResp, error) {
	var (
		key = fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
	)

	msgIds, err := rdb.ZRevRange(ctx, key, 0, 0).Result()
	if err != nil {
		return nil, err
	}

	if len(msgIds) > 0 {
		msgList, err := GetMsgCache(rdb, msgIds[0])
		if err != nil {
			return nil, err
		}
		return &msgList[0], nil
	}

	return nil, errors.New("not found")
}

// 会话下的最近消息总数
func GetConversationMsgCount(rdb *redis.Client, convId int64) (int64, error) {
	key := fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
	return rdb.ZCard(ctx, key).Result()
}

// 删除会话下的部分数量的消息(按seq升序)
func RemConversationMsg(rdb *redis.Client, convId, num int64) error {
	key := fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
	return rdb.ZRemRangeByRank(ctx, key, 0, num).Err()
}
