// global/global.go
package global

import (
	"elderly-care-backend/common/custom"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger         *zap.Logger
	Db             *gorm.DB
	SnowFlake      *snowflake.Node
	RedisClient    *redis.Client
	LockManager    *redsync.Redsync
	KafkaOperators map[string]*custom.KafkaOperator // 不用管
	// 注意：这里不再直接声明服务实例，避免循环导入
)
