package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService struct {
	messageModel      *mongo.Collection
	cloudinaryService *utils.CloudinaryService
}

func NewMessageService(messageModel *mongo.Collection, cloudinaryService *utils.CloudinaryService) *MessageService {
	return &MessageService{
		messageModel:      messageModel,
		cloudinaryService: cloudinaryService,
	}
}

func (s *MessageService) SendMessage(
	ctx context.Context,
	senderID string,
	receiverID string,
	content string,
	imagePath string,
) (models.MessageModel, error) {

	if receiverID == "" {
		return models.MessageModel{}, errors.New("receiver id is required")
	}

	if content == "" && imagePath == "" {
		return models.MessageModel{}, errors.New("message must contain text or image")
	}

	senderObj, err := primitive.ObjectIDFromHex(senderID)
	if err != nil {
		return models.MessageModel{}, errors.New("invalid sender id")
	}

	receiverObj, err := primitive.ObjectIDFromHex(receiverID)
	if err != nil {
		return models.MessageModel{}, errors.New("invalid receiver id")
	}

	imageURL := ""

	// upload image if exists
	if imagePath != "" {
		res, err := s.cloudinaryService.UploadFile(
			ctx,
			imagePath,
			"messages",
			"image/jpeg",
		)

		if err != nil {
			return models.MessageModel{}, err
		}

		imageURL = res.URL
	}

	message := models.MessageModel{
		ID:         primitive.NewObjectID(),
		SenderID:   senderObj,
		ReceiverID: receiverObj,
		Content:    content,
		Image:      imageURL,
		IsRead:     false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = s.messageModel.InsertOne(ctx, message)
	if err != nil {
		return models.MessageModel{}, err
	}

	return message, nil
}

func (s *MessageService) GetChatMessages(
	ctx context.Context,
	userID string,
	receiverID string,
) ([]models.MessageModel, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	receiverObj, err := primitive.ObjectIDFromHex(receiverID)
	if err != nil {
		return nil, errors.New("invalid receiver id")
	}

	filter := bson.M{
		"$or": []bson.M{
			{
				"senderId":   userObj,
				"receiverId": receiverObj,
			},
			{
				"senderId":   receiverObj,
				"receiverId": userObj,
			},
		},
	}

	cursor, err := s.messageModel.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.MessageModel

	for cursor.Next(ctx) {
		var msg models.MessageModel
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *MessageService) DeleteMessage(
	ctx context.Context,
	userID string,
	messageID string,
) (map[string]interface{}, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	msgObj, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, errors.New("invalid message id")
	}

	res, err := s.messageModel.DeleteOne(ctx, bson.M{
		"_id":      msgObj,
		"senderId": userObj,
	})

	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("message not found")
	}

	return map[string]interface{}{
		"success": "message deleted successfully",
	}, nil
}

func (s *MessageService) MarkAsRead(
	ctx context.Context,
	userID string,
	senderID string,
) (map[string]interface{}, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	senderObj, err := primitive.ObjectIDFromHex(senderID)
	if err != nil {
		return nil, errors.New("invalid sender id")
	}

	_, err = s.messageModel.UpdateMany(
		ctx,
		bson.M{
			"senderId":   senderObj,
			"receiverId": userObj,
			"isRead":     false,
		},
		bson.M{
			"$set": bson.M{"isRead": true},
		},
	)

	return nil, err
}
