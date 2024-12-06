package user_handler

import (
	"gin_template/internal/data/service_data"
	"gin_template/internal/service/user_service"
	"gin_template/internal/validation"
	"gin_template/pkg/logs"
	"gin_template/pkg/utils"
	"github.com/gin-gonic/gin"
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
	req := new(validation.Register)
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}
	utils.ResSuccess(ctx, nil)

}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req validation.Login
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}
	utils.ResSuccess(ctx, service_data.LoginResp{
		AccessToken: token,
	})
}

func (h *UserHandler) FindUser(ctx *gin.Context) {
	var req validation.FindUser
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}

	user, err := h.userService.FindUser(ctx, &req)
	if err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}
	utils.ResSuccess(ctx, user)
}

func (h *UserHandler) ModifyPassword(ctx *gin.Context) {
	var req validation.Login
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		utils.ResError(ctx, h.logger, err)
		return
	}
	utils.ResSuccess(ctx, service_data.LoginResp{
		AccessToken: token,
	})
}
