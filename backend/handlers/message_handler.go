package handlers

import (
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// SEND MESSAGE
func (h *MessageHandler) SendMessage(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	var body struct {
		ReceiverID string `form:"receiverId"`
		Content    string `form:"content"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	file, _ := c.FormFile("image")

	var filePath string
	if file != nil {
		filePath = "./uploads/messages/" + file.Filename

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
			return
		}
	}

	msg, err := h.messageService.SendMessage(
		c.Request.Context(),
		payload.UserID,
		body.ReceiverID,
		body.Content,
		filePath,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *MessageHandler) GetChatMessages(c *gin.Context) {

	userData, _ := c.Get("user")
	payload := userData.(*utils.TokenPayload)

	receiverId := c.Param("receiverId")

	messages, err := h.messageService.GetChatMessages(
		c.Request.Context(),
		payload.UserID,
		receiverId,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// DELETE MESSAGE
func (h *MessageHandler) DeleteMessage(c *gin.Context) {

	userData, _ := c.Get("user")
	payload := userData.(*utils.TokenPayload)

	messageId := c.Param("messageId")

	res, err := h.messageService.DeleteMessage(
		c.Request.Context(),
		payload.UserID,
		messageId,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// MARK AS READ
func (h *MessageHandler) MarkMessageAsRead(c *gin.Context) {

	userData, _ := c.Get("user")
	payload := userData.(*utils.TokenPayload)

	senderId := c.Param("senderId")

	res, err := h.messageService.MarkAsRead( 
		c.Request.Context(),
		payload.UserID,
		senderId,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
