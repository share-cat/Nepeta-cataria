package controllers

import (
	"net/http"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	userProfileService *services.UserProfileService
}

func NewUserProfileController(userProfileService *services.UserProfileService) *UserProfileController {
	return &UserProfileController{
		userProfileService: userProfileService,
	}
}

// GetAllUserProfiles 獲取所有用戶profile
// @Summary 獲取所有用戶profile
// @Description 獲取所有用戶profile的列表
// @Tags userProfiles
// @Accept json
// @Produce json
// @Success 200 {array} models.UserProfile
// @Failure 500 {object} map[string]string
// @Router /userProfiles [get]
func (c *UserProfileController) GetAllUserProfiles(ctx *gin.Context) {
	profiles, err := c.userProfileService.GetAllUserProfiles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, profiles)
}

// GetUserProfileByID 根據ID獲取用戶profile
// @Summary 獲取特定用戶profile
// @Description 根據ID獲取特定用戶profile
// @Tags userProfiles
// @Accept json
// @Produce json
// @Param id path string true "User Profile ID"
// @Success 200 {object} models.UserProfile
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userProfiles/{id} [get]
func (c *UserProfileController) GetUserProfileByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶profile ID",
		})
		return
	}

	profile, err := c.userProfileService.GetUserProfileByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用戶profile未找到",
		})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// GetUserProfileByUserID 根據用戶ID獲取用戶profile
// @Summary 獲取特定用戶profile
// @Description 根據用戶ID獲取特定用戶profile
// @Tags userProfiles
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} models.UserProfile
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userProfiles/user/{userID} [get]
func (c *UserProfileController) GetUserProfileByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	profile, err := c.userProfileService.GetUserProfileByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用戶profile未找到",
		})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// CreateUserProfile 創建用戶profile
// @Summary 創建用戶profile
// @Description 創建新的用戶profile
// @Tags userProfiles
// @Accept json
// @Produce json
// @Param userProfile body models.UserProfile true "User Profile"
// @Success 201 {object} models.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userProfiles [post]
func (c *UserProfileController) CreateUserProfile(ctx *gin.Context) {
	var userProfile models.UserProfile

	if err := ctx.ShouldBindJSON(&userProfile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	profile, err := c.userProfileService.CreateUserProfile(userProfile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, profile)
}

// UpdateUserProfile 更新用戶profile
// @Summary 更新用戶profile
// @Description 更新用戶profile
// @Tags userProfiles
// @Accept json
// @Produce json
// @Param id path string true "User Profile ID"
// @Param userProfile body models.UserProfile true "User Profile"
// @Success 200 {object} models.UserProfile
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userProfiles/{id} [put]
func (c *UserProfileController) UpdateUserProfile(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶profile ID",
		})
		return
	}

	var userProfile models.UserProfile
	if err := ctx.ShouldBindJSON(&userProfile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userProfile.ID = uint(id)

	err = c.userProfileService.UpdateUserProfile(userProfile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用戶profile已更新"})
}

// DeleteUserProfile 刪除用戶profile
// @Summary 刪除用戶profile
// @Description 刪除特定用戶profile
// @Tags userProfiles
// @Accept json
// @Produce json
// @Param id path string true "User Profile ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userProfiles/{id} [delete]
func (c *UserProfileController) DeleteUserProfile(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶profile ID",
		})
		return
	}

	err = c.userProfileService.DeleteUserProfile(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用戶profile已刪除"})
}

// UpdateAvatarURL 更新用戶profile的頭像URL
// @Summary 更新用戶profile的頭像URL
// @Description 更新用戶profile的頭像URL
// @Tags userProfiles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param avatarURL body string true "Avatar URL"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userProfiles/{user_id}/avatar [put]
func (c *UserProfileController) UpdateAvatarURL(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	var avatarURL string
	if err := ctx.ShouldBindJSON(&avatarURL); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.userProfileService.UpdateAvatarURL(uint(userID), avatarURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "頭像URL已更新"})
}
