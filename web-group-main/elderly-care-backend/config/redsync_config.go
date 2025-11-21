package config

import (
	"elderly-care-backend/global"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

func initRedSync() {
	pool := goredis.NewPool(global.RedisClient)
	lockManager := redsync.New(pool)
	global.LockManager = lockManager
}
