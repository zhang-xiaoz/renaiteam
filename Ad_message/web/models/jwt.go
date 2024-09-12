package models

import "github.com/dgrijalva/jwt-go"

var JWTconfig *JWTConfig = &JWTConfig{} //jwt保密验证信息

//jwt
type JWTConfig struct {
	Jwt string `json:"jwt"`
}

type CustomClaims struct { //jwt所用到的加密解密东西
	Mailbox string //用户唯一标识
	jwt.StandardClaims
}

type Jwt_Check struct {
	Access_token  string
	Refresh_token string
}
