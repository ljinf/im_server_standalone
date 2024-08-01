package cache

import (
	"encoding/json"
	"fmt"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
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
	MsgExpire          = 259200 //72 hour
)

func IncrConversationMsg(rdb *redis.Client, conversationId int64) int64 {
	key := fmt.Sprintf("%v%v", IncrConversationMsgPrefix, conversationId)
	return rdb.Incr(ctx, key).Val()
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

// 用户的会话设置  hash类型
func SetUserConversationCache(rdb *redis.Client, info *model.UserConversationList) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v%v", UserConversationInfoPrefix, info.UserId)
	if err = rdb.HSet(ctx, key, info.ConversationId, string(data)).Err(); err != nil {
		return err
	}
	return rdb.Expire(ctx, key, time.Duration(rand.Intn(randTime)+userConversationExpire)).Err()
}

func GetUserConversationCache(rdb *redis.Client, userId int64, convIds ...int64) ([]model.UserConversationList, error) {
	var (
		length = len(convIds)
		list   = make([]model.UserConversationList, 0, length)
		key    = fmt.Sprintf("%v%v", UserConversationInfoPrefix, userId)
	)

	if length > 0 {
		fields := make([]string, 0, length)
		for _, v := range convIds {
			fields = append(fields, strconv.Itoa(int(v)))
		}

		result, err := rdb.HMGet(ctx, key, fields...).Result()
		if err != nil {
			return nil, err
		}

		for _, v := range result {
			item := model.UserConversationList{}
			if err = json.Unmarshal([]byte(v.(string)), &item); err != nil {
				return nil, err
			}
			list = append(list, item)
		}
	}

	return list, nil
}

func DelUserConversationCache(rdb *redis.Client, userId int64, convIds ...string) error {
	key := fmt.Sprintf("%v%v", UserConversationInfoPrefix, userId)
	return rdb.HDel(ctx, key, convIds...).Err()
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

// 获取会话的用户ID列表
func GetConversationUserListCache(rdb *redis.Client, convId int64) ([]string, error) {
	key := fmt.Sprintf("%v%v", ConversationUserListPrefix, convId)
	return rdb.SMembers(ctx, key).Result()
}

// 会话下的用户数
func GetConversationUserCount(rdb *redis.Client, convId int64) (int64, error) {
	key := fmt.Sprintf("%v%v", ConversationUserListPrefix, convId)
	return rdb.SCard(ctx, key).Result()
}

// 会话消息链 ,最近消息
func AddConversationMsgCache(rdb *redis.Client, msg *model.MsgResp) error {
	key := fmt.Sprintf("%v%v", ConversationMsgListPrefix, msg.ConversationId)

	return rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(msg.Seq),
		Member: msg.MsgId,
	}).Err()
}

// 会话下的消息Id列表
func GetConversationMsgList(rdb *redis.Client, convId, pageNum, pageSize int64) ([]string, error) {
	var (
		key   = fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
		start = (pageNum - 1) * pageSize
		end   = start + pageNum - 1
	)
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

// 会话下的消息总数
func GetConversationMsgCount(rdb *redis.Client, convId int64) (int64, error) {
	key := fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
	return rdb.ZCard(ctx, key).Result()
}

// 删除会话下的部分数量的消息(按seq升序)
func RemConversationMsg(rdb *redis.Client, convId, num int64) error {
	key := fmt.Sprintf("%v%v", ConversationMsgListPrefix, convId)
	return rdb.ZRemRangeByRank(ctx, key, 0, num).Err()
}
