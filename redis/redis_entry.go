package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/paceew/go-redisson/bloomfilter"
	"github.com/paceew/go-redisson/glock"
	"github.com/paceew/go-redisson/limiter"
	"net"
	"time"
)

var (
	ObtainLockScript = redis.NewScript(2, `
		if (redis.call('exists', KEYS[1]) == 0) or (redis.call('hexists',KEYS[1],KEYS[2]) == 1) then
			redis.call('hincrby', KEYS[1], KEYS[2],1);
			redis.call('pexpire', KEYS[1], ARGV[1]);
			return 0;
		end;
		return redis.call('pttl', KEYS[1]);
	`)
	ReleaseLockScript = redis.NewScript(3, `
		if redis.call('hexists',KEYS[1],KEYS[2]) == 1 then
			local count = redis.call('hincrby', KEYS[1], KEYS[2], -1);
			if (count <= 0) then
				redis.call('del', KEYS[1]);
				redis.call('publish', KEYS[3], ARGV[1]);
			end;
			return 1;
		end;
		if redis.call('exists', KEYS[1]) == 0 then
			return 0;
		else
			return -1;
		end;
	`)
	RenewExpLocktScript = redis.NewScript(2, `
		if (redis.call('hexists', KEYS[1], KEYS[2]) == 1) then
			redis.call('pexpire', KEYS[1], ARGV[1]);
			return 1;
		end;
		if redis.call('exists', KEYS[1]) == 0 then
			return 0;
		else
			return -1;
		end;
	`)
)

type RedisEntry struct {
	pool *redis.Pool
}

func NewRedisEntry(redisAddr, redisPass string) RedisEntry {
	pool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   0,
		IdleTimeout: time.Minute,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				redisAddr,
				redis.DialConnectTimeout(10*time.Second),
				redis.DialReadTimeout(5*time.Second),
				redis.DialWriteTimeout(5*time.Second),
				redis.DialDatabase(0),
				redis.DialPassword(redisPass),
			)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	re := RedisEntry{
		pool: pool,
	}
	return re
}

func NewRedisEntryWithPool(pool *redis.Pool) RedisEntry {
	re := RedisEntry{
		pool: pool,
	}
	return re
}

func (re RedisEntry) ObtainLock(key, val string, ttl time.Duration) (int64, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int64(ObtainLockScript.Do(conn, key, val, ttl.Milliseconds()))
}

func (re RedisEntry) ReleaseLock(key, val, publishKey, publishMsg string) (int, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int(ReleaseLockScript.Do(conn, key, val, publishKey, publishMsg))
}

func (re RedisEntry) RenewLock(key, val string, ttl time.Duration) (int, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int(RenewExpLocktScript.Do(conn, key, val, ttl.Milliseconds()))
}

func (re RedisEntry) SubscribeAndReceiveMessage(publishKey string, timeout time.Duration) (msg string, err error) {
	conn := re.pool.Get()
	psc := redis.PubSubConn{Conn: conn}
	defer psc.Close()
	err = psc.Subscribe(publishKey)
	if err != nil {
		return "", err
	}
	for {
		switch v := psc.ReceiveWithTimeout(timeout).(type) {
		case redis.Message:
			return string(v.Data), nil
		case redis.Subscription:
		case error:
			if nerr, ok := v.(net.Error); ok && nerr.Timeout() {
				return "", glock.ErrSubTimeout
			}
			return "", v
		}
	}

}

var _ glock.RedisGlockOperate = NewRedisEntry("", "")

/************************************ RedisLimiterOperate *************************************/

func (re RedisEntry) IncrBy(key string, increment int64) (int64, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("INCRBY", key, increment))
}

func (re RedisEntry) PExpire(key string, millisecond time.Duration) (int64, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("PEXPIRE", key, millisecond.Milliseconds()))
}

func (re RedisEntry) PTTL(key string) (int64, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("PTTL", key))
}

var _ limiter.RedisLimiterOperate = NewRedisEntry("", "")

/************************************ RedisBloomFilterOperate *************************************/

func (re RedisEntry) HMGet(keys ...interface{}) (map[string]string, error) {
	conn := re.pool.Get()
	defer conn.Close()

	result := make(map[string]string, len(keys))
	values, err := redis.Values(conn.Do("HMGET", keys...))
	if len(keys)-1 != len(values) {
		return result, errors.New("error keys")
	}
	for i := 0; i < len(values); i++ {
		key, _ := keys[i+1].(string)
		v := values[i]
		if v != nil {
			bys, _ := v.([]byte)
			result[key] = string(bys)
		}
	}

	return result, err
}

func (re RedisEntry) HMSet(kvs ...interface{}) error {
	conn := re.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", kvs...)
	return err
}

func (re RedisEntry) SetBit(key string, offset uint64, value int) error {
	conn := re.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SETBIT", key, offset, value)
	return err
}

func (re RedisEntry) GetBit(key string, offset uint64) (int64, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GETBIT", key, offset))
}

func (re RedisEntry) Expire(key string, second time.Duration) (int64, error) {
	conn := re.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("PEXPIRE", key, second.Seconds()))
}

func (re RedisEntry) Del(key string) error {
	conn := re.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEl", key)
	return err
}

var _ bloomfilter.RedisBloomFilterOperate = NewRedisEntry("", "")
