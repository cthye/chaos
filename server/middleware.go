package server

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"nessaj/config"
	"strings"
)

func abortFail(c *gin.Context, code int, err string) {
	_ = c.AbortWithError(code, errors.New(fmt.Sprintf("%s, remote: %s", err, c.Request.RemoteAddr)))
}

func AuthenticationMiddleware(conf *config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		auth := c.Request.Header["Authorization"]
		if len(auth) == 0 {
			abortFail(c, 400, "Authorization missing")
			return
		}
		segs := strings.SplitN(auth[0], " ", 2)
		if len(segs) != 2 || segs[0] != "Bearer" {
			abortFail(c, 400, "Invalid Authorization header")
			return
		}
		tokenString := segs[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				abortFail(c, 400, fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			}

			return conf.Pubkey, nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			err = claims.Valid()
			if err != nil {
				abortFail(c, 403, fmt.Sprintf("claims no valid: %v", err.Error()))
			}
			if claims["exp"] == nil || claims["iat"] == nil {
				abortFail(c, 403, "exp or iat missing")
			}
		} else {
			if err != nil {
				abortFail(c, 403, fmt.Sprintf("claims no ok: %v", err.Error()))
			} else {
				abortFail(c, 403, "unknown error")
			}
		}
		c.Next()
	}
}
