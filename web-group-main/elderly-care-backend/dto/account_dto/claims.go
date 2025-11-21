package account_dto

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	AccountId uint   `json:"account_id"`
	Nickname  string `json:"nickname"`
	jwt.StandardClaims
}
