package controllers

import (
	"net/http"
	"share-cat/Nepeta-cataria/services"
	"share-cat/Nepeta-cataria/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userService *services.UserService
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

// Login 用戶登入
// @Summary 用戶登入
// @Description 用戶登入，返回用戶信息和token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登入請求"
// @Success 200 {object} LoginResponse "登入成功"
// @Failure 400 {object} map[string]string "請求參數無效"
// @Failure 401 {object} map[string]string "用戶名或密碼錯誤"
// @Failure 500 {object} map[string]string "伺服器內部錯誤"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "請求參數無效"})
		return
	}

	user, err := c.userService.GetUserByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用戶名或密碼錯誤"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用戶名或密碼錯誤"})
		return
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	auth, err := utils.CreateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失敗"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token:        auth.Token,
		RefreshToken: auth.RefreshToken,
		UserID:       user.ID,
		Username:     user.Username,
	})
}

// ReflashToken 刷新token
// @Summary 刷新token
// @Description 刷新token，返回新的token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新請求"
// @Success 200 {object} RefreshTokenResponse "刷新成功"
// @Failure 400 {object} map[string]string "請求參數無效"
// @Failure 401 {object} map[string]string "無效的 Refresh Token"
// @Failure 500 {object} map[string]string "伺服器內部錯誤"
// @Router /auth/refresh [post]
func (c *AuthController) ReflashToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "請求參數無效"})
		return
	}

	newToken, err := utils.ReflashToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 Refresh Token"})
		return
	}

	ctx.JSON(http.StatusOK, RefreshTokenResponse{
		Token: newToken,
	})
}
