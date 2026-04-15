package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	var body struct {
		TargetID        string            `json:"targetId"`
		TargetType      models.ContentType `json:"targetType"`
		Text            string            `json:"text"`
		ParentCommentID string            `json:"parentCommentId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	tType := models.ContentType(strings.ToLower(string(body.TargetType)))

	if !utils.IsValidTargetType(tType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target type"})
		return
	}

	comment, err := h.commentService.CreateComment(
		c.Request.Context(),
		payload.UserID,
		body.TargetID,
		tType,
		body.Text,
		body.ParentCommentID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	commentID := c.Param("commentId")

	var body struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	objID, _ := primitive.ObjectIDFromHex(payload.UserID)

	user := models.UserModel{
		ID:   objID,
		Role: payload.Role,
	}

	comment, err := h.commentService.UpdateComment(
		c.Request.Context(),
		commentID,
		user,
		body.Text,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	commentID := c.Param("commentId")

	objID, _ := primitive.ObjectIDFromHex(payload.UserID)

	user := models.UserModel{
		ID:   objID,
		Role: payload.Role,
	}

	_, err := h.commentService.DeleteComment(
		c.Request.Context(),
		commentID,
		user,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})
}

func (h *CommentHandler) GetComments(c *gin.Context) {

	targetID := c.Query("targetId")
	targetType := models.ContentType(c.Query("targetType"))

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	comments, err := h.commentService.GetComments(
		c.Request.Context(),
		targetID,
		targetType,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) GetReplies(c *gin.Context) {

	parentID := c.Param("commentId")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	replies, err := h.commentService.GetReplies(
		c.Request.Context(),
		parentID,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replies)
}

func (h *CommentHandler) GetCommentsByUser(c *gin.Context) {

	userID := c.Param("userId")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	comments, err := h.commentService.GetCommentsByUser(
		c.Request.Context(),
		userID,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
