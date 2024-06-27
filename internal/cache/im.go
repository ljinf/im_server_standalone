package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//会话中的消息序列号
	IncrConversationMsgPrefix = "im:conversationmsg:"
)

func IncrConversationMsg(rdb *redis.Client, conversationId int64) int64 {
	key := fmt.Sprintf("%v%v", IncrConversationMsgPrefix, conversationId)
	return rdb.Incr(ctx, key).Val()
}
