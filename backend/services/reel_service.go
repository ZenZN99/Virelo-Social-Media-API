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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReelService struct {
	reelModel         *mongo.Collection
	followModel       *mongo.Collection
	notificationModel *mongo.Collection
	cloudinaryService *utils.CloudinaryService
}

func NewReelService(
	reelModel *mongo.Collection,
	followModel *mongo.Collection,
	notificationModel *mongo.Collection,
	cloudinaryService *utils.CloudinaryService,
) *ReelService {
	return &ReelService{
		reelModel:         reelModel,
		followModel:       followModel,
		notificationModel: notificationModel,
		cloudinaryService: cloudinaryService,
	}
}
func (s *ReelService) CreateReel(
	ctx context.Context,
	userID string,
	content string,
	filePath string,
) (models.ReelModel, error) {

	if filePath == "" {
		return models.ReelModel{}, errors.New("reel file is required")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.ReelModel{}, errors.New("invalid user id")
	}

	uploadResult, err := s.cloudinaryService.UploadFile(
		ctx,
		filePath,
		"reels",
		"video/mp4",
	)
	if err != nil {
		return models.ReelModel{}, err
	}

	reel := models.ReelModel{
		ID:        primitive.NewObjectID(),
		UserID:    objectID,
		Reel:      uploadResult.URL,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.reelModel.InsertOne(ctx, reel)
	if err != nil {
		return models.ReelModel{}, err
	}

	followers, err := s.followModel.Find(ctx, bson.M{
		"following": objectID,
	})
	if err != nil {
		return models.ReelModel{}, err
	}
	defer followers.Close(ctx)

	for followers.Next(ctx) {

		var follow models.FollowModel
		if err := followers.Decode(&follow); err != nil {
			return models.ReelModel{}, err
		}

		notification := models.NotificationModel{
			ID:               primitive.NewObjectID(),
			SenderID:         objectID,
			ReceiverID:       follow.Follower,
			NotificationType: models.REEL_,
			TargetID:         reel.ID,
			IsRead:           false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		_, err := s.notificationModel.InsertOne(ctx, notification)
		if err != nil {
			return models.ReelModel{}, err
		}

		r := utils.GetRedis()
		_ = r.Incr(ctx, "unread_count:"+follow.Follower.Hex()).Err()
		_ = r.LPush(ctx, "notifications:"+follow.Follower.Hex(), notification.ID.Hex()).Err()
	}

	return reel, nil
}

func (s *ReelService) UpdateReel(
	ctx context.Context,
	reelID string,
	user *utils.TokenPayload,
	content string,
) (models.ReelModel, error) {

	objectID, err := primitive.ObjectIDFromHex(reelID)
	if err != nil {
		return models.ReelModel{}, errors.New("invalid reel id")
	}

	var reel models.ReelModel
	err = s.reelModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&reel)
	if err != nil {
		return models.ReelModel{}, errors.New("reel not found")
	}

	isOwner := reel.UserID.Hex() == user.UserID
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		return models.ReelModel{}, errors.New("not allowed")
	}

	update := bson.M{}

	if content != "" {
		update["content"] = content
	}

	update["updatedAt"] = time.Now()

	_, err = s.reelModel.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)

	if err != nil {
		return models.ReelModel{}, err
	}

	var updatedReel models.ReelModel
	err = s.reelModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedReel)
	if err != nil {
		return models.ReelModel{}, err
	}

	return updatedReel, nil
}

func (s *ReelService) DeleteReel(
	ctx context.Context,
	reelID string,
	user *utils.TokenPayload,
) (map[string]interface{}, error) {

	objectID, err := primitive.ObjectIDFromHex(reelID)
	if err != nil {
		return nil, errors.New("invalid reel id")
	}

	var reel models.ReelModel
	err = s.reelModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&reel)
	if err != nil {
		return nil, errors.New("reel not found")
	}

	isOwner := reel.UserID.Hex() == user.UserID
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		return nil, errors.New("not allowed")
	}

	_, err = s.reelModel.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Reel deleted successfully",
	}, nil
}

func (s *ReelService) GetAllReels(
	ctx context.Context,
	page int,
	limit int,
) ([]models.ReelModel, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.reelModel.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reels []models.ReelModel

	for cursor.Next(ctx) {
		var reel models.ReelModel
		if err := cursor.Decode(&reel); err != nil {
			return nil, err
		}
		reels = append(reels, reel)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return reels, nil
}

func (s *ReelService) GetReelById(
	ctx context.Context,
	reelID string,
) (models.ReelModel, error) {

	objectID, err := primitive.ObjectIDFromHex(reelID)
	if err != nil {
		return models.ReelModel{}, errors.New("invalid reel id")
	}

	var reel models.ReelModel

	err = s.reelModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&reel)
	if err != nil {
		return models.ReelModel{}, errors.New("reel not found")
	}

	return reel, nil
}

func (s *ReelService) GetReelsByUser(
	ctx context.Context,
	userID string,
	page int,
	limit int,
) ([]models.ReelModel, error) {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	filter := bson.M{"userId": objectID}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.reelModel.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reels []models.ReelModel

	for cursor.Next(ctx) {
		var reel models.ReelModel
		if err := cursor.Decode(&reel); err != nil {
			return nil, err
		}
		reels = append(reels, reel)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return reels, nil
}
