package routes

import (
	"gin_template/internal/handler"
	"gin_template/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

var ProviderSet = wire.NewSet(NewRouter)

func NewRouter(
	userHandler *handler.UserHandler,
	recoveryM *middleware.Recovery,
	corsM *middleware.Cors,
	logM *middleware.LogM,
) *gin.Engine {
	router := gin.New()
	if true {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(recoveryM.Handler())
	publicGroup := router.Group("/api")
	privateGroup := router.Group("/api")

	{
		// 健康监测
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	{
		userRouter(publicGroup, privateGroup, userHandler)
	}
	return router
}
