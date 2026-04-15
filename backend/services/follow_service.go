package services

import (
	"backend/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowService struct {
	followModel       *mongo.Collection
	notificationModel *mongo.Collection
}

func NewFollowService(followModel *mongo.Collection, notificationModel *mongo.Collection) *FollowService {
	return &FollowService{
		followModel:       followModel,
		notificationModel: notificationModel,
	}
}

func (s *FollowService) FollowUser(
	ctx context.Context,
	followerID string,
	followingID string,
) (models.FollowModel, error) {

	if followerID == followingID {
		return models.FollowModel{}, errors.New("you can't follow yourself")
	}

	followerObj, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return models.FollowModel{}, errors.New("invalid follower id")
	}

	followingObj, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return models.FollowModel{}, errors.New("invalid following id")
	}

	count, err := s.followModel.CountDocuments(ctx, bson.M{
		"follower":  followerObj,
		"following": followingObj,
	})

	if err != nil {
		return models.FollowModel{}, err
	}

	if count > 0 {
		return models.FollowModel{}, errors.New("already following this user")
	}

	follow := models.FollowModel{
		ID:          primitive.NewObjectID(),
		Follower:    followerObj,
		Following:   followingObj,
		IsFollowing: true,
	}

	_, err = s.followModel.InsertOne(ctx, follow)
	if err != nil {
		return models.FollowModel{}, err
	}

	notification := models.NotificationModel{
		ID:               primitive.NewObjectID(),
		SenderID:         followerObj,
		ReceiverID:       followingObj,
		NotificationType: models.FOLLOW_,
		TargetID:         follow.ID,
		IsRead:           false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, _ = s.notificationModel.InsertOne(ctx, notification)

	return follow, nil
}
func (s *FollowService) UnfollowUser(
	ctx context.Context,
	followerID string,
	followingID string,
) error {

	followerObj, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return errors.New("invalid follower id")
	}

	followingObj, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return errors.New("invalid following id")
	}

	res, err := s.followModel.DeleteOne(ctx, bson.M{
		"follower":    followerObj,
		"following":   followingObj,
		"isFollowing": false,
	})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("follow not found")
	}

	return nil
}

func (s *FollowService) IsFollowing(
	ctx context.Context,
	followerID string,
	followingID string,
) (bool, error) {

	followerObj, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return false, errors.New("invalid follower id")
	}

	followingObj, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return false, errors.New("invalid following id")
	}

	err = s.followModel.FindOne(ctx, bson.M{
		"follower":  followerObj,
		"following": followingObj,
	}).Err()

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *FollowService) GetFollowers(
	ctx context.Context,
	userID string,
) ([]bson.M, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"following": userObj,
		}}},
	}

	cursor, err := s.followModel.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M

	for cursor.Next(ctx) {
		var item bson.M
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *FollowService) GetFollowing(
	ctx context.Context,
	userID string,
) ([]bson.M, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	cursor, err := s.followModel.Find(ctx, bson.M{
		"follower": userObj,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M

	for cursor.Next(ctx) {
		var item bson.M
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *FollowService) GetFollowersCount(
	ctx context.Context,
	userID string,
) (int64, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, errors.New("invalid user id")
	}

	return s.followModel.CountDocuments(ctx, bson.M{
		"following": userObj,
	})
}

func (s *FollowService) GetFollowingCount(
	ctx context.Context,
	userID string,
) (int64, error) {

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, errors.New("invalid user id")
	}

	return s.followModel.CountDocuments(ctx, bson.M{
		"follower": userObj,
	})
}
