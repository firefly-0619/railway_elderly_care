package server_error

import "errors"

var (
	JwtExpireError = errors.New("jwt expire error")
)
