package middleware

import (
	"errors"
	"gin_template/pkg/jwt"
	"gin_template/pkg/logs"
	"gin_template/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthM struct {
	logger *logs.Logger
	j      *jwt.JWT
}

func NewAuthM(logger *logs.Logger, j *jwt.JWT) *AuthM {
	return &AuthM{
		logger: logger,
		j:      j,
	}
}

func (m *AuthM) StrictAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			m.logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			utils.ResError(ctx, m.logger, errors.New("token is empty"))
			ctx.Abort()
			return
		}

		claims, err := m.j.ParseToken(tokenString)
		if err != nil {
			m.logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			utils.ResError(ctx, m.logger, errors.New("token is valid"))
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, m.logger)
		ctx.Next()
	}
}

func (m *AuthM) NoStrictAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			ctx.Next()
			return
		}

		claims, err := m.j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, m.logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *logs.Logger) {
	if userInfo, ok := ctx.MustGet("claims").(*jwt.MyCustomClaims); ok {
		logger.WithValue(ctx, zap.String("UserId", userInfo.UserID))
	}
}
