package redisconn

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

func (r *R) GetInt(k string) (int, error) {
	i, err := r.GetInt64(k)
	return int(i), err
}

func (r *R) GetInt32(k string) (int32, error) {
	i, err := r.GetInt64(k)
	return int32(i), err
}

func (r *R) GetInt64(k string) (int64, error) {

	// return 0 if no record
	rData, err := r.Do("GET", k)
	if err != nil || rData == nil {
		return int64(0), err
	}

	// switch type
	switch intVal := rData.(type) {
	case int64:
		return intVal, nil
	case string:
		return strconv.ParseInt(intVal, 10, 64)
	case []byte:
		return strconv.ParseInt(fmt.Sprintf("%s", string(intVal)), 10, 64)
	case nil:
		return int64(0), nil
	case redis.Error:
		return int64(0), intVal
	default:
		return int64(0), errors.New("redis command return type not unknow")
	}
}

func (r *R) GetString(k string) (str string, err error) {

	// return "" if no record
	rData, err := r.Do("GET", k)
	if err != nil || rData == nil {
		return "", err
	}

	// switch type
	switch strVal := rData.(type) {
	case string:
		return strVal, nil
	case int64:
		return fmt.Sprintf("%d", strVal), nil
	case []byte:
		return string(strVal), nil
	case nil:
		return "", nil
	case redis.Error:
		return "", strVal
	default:
		return "", errors.New("redis command return type not unknow")
	}
}
