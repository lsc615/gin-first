package common

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtCustomClaims struct {
	ID                   int
	Name                 string
	jwt.RegisteredClaims // 注意!这是jwt-go的v4版本新增的，原先是jwt.StandardClaims
}

var stSignKey = []byte("手写的从前") // 定义secret，后面会用到

func GenerateToken(id int, name string) (tokenString string, err error) {
	// 初始化
	claim := JwtCustomClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(stSignKey)          // stsignKey必须是byte字节的,所以我们在设置签名秘钥,要使用byte强转
	println(token)
	return tokenString, err
}
