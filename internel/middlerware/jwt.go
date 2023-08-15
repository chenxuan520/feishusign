package middlerware

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type jwtInfo struct {
	UserId string
	jwt.StandardClaims
}

func GenerateJwt(uid string) (string, error) {
	claims := jwtInfo{
		UserId: uid,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.GlobalConfig.Sign.ExpireDuration).Unix(),
			Issuer:    config.GlobalConfig.Sign.Issuer,
		},
	}
	tokenCla := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenCla.SignedString([]byte(config.GlobalConfig.Sign.JwtToken)); err == nil {
		token = "Bearer " + token
		return token, nil
	} else {
		return "", err
	}
}

func VerifyJwt(auth string) (string, error) {
	arr := strings.Fields(auth)
	if len(arr) < 2 {
		err := fmt.Errorf("wrong token")
		return "", err
	}
	auth = arr[1]
	token, err := jwt.ParseWithClaims(auth, &jwtInfo{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.Sign.JwtToken), nil
	})
	if err != nil {
		return "", err
	}
	return token.Claims.(*jwtInfo).UserId, nil
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			context.Abort()
			err := fmt.Errorf("no auth existing in header")
			logger.GetLogger().Error(err.Error())
			response.Error(context, http.StatusUnauthorized, err)
			return
		}
		userId, err := VerifyJwt(auth)
		if err != nil {
			context.Abort()
			logger.GetLogger().Error(err.Error())
			response.Error(context, http.StatusForbidden, err)
		}
		context.Set("uid", userId)
		context.Next()
	}
}

func Debug() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Set("uid", "123456")
		g.Next()
	}
}
