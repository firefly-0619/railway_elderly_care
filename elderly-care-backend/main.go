package main

import (
	"elderly-care-backend/common/factories"
	"elderly-care-backend/config"
	_ "elderly-care-backend/docs"
	"elderly-care-backend/global"
	"elderly-care-backend/routes"
	"log"
	"os"
)

// @title elderly-care-backend
// @version 1.0
// @description 养老服务平台接口文档
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	defer global.Logger.Sync()
	config.InitConfig()

	// 初始化oss工厂
	factories.InitOssFactory()

	// 初始化服务
	initServices()

	r := routes.SetUpRouter()

	// 获取端口，Railway 会提供 PORT 环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = config.Config.Server.Port // 默认使用配置文件的端口
	}

	log.Printf("🚀 服务器启动在端口 %s", port)
	r.Run(":" + port)
}

func initServices() {
	// 初始化高德地图服务（从配置读取）
	if config.Config.Map.AMap.Enable {
		log.Println("✅ 初始化高德地图服务...")
		// 这里可以赋值给全局变量或在路由中传递
	} else {
		log.Println("⚠️ 高德地图服务未启用")
	}
}
