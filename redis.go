package redisconn

import (
	"github.com/garyburd/redigo/redis"
)

type R struct {
	network string
	address string
	conn    redis.Conn
}

func (r *R) reconnect() error {

	var err error

	// reconnect if conn is nil
	if r.conn != nil {
		if r.conn.Err() != nil {
			// close old conn
			r.conn.Close()
			r.conn, err = redis.Dial(r.network, r.address)
			return err
		}
	} else {
		r.conn, err = redis.Dial(r.network, r.address)
		return err
	}
	return nil
}

func Open(network, address string) (*R, error) {

	redis, err := redis.Dial(network, address)

	return &R{
		network: network,
		address: address,
		conn:    redis,
	}, err
}

func (r *R) Do(cmd string, args ...interface{}) (reply interface{}, err error) {

	// reconnect before use
	if err = r.reconnect(); err != nil {
		return nil, err
	}

	return r.conn.Do(cmd, args...)
}

// TODO
// func Send
// func Receive
// func Flush
