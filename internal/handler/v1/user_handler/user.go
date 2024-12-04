package user_handler

import (
	"gin_template/internal/model/request"
	"gin_template/internal/model/response"
	"gin_template/internal/service/user_service"
	"gin_template/pkg/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	logger      *logs.Logger
	userService user_service.UserService
}

func NewUserHandler(logger *logs.Logger, userService user_service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(request.Register)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		response.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	response.HandleSuccess(ctx, nil)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req request.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		response.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	response.HandleSuccess(ctx, response.Login{
		AccessToken: token,
	})
}
