package user_handler

import (
	"gin_template/internal/data/service_data"
	"gin_template/internal/service/user_service"
	"gin_template/internal/validation"
	"gin_template/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user_service.UserService
}

func NewUserHandler(userService user_service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(validation.Register)
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, err)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		utils.ResError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, nil)

}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req validation.Login
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, err)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		utils.ResError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, service_data.LoginResp{
		AccessToken: token,
	})
}

func (h *UserHandler) FindUser(ctx *gin.Context) {
	var req validation.FindUser
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, err)
		return
	}

	user, err := h.userService.FindUser(ctx, &req)
	if err != nil {
		utils.ResError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, user)
}

func (h *UserHandler) ModifyPassword(ctx *gin.Context) {
	var req validation.Login
	if err := utils.ParseJSON(ctx, &req); err != nil {
		utils.ResError(ctx, err)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		utils.ResError(ctx, err)
		return
	}
	utils.ResSuccess(ctx, service_data.LoginResp{
		AccessToken: token,
	})
}
