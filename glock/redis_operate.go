package glock

import (
	"errors"
	"time"
)

const (
	MSG_UNLOCK string = "RELEASE"
)

var (
	ErrSubTimeout = errors.New("subscribe timeout")
)

type RedisGlockOperate interface {
	/*可重入方式获取锁，获取成功返回0，获取失败返回锁对应pttl
	获取锁应该设置对应的key val 并加上过期时间ttl,如果锁已存在则应该计数+1*/
	ObtainLock(key, val string, ttl time.Duration) (int64, error)
	/*释放锁，释放成功返回1，锁不存在返回0，val不对返回-1
	释放锁应该计数-1，如果计数归零应删除对应key，并向publishKey推送publishMsg消息*/
	ReleaseLock(key, val, publishKey, publishMsg string) (int, error)

	/*刷新锁，刷新成功返回1，锁不存在返回0，val不对返回-1
	刷新锁应该刷新key val 对应过期时间为 ttl*/
	RenewLock(key, val string, ttl time.Duration) (int, error)
	/*订阅channel渠道并接受相应订阅消息
	返回ErrSubTimeout错误为订阅超时*/
	SubscribeAndReceiveMessage(publishKey string, timeout time.Duration) (msg string, err error)
}
