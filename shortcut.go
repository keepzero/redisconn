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

	return convertToInt64(rData)
}

func (r *R) GetString(k string) (str string, err error) {

	// return "" if no record
	rData, err := r.Do("GET", k)
	if err != nil || rData == nil {
		return "", err
	}

	// switch type
	switch strVal := rData.(type) {
	case redis.Error:
		return "", strVal
	case string:
		return strVal, nil
	case int64:
		return fmt.Sprintf("%d", strVal), nil
	case []byte:
		return string(strVal), nil
	case nil:
		return "", nil
	default:
		return "", errors.New("Redis command GET return type unknow")
	}
}

func (r *R) GetInt64List(k string, start int, stop int) ([]int64, error) {

	ret := []int64{}
	rDatas, err := r.Do("LRANGE", k, start, stop)
	if err != nil || rDatas == nil {
		return ret, err
	}

	if datas, ok := rDatas.([]interface{}); ok {
		for _, data := range datas {
			i, err := convertToInt64(data)
			if err != nil {
				return ret, err
			}
			ret = append(ret, i)
		}
	} else {
		return ret, errors.New("Redis command LRANGE return type not []interface{}")
	}

	return ret, nil
}

func (r *R) GetIntList(k string, start int, stop int) ([]int, error) {

	ret := []int{}
	retInt64, err := r.GetInt64List(k, start, stop)
	if err != nil {
		return ret, err
	}

	for _, i := range retInt64 {
		ret = append(ret, int(i))
	}

	return ret, nil
}

func (r *R) GetInt32List(k string, start int, stop int) ([]int32, error) {

	ret := []int32{}
	retInt64, err := r.GetInt64List(k, start, stop)
	if err != nil {
		return ret, err
	}

	for _, i := range retInt64 {
		ret = append(ret, int32(i))
	}

	return ret, nil
}

// Do with int value returned
// example: DoInt("INCR", "key")
func (r *R) DoInt(cmd string, args ...interface{}) (int, error) {
	//return 0, nil
	rData, err := r.Do(cmd, args...)
	if err != nil || rData == nil {
		return 0, err
	}
	ret, err := convertToInt64(rData)
	return int(ret), err
}

// TODO
// DoString

func convertToInt64(val interface{}) (int64, error) {

	// switch type
	switch intVal := val.(type) {
	case redis.Error:
		return 0, intVal
	case int64:
		return intVal, nil
	case string:
		return strconv.ParseInt(intVal, 10, 64)
	case []byte:
		return strconv.ParseInt(fmt.Sprintf("%s", string(intVal)), 10, 64)
	case nil:
		return 0, nil
	case []interface{}:
		return 0, errors.New("Cant't convert []interface{} to int64")
	default:
		return 0, errors.New("Redis return type unknow")
	}
}
