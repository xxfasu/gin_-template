package middleware

import (
	"bytes"
	"gin_template/pkg/logs"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

type LogM struct {
	logger *logs.Logger
}

func NewLogM(logger *logs.Logger) *LogM {
	return &LogM{
		logger: logger,
	}
}

func (m *LogM) RequestLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// The configuration is initialized once per request
		trace := uuid.NewString()
		m.logger.WithValue(ctx, zap.String("trace", trace))
		m.logger.WithValue(ctx, zap.String("request_method", ctx.Request.Method))
		m.logger.WithValue(ctx, zap.Any("request_headers", ctx.Request.Header))
		m.logger.WithValue(ctx, zap.String("request_url", ctx.Request.URL.String()))
		if ctx.Request.Body != nil {
			bodyBytes, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			m.logger.WithValue(ctx, zap.String("request_params", string(bodyBytes)))
		}
		m.logger.WithContext(ctx).Info("Request")
		ctx.Next()
	}
}
func (m *LogM) ResponseLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime).String()
		m.logger.WithContext(ctx).Info("Response", zap.Any("response_body", blw.body.String()), zap.Any("time", duration))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
