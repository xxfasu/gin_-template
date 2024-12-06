package routes

import (
	"gin_template/internal/handler/v1/user_handler"
	"github.com/gin-gonic/gin"
)

func userRouter(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup, userHandler *user_handler.UserHandler) {
	publicRouter = publicRouter.Group("/user")
	privateRouter = privateRouter.Group("/user")
	{
		publicRouter.POST("/register", userHandler.Register)
		publicRouter.POST("/login", userHandler.Login)
		publicRouter.POST("/find", userHandler.FindUser)
	}
}
