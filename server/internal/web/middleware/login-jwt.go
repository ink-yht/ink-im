package middlelware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"ink-im-server/internal/web/user_web"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 步骤三
		// 不需要登录的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// Bearer xxx
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 || segs[0] != "Bearer" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc user_web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return user_web.JWTKey, nil
		})
		if err != nil {
			// token 不对, token 是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			// token 不对，token 解析出来了，但可能是非法的，或者过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 增强系统的安全性
		if uc.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题  加监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		// 不判定也可以
		if expireTime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 剩余过期时间 《 1h 刷新
		if expireTime.Sub(time.Now()) < 1*time.Hour {
			// 刷新
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err := token.SignedString(user_web.JWTKey)
			if err != nil {
				// 这边不要中断，仅仅是过期时间没有刷新，但是用户登陆了
				log.Panicln(err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		// 缓存
		ctx.Set("claims", uc)
	}
}
