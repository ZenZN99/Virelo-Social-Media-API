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

type NotificationService struct {
	notificationModel *mongo.Collection
}

func NewNotificationService(
	notificationModel *mongo.Collection,
) *NotificationService {
	return &NotificationService{
		notificationModel: notificationModel,
	}
}

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	data models.NotificationModel,
) error {

	if data.ReceiverID == data.SenderID {
		return nil
	}

	data.ID = primitive.NewObjectID()
	data.IsRead = false
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	_, err := s.notificationModel.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	r := utils.GetRedis()

	receiverID := data.ReceiverID.Hex()

	_ = r.Incr(ctx, "unread_count:"+receiverID).Err()

	_ = r.LPush(ctx,
		"notifications:"+receiverID,
		data.ID.Hex(),
	).Err()

	_ = r.LTrim(ctx,
		"notifications:"+receiverID,
		0,
		19,
	).Err()

	return nil
}

func (s *NotificationService) GetUserNotifications(
	ctx context.Context,
	userId string,
) ([]models.NotificationModel, error) {

	uID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	cursor, err := s.notificationModel.Find(ctx, bson.M{
		"receiverId": uID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []models.NotificationModel

	for cursor.Next(ctx) {
		var n models.NotificationModel
		if err := cursor.Decode(&n); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (s *NotificationService) GetUnreadCount(
	ctx context.Context,
	userId string,
) (int64, error) {

	r := utils.GetRedis()

	key := "unread_count:" + userId

	_, err := r.Get(ctx, key).Result()
	if err == nil {
		return r.Get(ctx, key).Int64()
	}

	count, err := s.notificationModel.CountDocuments(ctx, bson.M{
		"receiverId": userId,
		"isRead":     false,
	})
	if err != nil {
		return 0, err
	}

	_ = r.Set(ctx, key, count, 0).Err()

	return count, nil
}

func (s *NotificationService) MarkAsRead(
	ctx context.Context,
	notificationId string,
) (map[string]interface{}, error) {

	nID, err := primitive.ObjectIDFromHex(notificationId)
	if err != nil {
		return nil, errors.New("invalid notification id")
	}

	var notification models.NotificationModel

	err = s.notificationModel.FindOneAndUpdate(
		ctx,
		bson.M{"_id": nID},
		bson.M{"$set": bson.M{"isRead": true}},
	).Decode(&notification)

	if err != nil {
		return nil, err
	}

	r := utils.GetRedis()
	_ = r.Decr(ctx, "unread_count:"+notification.ReceiverID.Hex()).Err()

	return nil, nil
}

func (s *NotificationService) MarkAllAsRead(
	ctx context.Context,
	userId string,
) (map[string]interface{}, error) {

	uID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	_, err = s.notificationModel.UpdateMany(
		ctx,
		bson.M{
			"receiverId": uID,
			"isRead":     false,
		},
		bson.M{
			"$set": bson.M{"isRead": true},
		},
	)

	if err != nil {
		return nil, err
	}

	r := utils.GetRedis()
	_ = r.Set(ctx, "unread_count:"+userId, 0, 0).Err()

	return nil, nil
}

func (s *NotificationService) DeleteNotification(
	ctx context.Context,
	notificationId string,
) (map[string]interface{}, error) {

	nID, err := primitive.ObjectIDFromHex(notificationId)
	if err != nil {
		return nil, errors.New("invalid notification id")
	}

	var notification models.NotificationModel

	err = s.notificationModel.FindOne(ctx, bson.M{
		"_id": nID,
	}).Decode(&notification)

	if err != nil {
		return nil, errors.New("notification not found")
	}

	_, err = s.notificationModel.DeleteOne(ctx, bson.M{
		"_id": nID,
	})
	if err != nil {
		return nil, err
	}

	r := utils.GetRedis()

	_ = r.LRem(
		ctx,
		"notifications:"+notification.ReceiverID.Hex(),
		0,
		notification.ID.Hex(),
	).Err()

	if !notification.IsRead {
		_ = r.Decr(ctx, "unread_count:"+notification.ReceiverID.Hex()).Err()
	}

	return nil, nil
}
