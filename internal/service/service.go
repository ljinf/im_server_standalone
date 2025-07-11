package service

import (
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"github.com/ljinf/im_server_standalone/pkg/jwt"
	"github.com/ljinf/im_server_standalone/pkg/log"
	"github.com/ljinf/im_server_standalone/pkg/sid"
	"math/rand"
	"time"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	cache  *cache.Cache
	tm     repository.Transaction
}

func NewService(
	cache *cache.Cache,
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
		cache:  cache,
	}
}

// 随机数字
func GenerateRandomNumberString(length int) string {
	const digits = "0123456789"
	rand.Seed(time.Now().UnixNano()) // 使用当前时间的纳秒数来初始化随机数生成器
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
