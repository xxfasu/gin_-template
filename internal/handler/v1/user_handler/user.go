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

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body response.Register true "params"
// @Success 200 {object} response.Response
// @Router /register [post]
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

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body response.Login true "params"
// @Success 200 {object} response.LoginResponse
// @Router /login [post]
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
