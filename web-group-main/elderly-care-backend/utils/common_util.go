package utils

import (
	. "elderly-care-backend/global"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

func WithDefault[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

func FromStringSliceToUintSlice(slice []string) []uint {
	res := make([]uint, len(slice))
	for i, v := range slice {
		t, _ := strconv.Atoi(v)
		res[i] = uint(t)
	}
	return res
}

// 通用重试函数
func RetryWhenError[T any](retryMaxTime time.Duration, retryInterval time.Duration, function func(T) error, args T) {
	timer := time.NewTimer(retryMaxTime)
	ticker := time.NewTicker(retryInterval)
	var err error
loop:
	for {
		select {
		case <-timer.C:

			Logger.Error(fmt.Sprintf("Retry %s exhausted", getFunctionName(function)), zap.Error(err))
			break loop
		case <-ticker.C:
			err = function(args)
			if err == nil {
				timer.Stop()
				ticker.Stop()
				break loop
			}
		}
	}
}

// 获取函数名的辅助函数
func getFunctionName(f interface{}) string {
	rv := reflect.ValueOf(f)
	if rv.Kind() != reflect.Func {
		return ""
	}

	// 获取函数指针并解析函数名
	ptr := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	if ptr != nil {
		return ptr.Name()
	}
	return ""
}

func CoverStr2Int(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}
