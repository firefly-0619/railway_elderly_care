package middlewares

import (
	"elderly-care-backend/common/constants"
	"elderly-care-backend/common/custom"
	"elderly-care-backend/config"
	"elderly-care-backend/utils"
	"elderly-care-backend/vo"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyMiddleware() gin.HandlerFunc {
	whiteList := []interface{}{
		"/account/login",
		"/account/register",
		"/account/checkPhone",
		"/music/list",
		"/account/refresh",
	}

	whiteSet := custom.NewHashSet(whiteList...)

	return func(c *gin.Context) {
		// === 添加调试信息 ===
		token := c.GetHeader("Authorization")
		if whiteSet.Contains(c.Request.URL.Path) {
			c.Next()
			return
		}

		if token == "" {
			fmt.Println(c.Request.URL.Path)
			c.JSON(http.StatusGone, vo.Fail(constants.NOT_LOGIN))
			c.Abort()
			return
		}

		// 添加这行：去除 "Bearer " 前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
			fmt.Printf("去除Bearer前缀后: %s\n", token)
		}

		claims, err := utils.ParseToken(token, config.Config.Jwt.SecretKey)
		if err != nil {
			fmt.Printf("Token解析错误: %v\n", err)
			c.JSON(http.StatusGone, vo.Fail(constants.INVALID_TOKEN))
			c.Abort()
			return
		}
		fmt.Printf("Token验证成功，用户ID: %d\n", claims.AccountId)

		utils.SetAccountIdInContext(c, claims.AccountId)
		utils.SetNickNameInContext(c, claims.Nickname)
		c.Next()
	}
}
