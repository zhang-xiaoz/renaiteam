package tool

import (
	"web/models"
	modelss "web/models/user"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	SigningKey []byte
}

// 生成一个密码串
func NewJWT() *JWT {
	return &JWT{
		[]byte(models.JWTconfig.Jwt),
	}
}

// 创建一个token//把数据进行加密
func (j *JWT) CreateToken(claims modelss.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//解密
