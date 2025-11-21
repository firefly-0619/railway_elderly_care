package utils

import (
	"elderly-care-backend/common/server_error"
	"elderly-care-backend/dto/account_dto"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(fromPassword), err
}

func CheckPasswordHash(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func GenToken(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString string, secretKey string) (*account_dto.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &account_dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
		fmt.Printf("Token签名方法: %v\n", token.Method)
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Printf("Token解析错误: %v\n", err) // 这里会显示具体错误
		return nil, err
	}

	if claims, ok := token.Claims.(*account_dto.Claims); ok && token.Valid {
		fmt.Printf("Token验证成功，用户ID: %d\n", claims.AccountId)
		return claims, nil
	} else {
		fmt.Printf("Token验证失败，Token有效: %t\n", token.Valid)
		return nil, server_error.JwtExpireError
	}
}

func SetAccountIdInContext(c *gin.Context, accountId uint) {
	c.Set("account_id", accountId)
}
func SetNickNameInContext(c *gin.Context, nickName string) {
	c.Set("nickname", nickName)
}

func GetAccountIdInContext(c *gin.Context) uint {

	value, exists := c.Get("account_id")

	if !exists {
		return 0
	}
	return value.(uint)
}

func GetNickNameInContext(c *gin.Context) string {

	value, exists := c.Get("nickname")

	if !exists {
		return ""
	}
	return value.(string)
}
