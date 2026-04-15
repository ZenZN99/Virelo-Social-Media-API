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

type CommentService struct {
	commentModel      *mongo.Collection
	notificationModel *mongo.Collection
	postModel         *mongo.Collection
	reelModel         *mongo.Collection
	commentModelRef   *mongo.Collection
}

func NewCommentService(
	commentModel *mongo.Collection,
	notificationModel *mongo.Collection,
	postModel *mongo.Collection,
	reelModel *mongo.Collection,
) *CommentService {
	return &CommentService{
		commentModel:      commentModel,
		notificationModel: notificationModel,
		postModel:         postModel,
		reelModel:         reelModel,
		commentModelRef:   commentModel,
	}
}

func (s *CommentService) CreateComment(
	ctx context.Context,
	userID string,
	targetID string,
	targetType models.ContentType,
	text string,
	parentCommentID string,
) (models.CommentModel, error) {

	if text == "" {
		return models.CommentModel{}, errors.New("text is required")
	}

	if !utils.IsValidTargetType(targetType) {
		return models.CommentModel{}, errors.New("invalid target type")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.CommentModel{}, errors.New("invalid user id")
	}

	targetObjectID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return models.CommentModel{}, errors.New("invalid target id")
	}

	var parentID *primitive.ObjectID
	if parentCommentID != "" {
		id, err := primitive.ObjectIDFromHex(parentCommentID)
		if err == nil {
			parentID = &id
		}
	}

	comment := models.CommentModel{
		ID:              primitive.NewObjectID(),
		UserID:          userObjectID,
		TargetID:        targetObjectID,
		TargetType:      models.ContentType(targetType),
		Text:            text,
		ParentCommentId: parentID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = s.commentModel.InsertOne(ctx, comment)
	if err != nil {
		return models.CommentModel{}, err
	}

	ownerID, err := s.getTargetOwner(ctx, targetObjectID, targetType)
	if err != nil {
		return models.CommentModel{}, err
	}

	if ownerID.Hex() != userObjectID.Hex() {

		notification := models.NotificationModel{
			ID:               primitive.NewObjectID(),
			SenderID:         userObjectID,
			ReceiverID:       ownerID,
			NotificationType: models.COMMENT_,
			TargetID:         comment.ID,
			IsRead:           false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		_, err = s.notificationModel.InsertOne(ctx, notification)
		if err != nil {
			return models.CommentModel{}, err
		}
	}

	return comment, nil
}
func (s *CommentService) UpdateComment(
	ctx context.Context,
	commentID string,
	user models.UserModel,
	text string,
) (models.CommentModel, error) {

	cID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return models.CommentModel{}, errors.New("invalid comment id")
	}

	var comment models.CommentModel

	err = s.commentModel.FindOne(ctx, bson.M{"_id": cID}).Decode(&comment)
	if err != nil {
		return models.CommentModel{}, errors.New("comment not found")
	}

	isOwner := comment.UserID.Hex() == user.ID.Hex()
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		return models.CommentModel{}, errors.New("not allowed")
	}

	update := bson.M{
		"text":      text,
		"updatedAt": time.Now(),
	}

	_, err = s.commentModel.UpdateOne(
		ctx,
		bson.M{"_id": cID},
		bson.M{"$set": update},
	)

	if err != nil {
		return models.CommentModel{}, err
	}

	comment.Text = text
	comment.UpdatedAt = time.Now()

	return comment, nil
}

func (s *CommentService) DeleteComment(
	ctx context.Context,
	commentID string,
	user models.UserModel,
) (models.CommentModel, error) {

	cID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return models.CommentModel{}, errors.New("invalid comment id")
	}

	var comment models.CommentModel

	err = s.commentModel.FindOne(ctx, bson.M{"_id": cID}).Decode(&comment)
	if err != nil {
		return models.CommentModel{}, errors.New("comment not found")
	}

	isOwner := comment.UserID.Hex() == user.ID.Hex()
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		return models.CommentModel{}, errors.New("not allowed")
	}

	_, err = s.commentModel.DeleteOne(ctx, bson.M{"_id": cID})
	return models.CommentModel{}, err
}

func (s *CommentService) GetComments(
	ctx context.Context,
	targetID string,
	targetType models.ContentType,
	page int,
	limit int,
) ([]models.CommentModel, error) {

	tID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return nil, errors.New("invalid target id")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	filter := bson.M{
		"targetId":   tID,
		"targetType": targetType,
		"parentCommentId": bson.M{
			"$in": []interface{}{"", nil},
		},
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.commentModel.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.CommentModel

	for cursor.Next(ctx) {
		var comment models.CommentModel
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) GetReplies(
	ctx context.Context,
	parentCommentID string,
	page int,
	limit int,
) ([]models.CommentModel, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	filter := bson.M{
		"parentCommentId": parentCommentID,
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.commentModel.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var replies []models.CommentModel

	for cursor.Next(ctx) {
		var reply models.CommentModel
		if err := cursor.Decode(&reply); err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return replies, nil
}

func (s *CommentService) GetCommentsByUser(
	ctx context.Context,
	userID string,
	page int,
	limit int,
) ([]models.CommentModel, error) {

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

	filter := bson.M{
		"userId": objectID,
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.commentModel.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.CommentModel

	for cursor.Next(ctx) {
		var comment models.CommentModel
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) getTargetOwner(
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
