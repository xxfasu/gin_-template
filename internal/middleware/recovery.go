package middleware

import (
	"fmt"
	"gin_template/pkg/errors"
	"gin_template/pkg/logs"
	"gin_template/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http/httputil"
	"strings"
	"time"
)

type Recovery struct {
}

func NewRecoveryM() *Recovery {
	return &Recovery{}
}

func (m *Recovery) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rv := recover(); rv != nil { // 捕获任何发生的恐慌（panic）
				ctx := c.Request.Context()

				var fields []zap.Field
				fields = append(fields, zap.Strings("error", []string{fmt.Sprintf("%v", rv)})) // 添加错误信息字段
				fields = append(fields, zap.StackSkip("stack", 0))                             // 添加堆栈信息，跳过指定的调用层级

				if gin.IsDebugging() { // 如果处于调试模式
					httpRequest, _ := httputil.DumpRequest(c.Request, false) // 转储HTTP请求（不包含请求体）
					headers := strings.Split(string(httpRequest), "\r\n")    // 按行分割请求头
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" { // 如果是授权头
							headers[idx] = current[0] + ": *" // 隐藏授权信息
						}
					}
					fields = append(fields, zap.Strings("headers", headers)) // 添加处理后的请求头字段
				}

				// 记录错误日志，包含时间戳和相关字段
				logs.Log.WithContext(ctx).Error(fmt.Sprintf("[Recovery] %s panic recovered", time.Now().Format("2006/01/02 - 15:04:05")), fields...)
				// 返回内部服务器错误响应给客户端
				utils.ResError(c, errors.InternalServerError("", "Internal server error, please try again later"))
			}
		}()

		c.Next() // 继续处理下一个中间件或请求处理函数
	}
}
