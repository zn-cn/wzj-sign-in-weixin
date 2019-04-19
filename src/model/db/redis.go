package db

import (
	"config"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// redis controller
type RedisDBCntlr struct {
	conn redis.Conn
}

var globalRedisPool *redis.Pool
var redisURL string
var redisPW string

func init() {
	redisConf := config.Conf.Redis
	redisURL = fmt.Sprintf("%s:%s", redisConf.Host, redisConf.Port)
	redisPW = redisConf.PW
	globalRedisPool = GetRedisPool()
}

// GetRedisPool get the client pool of redis
func GetRedisPool() *redis.Pool {
	pool := &redis.Pool{ // 实例化一个连接池
		MaxIdle:     30, // 最大的连接数量
		MaxActive:   0,  // 连接池最大连接数量,不确定可以用0（0表示自动定义）
		IdleTimeout: 60, // 连接关闭时间 60秒 （60秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { // 要连接的redis数据库
			conn, err := redis.Dial("tcp", redisURL)
			if err != nil {
				return nil, err
			}
			if redisPW != "" {
				_, err = conn.Do("AUTH", redisPW)
			}
			return conn, err
		},
	}
	return pool
}

/********************************************* RedisDBCntlr *******************************************/

func NewRedisDBCntlr() *RedisDBCntlr {
	return &RedisDBCntlr{
		conn: globalRedisPool.Get(),
	}
}

func (this *RedisDBCntlr) Close() {
	this.conn.Close()
}

func (this *RedisDBCntlr) GetConn() redis.Conn {
	return this.conn
}

func (this *RedisDBCntlr) Send(commandName string, args ...interface{}) error {
	return this.conn.Send(commandName, args...)
}
