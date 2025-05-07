package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type DiscussionController struct {
	discussionService *services.DiscussionService
}

func NewDiscussionController(discussionService *services.DiscussionService) *DiscussionController {
	return &DiscussionController{
		discussionService: discussionService,
	}
}

// GetDiscussions 獲取所有討論
// @Summary 獲取所有討論
// @Description 獲取所有討論的列表
// @Tags discussions
// @Accept json
// @Produce json
// @Param type query string false "討論類型"
// @Param tags query string false "標籤"
// @Param keyword query string false "關鍵字"
// @Param most_viewed query string false "是否按瀏覽量排序"
// @Param limit query string false "限制數量"
// @Success 200 {array} models.Discussion
// @Failure 500 {object} map[string]string
// @Router /discussions [get]
func (c *DiscussionController) GetDiscussions(ctx *gin.Context) {
	discussionType := ctx.Query("type")
	tagsStr := ctx.Query("tags")
	keyword := ctx.Query("keyword")
	mostViewedStr := ctx.Query("most_viewed")
	limitStr := ctx.Query("limit")

	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
	}

	limit := 0
	if mostViewedStr == "true" {
		limit = 10
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}
	}

	discussions, err := c.discussionService.GetDiscussionsWithFilters(discussionType, tags, keyword, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, discussions)
}

// GetDiscussion 根據ID獲取討論
// @Summary 獲取特定討論
// @Description 根據ID獲取特定討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Success 200 {object} models.Discussion
// @Failure 400,404 {object} map[string]string
// @Router /discussions/{id} [get]
func (c *DiscussionController) GetDiscussion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	discussion, err := c.discussionService.GetDiscussionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "discussion not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, discussion)
}

// GetUserDiscussions 獲取特定用戶的所有討論
// @Summary 獲取特定用戶的所有討論
// @Description 獲取特定用戶的所有討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Discussion
// @Failure 400,404 {object} map[string]string
// @Router /discussions/user/{user_id} [get]
func (c *DiscussionController) GetUserDiscussions(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	discussions, err := c.discussionService.GetUserDiscussions(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, discussions)
}

// GetChildDiscussions 獲取子討論
// @Summary 獲取子討論
// @Description 根據父討論 ID 獲取子討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param parent_id path int true "Parent Discussion ID"
// @Param parent_type path string true "Parent Type"
// @Success 200 {array} models.Discussion
// @Failure 400 {object} map[string]string
// @Router /discussions/parent/{parent_type}/{parent_id} [get]
func (c *DiscussionController) GetChildDiscussions(ctx *gin.Context) {
	parentID, err := strconv.Atoi(ctx.Param("parent_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parent ID",
		})
		return
	}

	parentType := ctx.Param("parent_type")
	if parentType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "parent type is required",
		})
		return
	}

	discussions, err := c.discussionService.GetChildDiscussions(parentID, parentType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, discussions)
}

// CreateDiscussion 創建新討論
// @Summary 創建新討論
// @Description 創建新的討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param discussion body models.Discussion true "Discussion"
// @Success 201 {object} models.Discussion
// @Failure 400 {object} map[string]string
// @Router /discussions [post]
func (c *DiscussionController) CreateDiscussion(ctx *gin.Context) {
	var discussion models.Discussion

	if err := ctx.ShouldBindJSON(&discussion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if exists {
		discussion.AuthorID = userID.(int)
	}

	createdDiscussion, err := c.discussionService.CreateDiscussion(discussion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdDiscussion)
}

// UpdateDiscussion 更新討論
// @Summary 更新討論
// @Description 更新指定ID的討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Param discussion body models.Discussion true "Discussion"
// @Success 200 {object} models.Discussion
// @Failure 400,404 {object} map[string]string
// @Router /discussions/{id} [put]
func (c *DiscussionController) UpdateDiscussion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var discussionUpdate models.Discussion
	if err := ctx.ShouldBindJSON(&discussionUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.discussionService.UpdateDiscussion(uint(id), discussionUpdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "discussion updated successfully",
	})
}

// DeleteDiscussion 刪除討論
// @Summary 刪除討論
// @Description 根據ID刪除討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param id path int true "Discussion ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /discussions/{id} [delete]
func (c *DiscussionController) DeleteDiscussion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	// 可能需要檢查權限，確保只有作者或管理員可以刪除討論

	err = c.discussionService.DeleteDiscussion(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "discussion deleted successfully",
	})
}

// SearchDiscussions 搜索討論
// @Summary 搜索討論
// @Description 根據關鍵字搜索討論
// @Tags discussions
// @Accept json
// @Produce json
// @Param keyword query string true "Search Keyword"
// @Success 200 {array} models.Discussion
// @Failure 400 {object} map[string]string
// @Router /discussions/search [get]
func (c *DiscussionController) SearchDiscussions(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "search keyword is required",
		})
		return
	}

	discussions, err := c.discussionService.SearchDiscussions(keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, discussions)
}
