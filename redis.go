package redisconn

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type R struct {
	pool *redis.Pool
}

// Open a new *R
func Open(network, address, password string) (*R, error) {

	// Test connection
	_, err := redis.Dial(network, address)

	// redigo pool
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, address)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// Return *R
	return &R{
		pool: pool,
	}, err
}

// Do something
func (r *R) Do(cmd string, args ...interface{}) (reply interface{}, err error) {

	conn := r.pool.Get()
	defer conn.Close()

	return conn.Do(cmd, args...)
}

// get a redis connection
func (r *R) Get() redis.Conn {
	return r.pool.Get()
}
