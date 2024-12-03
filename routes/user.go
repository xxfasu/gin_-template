package routes

import (
	"gin_template/internal/handler"
	"github.com/gin-gonic/gin"
)

func userRouter(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup, userHandler *handler.UserHandler) {
	publicRouter = publicRouter.Group("/user")
	privateRouter = privateRouter.Group("/user")
	{
		publicRouter.POST("/user/login", userHandler.Login)
	}
}
