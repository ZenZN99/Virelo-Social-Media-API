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

type LikeService struct {
	likeModel         *mongo.Collection
	postModel         *mongo.Collection
	reelModel         *mongo.Collection
	commentModel      *mongo.Collection
	notificationModel *mongo.Collection
}

func NewLikeService(
	likeModel *mongo.Collection,
	postModel *mongo.Collection,
	reelModel *mongo.Collection,
	commentModel *mongo.Collection,
	notificationModel *mongo.Collection,
) *LikeService {
	return &LikeService{
		likeModel:         likeModel,
		postModel:         postModel,
		reelModel:         reelModel,
		commentModel:      commentModel,
		notificationModel: notificationModel,
	}
}

func (s *LikeService) ToggleLike(
	ctx context.Context,
	userID string,
	targetID string,
	targetType models.ContentType,
) (map[string]interface{}, error) {

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	targetObjectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return nil, errors.New("invalid target id")
	}

	if !utils.IsValidLikeType(targetType) {
		return nil, errors.New("invalid target type")
	}

	filter := bson.M{
		"userId":     userObjectID,
		"targetId":   targetObjectID,
		"targetType": targetType,
	}

	var existing models.LikeModel
	err = s.likeModel.FindOne(ctx, filter).Decode(&existing)

	if err == nil {
		_, err := s.likeModel.DeleteOne(ctx, filter)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"liked": false,
		}, nil
	}

	like := models.LikeModel{
		ID:         primitive.NewObjectID(),
		UserID:     userObjectID,
		TargetID:   targetObjectID,
		TargetType: targetType,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = s.likeModel.InsertOne(ctx, like)
	if err != nil {
		return nil, err
	}

	ownerID, err := s.getTargetOwner(ctx, targetObjectID, targetType)
	if err != nil {
		return nil, err
	}

	if ownerID != userObjectID {

		notification := models.NotificationModel{
			ID:               primitive.NewObjectID(),
			SenderID:         userObjectID,
			ReceiverID:       ownerID,
			NotificationType: models.LIKE_,
			TargetID:         targetObjectID,
			IsRead:           false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		_, err = s.notificationModel.InsertOne(ctx, notification)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"liked": true,
		"like":  like,
	}, nil
}

func (s *LikeService) CountLikes(
	ctx context.Context,
	targetID string,
	targetType models.ContentType,
) (int64, error) {

	objectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return 0, errors.New("invalid target id")
	}

	count, err := s.likeModel.CountDocuments(ctx, bson.M{
		"targetId":   objectID,
		"targetType": targetType,
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *LikeService) IsLiked(
	ctx context.Context,
	userID string,
	targetID string,
	targetType models.ContentType,
) (bool, error) {

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user id")
	}

	targetObjectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return false, errors.New("invalid target id")
	}

	err = s.likeModel.FindOne(ctx, bson.M{
		"userId":     userObjectID,
		"targetId":   targetObjectID,
		"targetType": targetType,
	}).Err()

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *LikeService) GetLikes(
	ctx context.Context,
	targetID string,
	targetType models.ContentType,
) ([]bson.M, error) {

	objectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return nil, errors.New("invalid target id")
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"targetId":   objectID,
			"targetType": targetType,
		}}},
		{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
	}

	cursor, err := s.likeModel.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var likes []bson.M

	for cursor.Next(ctx) {
		var like bson.M
		if err := cursor.Decode(&like); err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func (s *LikeService) getTargetOwner(
	ctx context.Context,
	targetID primitive.ObjectID,
	targetType models.ContentType,
) (primitive.ObjectID, error) {

	switch targetType {

	case models.Post:
		var post struct {
			UserID primitive.ObjectID `bson:"userId"`
		}

		err := s.postModel.FindOne(ctx, bson.M{"_id": targetID}).Decode(&post)
		if err != nil {
			return primitive.NilObjectID, errors.New("post not found")
		}

		return post.UserID, nil

	case models.Reel:
		var reel struct {
			UserID primitive.ObjectID `bson:"userId"`
		}

		err := s.reelModel.FindOne(ctx, bson.M{"_id": targetID}).Decode(&reel)
		if err != nil {
			return primitive.NilObjectID, errors.New("reel not found")
		}

		return reel.UserID, nil

	case models.Comment:
		var comment struct {
			UserID primitive.ObjectID `bson:"userId"`
		}

		err := s.commentModel.FindOne(ctx, bson.M{"_id": targetID}).Decode(&comment)
		if err != nil {
			return primitive.NilObjectID, errors.New("comment not found")
		}

		return comment.UserID, nil

	default:
		return primitive.NilObjectID, errors.New("invalid target type")
	}
}
