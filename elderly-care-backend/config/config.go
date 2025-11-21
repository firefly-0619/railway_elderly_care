package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Server struct {
		Name    string `mapstructure:"name"`
		Port    string `mapstructure:"port"`
		Profile string `mapstructure:"profile"`
		NodeId  int64  `mapstructure:"node_id"`
	}
	Database struct {
		Host            string `mapstructure:"host"`
		Port            string `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		Db              string `mapstructure:"db"`
		IdleConnections int    `mapstructure:"idle_connections"`
		MaxConnections  int    `mapstructure:"max_connections"`
		MaxLifetime     int    `mapstructure:"max_lifetime"`
	}
	Redis struct {
		Host              string `mapstructure:"host"`
		Port              string `mapstructure:"port"`
		Password          string `mapstructure:"password"`
		Db                int    `mapstructure:"db"`
		ConnectionTimeout int    `mapstructure:"connection_timeout"`
	}
	Mongodb struct {
		Uri            string `mapstructure:"uri"`
		ConnectTimeout int    `mapstructure:"connect_timeout"`
		Db             string `mapstructure:"db"`
	}
	Kafka struct {
		Async     bool     `mapstructure:"async"`
		Addresses []string `mapstructure:"address"`
	}
	Oss struct {
		Aliyun Oss
		Minio  Oss
	}
	Log struct {
		Output string `mapstructure:"output"`
	}
	Jwt struct {
		SecretKey     string `mapstructure:"secret_key"`
		Expire        int    `mapstructure:"expire"`
		RefreshExpire int    `mapstructure:"refresh_expire"`
	}
	// 新增 Map 配置
	Map struct {
		AMap struct {
			Enable  bool   `mapstructure:"enable"`
			APIKey  string `mapstructure:"api_key"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"amap"`
	} `mapstructure:"map"`
}

type Oss struct {
	Enable          bool   `mapstructure:"enable"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSsl          bool   `mapstructure:"use_ssl"`
}

var Config *AppConfig

var logFile *os.File

func InitConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config,error: %v", err)
	}
	Config = &AppConfig{}
	if err = viper.Unmarshal(Config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// 🆕 新增：环境变量覆盖（Railway 部署用）
	overrideConfigWithEnv()

	logFile, err = os.OpenFile(Config.Log.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//初始化数据库
	initDB()
	//初始化logger
	initLogger()
	//初始化雪花生成器
	initSnowFlake()
	//初始化redis客户端
	initRedis()
	//初始化redis分布式锁
	initRedSync()
	//初始化kafka客户端
	initKafka()
}

// 🆕 修正函数：使用环境变量覆盖配置
func overrideConfigWithEnv() {
	// Server 配置
	if port := os.Getenv("PORT"); port != "" {
		Config.Server.Port = port
	}
	if nodeID := os.Getenv("NODE_ID"); nodeID != "" {
		if id, err := strconv.ParseInt(nodeID, 10, 64); err == nil { // 修正：使用 ParseInt 而不是 Atoi
			Config.Server.NodeId = id
		}
	}

	// Database 配置
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		Config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		Config.Database.Port = dbPort // 修正：Database.Port 是 string 类型
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		Config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		Config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		Config.Database.Db = dbName
	}

	// Redis 配置
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		Config.Redis.Host = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		Config.Redis.Port = redisPort // 修正：Redis.Port 是 string 类型
	}
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		if db, err := strconv.Atoi(redisDB); err == nil {
			Config.Redis.Db = db
		}
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		Config.Redis.Password = redisPassword
	}

	// JWT 配置
	if jwtSecret := os.Getenv("JWT_SECRET_KEY"); jwtSecret != "" {
		Config.Jwt.SecretKey = jwtSecret
	}
	if jwtExpire := os.Getenv("JWT_EXPIRE"); jwtExpire != "" {
		if expire, err := strconv.Atoi(jwtExpire); err == nil {
			Config.Jwt.Expire = expire
		}
	}

	// 高德地图配置
	if amapKey := os.Getenv("AMAP_API_KEY"); amapKey != "" {
		Config.Map.AMap.APIKey = amapKey
	}
	if amapEnable := os.Getenv("AMAP_ENABLE"); amapEnable != "" {
		Config.Map.AMap.Enable = (amapEnable == "true")
	}
	if amapBaseURL := os.Getenv("AMAP_BASE_URL"); amapBaseURL != "" {
		Config.Map.AMap.BaseURL = amapBaseURL
	}

	// MinIO 配置
	if minioEndpoint := os.Getenv("MINIO_ENDPOINT"); minioEndpoint != "" {
		Config.Oss.Minio.Endpoint = minioEndpoint
	}
	if minioAccessKey := os.Getenv("MINIO_ACCESS_KEY"); minioAccessKey != "" {
		Config.Oss.Minio.AccessKeyId = minioAccessKey
	}
	if minioSecretKey := os.Getenv("MINIO_SECRET_KEY"); minioSecretKey != "" {
		Config.Oss.Minio.SecretAccessKey = minioSecretKey
	}
	if minioEnable := os.Getenv("MINIO_ENABLE"); minioEnable != "" {
		Config.Oss.Minio.Enable = (minioEnable == "true")
	}

	// Aliyun OSS 配置
	if aliyunEnable := os.Getenv("ALIYUN_OSS_ENABLE"); aliyunEnable != "" {
		Config.Oss.Aliyun.Enable = (aliyunEnable == "true")
	}
	if aliyunEndpoint := os.Getenv("ALIYUN_OSS_ENDPOINT"); aliyunEndpoint != "" {
		Config.Oss.Aliyun.Endpoint = aliyunEndpoint
	}
	if aliyunAccessKey := os.Getenv("ALIYUN_OSS_ACCESS_KEY"); aliyunAccessKey != "" {
		Config.Oss.Aliyun.AccessKeyId = aliyunAccessKey
	}
	if aliyunSecretKey := os.Getenv("ALIYUN_OSS_SECRET_KEY"); aliyunSecretKey != "" {
		Config.Oss.Aliyun.SecretAccessKey = aliyunSecretKey
	}

	// Kafka 配置
	if kafkaAddress := os.Getenv("KAFKA_ADDRESS"); kafkaAddress != "" {
		Config.Kafka.Addresses = []string{kafkaAddress} // 修正：字段名是 Addresses 不是 Address
	}
	if kafkaAsync := os.Getenv("KAFKA_ASYNC"); kafkaAsync != "" {
		Config.Kafka.Async = (kafkaAsync == "true")
	}

	// MongoDB 配置（如果需要）
	if mongodbURI := os.Getenv("MONGODB_URI"); mongodbURI != "" {
		Config.Mongodb.Uri = mongodbURI
	}
	if mongodbDB := os.Getenv("MONGODB_DB"); mongodbDB != "" {
		Config.Mongodb.Db = mongodbDB
	}

	// 日志配置
	if logOutput := os.Getenv("LOG_OUTPUT"); logOutput != "" {
		Config.Log.Output = logOutput
	}
}
