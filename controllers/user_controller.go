package controllers

import (
	"net/http"
	"strconv"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUsers 獲取所有用戶
// @Summary 獲取所有用戶
// @Description 獲取所有用戶的列表
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUser 根據ID獲取用戶
// @Summary 獲取特定用戶
// @Description 根據ID獲取特定用戶
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用戶不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// CreateUser 創建新用戶
// @Summary 創建新用戶
// @Description 創建一個新用戶
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User Info"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdUser, err := c.userService.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

// UpdateUser 更新用戶
// @Summary 更新用戶
// @Description 更新特定用戶的信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User Info"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	var userUpdate models.User
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.userService.UpdateUser(uint(id), userUpdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用戶更新成功",
	})
}

// UpdatePassword 更新用戶密碼
// @Summary 更新用戶密碼
// @Description 更新特定用戶的密碼
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param passwords body map[string]string true "Old and New Passwords"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /users/{id}/password [put]
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	var passwords struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&passwords); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.userService.UpdateUserPassword(uint(id), passwords.OldPassword, passwords.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "密碼更新成功",
	})
}

// DeleteUser 刪除用戶
// @Summary 刪除用戶
// @Description 刪除特定用戶
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用戶刪除成功",
	})
}

// GetUserProfile 獲取用戶資料
// @Summary 獲取用戶資料
// @Description 獲取特定用戶的詳細資料
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserProfile
// @Failure 400,404 {object} map[string]string
// @Router /users/{id}/profile [get]
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	profile, err := c.userService.GetUserProfile(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用戶資料不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
