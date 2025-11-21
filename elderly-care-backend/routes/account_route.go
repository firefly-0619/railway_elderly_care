package routes

import (
	"elderly-care-backend/controllers"
	"fmt" // 添加这行

	"github.com/gin-gonic/gin"
)

func AccountRoute(e *gin.Engine) {
	fmt.Println("   📍 注册账户路由组: /account")
	controller := controllers.AccountController{}
	accountRoute := e.Group("/account")
	{
		accountRoute.POST("/register", controller.Register)
		fmt.Println("     ✅ POST /account/register")

		accountRoute.POST("/login", controller.Login)
		fmt.Println("     ✅ POST /account/login")

		accountRoute.PUT("", controller.UpdateAccount)
		fmt.Println("     ✅ PUT /account")

		accountRoute.PUT("/changePassword", controller.ChangePassword)
		fmt.Println("     ✅ PUT /account/changePassword")

		accountRoute.GET("/checkPhone", controller.CheckPhoneIsExists)
		fmt.Println("     ✅ GET /account/checkPhone")

		accountRoute.GET("", controller.GetAccountInfo)
		fmt.Println("     ✅ GET /account")
		accountRoute.GET("/:accountID", controller.GetAccountInfoByAccountID)
	}
	fmt.Println("   ✅ 账户路由注册完成")
}
