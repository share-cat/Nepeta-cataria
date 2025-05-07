package controllers

import (
	"net/http"
	"strconv"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type ChatSessionController struct {
	chatSessionService *services.ChatSessionService
}

func NewChatSessionController(chatSessionService *services.ChatSessionService) *ChatSessionController {
	return &ChatSessionController{
		chatSessionService: chatSessionService,
	}
}

// GetChatSessions 獲取所有聊天會話
// @Summary 獲取所有聊天會話
// @Description 獲取所有聊天會話的列表
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Success 200 {array} models.ChatSession
// @Failure 500 {object} map[string]string
// @Router /chat_sessions [get]
func (c *ChatSessionController) GetChatSessions(ctx *gin.Context) {
	sessions, err := c.chatSessionService.GetAllChatSessions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}

// GetChatSession 根據ID獲取聊天會話
// @Summary 獲取特定聊天會話
// @Description 根據ID獲取特定聊天會話
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param id path int true "Chat Session ID"
// @Success 200 {object} models.ChatSession
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/{id} [get]
func (c *ChatSessionController) GetChatSession(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	session, err := c.chatSessionService.GetChatSessionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "chat session not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

// GetUserChatSessions 獲取特定用戶的所有聊天會話
// @Summary 獲取特定用戶的所有聊天會話
// @Description 獲取特定用戶的所有聊天會話
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.ChatSession
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/user/{user_id} [get]
func (c *ChatSessionController) GetUserChatSessions(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	sessions, err := c.chatSessionService.GetUserChatSessions(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}

// GetJudgeChatSessions 獲取特定評審的所有聊天會話
// @Summary 獲取特定評審的所有聊天會話
// @Description 獲取特定評審的所有聊天會話
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param judge_id path int true "Judge ID"
// @Success 200 {array} models.ChatSession
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/judge/{judge_id} [get]
func (c *ChatSessionController) GetJudgeChatSessions(ctx *gin.Context) {
	judgeID, err := strconv.ParseUint(ctx.Param("judge_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid judge ID",
		})
		return
	}

	sessions, err := c.chatSessionService.GetJudgeChatSessions(uint(judgeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}

// CreateChatSession 創建新聊天會話
// @Summary 創建新聊天會話
// @Description 創建新的聊天會話
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param chat_session body models.ChatSession true "Chat Session"
// @Success 201 {object} models.ChatSession
// @Failure 400 {object} map[string]string
// @Router /chat-sessions [post]
func (c *ChatSessionController) CreateChatSession(ctx *gin.Context) {
	var session models.ChatSession

	if err := ctx.ShouldBindJSON(&session); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdSession, err := c.chatSessionService.CreateChatSession(session)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdSession)
}

// UpdateChatSession 更新聊天會話
// @Summary 更新聊天會話
// @Description 更新聊天會話的詳細信息
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param id path int true "Chat Session ID"
// @Param chat_session body models.ChatSession true "Chat Session"
// @Success 200 {object} models.ChatSession
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/{id} [put]
func (c *ChatSessionController) UpdateChatSession(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var sessionUpdate models.ChatSession
	if err := ctx.ShouldBindJSON(&sessionUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.chatSessionService.UpdateChatSession(uint(id), sessionUpdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "chat session updated successfully",
	})
}

// UpdateChatHistory 更新聊天歷史
// @Summary 更新聊天歷史
// @Description 更新聊天會話的歷史記錄
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param id path int true "Chat Session ID"
// @Param chat_history body interface{} true "Chat History"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/{id}/history [put]
func (c *ChatSessionController) UpdateChatHistory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var request struct {
		ChatHistory interface{} `json:"chat_history" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.chatSessionService.UpdateChatHistory(uint(id), request.ChatHistory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "chat history updated successfully",
	})
}

// UpdateTitle 更新標題
// @Summary 更新標題
// @Description 更新聊天會話的標題
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param id path int true "Chat Session ID"
// @Param title body string true "Title"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/{id}/title [put]
func (c *ChatSessionController) UpdateTitle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var request struct {
		Title string `json:"title" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.chatSessionService.UpdateTitle(uint(id), request.Title)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "title updated successfully",
	})
}

// DeleteChatSession 刪除聊天會話
// @Summary 刪除聊天會話
// @Description 根據ID刪除聊天會話
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param id path int true "Chat Session ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/{id} [delete]
func (c *ChatSessionController) DeleteChatSession(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	err = c.chatSessionService.DeleteChatSession(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "chat session deleted successfully",
	})
}

// GetUserJudgeChatSessions 獲取特定用戶和特定評審的聊天會話
// @Summary 獲取特定用戶和特定評審的聊天會話
// @Description 獲取特定用戶和特定評審的聊天會話
// @Tags chat_sessions
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param judge_id path int true "Judge ID"
// @Success 200 {array} models.ChatSession
// @Failure 400,404 {object} map[string]string
// @Router /chat-sessions/user/{user_id}/judge/{judge_id} [get]
func (c *ChatSessionController) GetUserJudgeChatSessions(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	judgeID, err := strconv.ParseUint(ctx.Param("judge_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid judge ID",
		})
		return
	}

	sessions, err := c.chatSessionService.GetUserJudgeChatSessions(uint(userID), uint(judgeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}
