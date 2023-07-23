package tools

import (
	"github.com/dgrijalva/jwt-go"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"time"
)

const (
	expireDuration = 300 * time.Second
	issuer         = "DianGroup"
)

type Claims struct {
	UserId   string
	Username string
	jwt.StandardClaims
}

func GenerateJwtToken(userId string, username string) (string, error) {
	expireTime := time.Now().Add(expireDuration)
	claims := Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.GlobalConfig.Feishu.EncryptKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseJwtToken 解析jwtToken，可以拿到其中的username和userId
func ParseJwtToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.Feishu.EncryptKey), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {

			return claims, nil
		}
	}
	return nil, err
}

// VerifyJwtToken 验证jwtToken是否正确
func VerifyJwtToken(token string) (bool, error) {
	_, err := ParseJwtToken(token)
	if err != nil {
		return false, err
	}

	return true, nil
}
