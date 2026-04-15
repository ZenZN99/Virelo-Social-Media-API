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

type PostService struct {
	postModel         *mongo.Collection
	followModel       *mongo.Collection
	notificationModel *mongo.Collection
	cloudinaryService *utils.CloudinaryService
}

func NewPostService(
	postModel *mongo.Collection,
	followModel *mongo.Collection,
	notificationModel *mongo.Collection,
	cloudinaryService *utils.CloudinaryService,
) *PostService {
	return &PostService{
		postModel:         postModel,
		followModel:       followModel,
		notificationModel: notificationModel,
		cloudinaryService: cloudinaryService,
	}
}

func (s *PostService) CreatePost(
	ctx context.Context,
	userID string,
	content string,
	files []string,
) (models.PostModel, error) {

	if len(files) == 0 {
		return models.PostModel{}, errors.New("images are required")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.PostModel{}, errors.New("invalid user id")
	}

	var uploadedImages []models.PostImage

	for _, file := range files {

		res, err := s.cloudinaryService.UploadFile(
			ctx,
			file,
			"posts",
			"image/jpeg",
		)

		if err != nil {
			return models.PostModel{}, err
		}

		uploadedImages = append(uploadedImages, models.PostImage{
			URL:      res.URL,
			PublicID: res.PublicID,
		})
	}

	post := models.PostModel{
		ID:        primitive.NewObjectID(),
		Content:   content,
		Images:    uploadedImages,
		UserID:    objectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.postModel.InsertOne(ctx, post)
	if err != nil {
		return models.PostModel{}, err
	}

	followers, err := s.followModel.Find(ctx, bson.M{
		"following": objectID,
	})
	if err != nil {
		return models.PostModel{}, err
	}
	defer followers.Close(ctx)

	for followers.Next(ctx) {

		var follow models.FollowModel
		if err := followers.Decode(&follow); err != nil {
			return models.PostModel{}, err
		}

		notification := models.NotificationModel{
			ID:               primitive.NewObjectID(),
			SenderID:         objectID,
			ReceiverID:       follow.Follower,
			NotificationType: models.POST_,
			TargetID:         post.ID,
			IsRead:           false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		_, err := s.notificationModel.InsertOne(ctx, notification)
		if err != nil {
			return models.PostModel{}, err
		}

		r := utils.GetRedis()
		_ = r.Incr(ctx, "unread_count:"+follow.Follower.Hex()).Err()
		_ = r.LPush(ctx, "notifications:"+follow.Follower.Hex(), notification.ID.Hex()).Err()
	}

	return post, nil
}

func (s *PostService) UpdatePost(
	ctx context.Context,
	postID string,
	user *utils.TokenPayload,
	content string,
) (models.PostModel, error) {

	pID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return models.PostModel{}, errors.New("invalid post id")
	}

	var post models.PostModel

	err = s.postModel.FindOne(ctx, bson.M{"_id": pID}).Decode(&post)

	if err != nil {
		return models.PostModel{}, errors.New("post not found")
	}

	isOwner := post.UserID.Hex() == user.UserID
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		return models.PostModel{}, errors.New("not allowed")
	}

	if content != "" {
		post.Content = content
	}

	post.UpdatedAt = time.Now()

	_, err = s.postModel.UpdateOne(
		ctx,
		bson.M{"_id": pID},
		bson.M{
			"$set": bson.M{
				"content":   post.Content,
				"updatedAt": post.UpdatedAt,
			},
		},
	)

	if err != nil {
		return models.PostModel{}, err
	}

	return post, nil
}

func (s *PostService) DeletePost(
	ctx context.Context,
	postID string,
	user *utils.TokenPayload,
) (models.PostModel, error) {
	pID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return models.PostModel{}, errors.New("invalid post id")
	}

	var post models.PostModel

	err = s.postModel.FindOne(ctx, bson.M{"_id": pID}).Decode(&post)
	if err != nil {
		return models.PostModel{}, errors.New("post not found")
	}

	isOwner := post.UserID.Hex() == user.UserID
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		return models.PostModel{}, errors.New("not allowed")
	}

	for _, img := range post.Images {
		if img.PublicID != "" {
			_ = s.cloudinaryService.Delete(img.PublicID)
		}
	}

	_, err = s.postModel.DeleteOne(ctx, bson.M{"_id": pID})
	if err != nil {
		return models.PostModel{}, err
	}

	return models.PostModel{}, nil
}

func (s *PostService) GetAllPosts(
	ctx context.Context,
	page int64,
	limit int64,
) ([]models.PostModel, int64, error) {

	skip := (page - 1) * limit

	cursor, err := s.postModel.Find(
		ctx,
		bson.M{},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []models.PostModel

	for cursor.Next(ctx) {
		var post models.PostModel
		if err := cursor.Decode(&post); err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	total, err := s.postModel.CountDocuments(ctx, bson.M{})

	if err != nil {
		return nil, 0, err
	}

	// manual pagination
	start := skip
	end := skip + limit

	if start > int64(len(posts)) {
		return []models.PostModel{}, total, nil
	}

	if end > int64(len(posts)) {
		end = int64(len(posts))
	}

	return posts[start:end], total, nil
}

func (s *PostService) GetPostByID(
	ctx context.Context,
	postID string,
) (models.PostModel, error) {

	pID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return models.PostModel{}, errors.New("invalid post id")
	}

	var post models.PostModel

	err = s.postModel.FindOne(ctx, bson.M{"_id": pID}).Decode(&post)
	if err != nil {
		return models.PostModel{}, errors.New("post not found")
	}

	return post, nil
}

func (s *PostService) GetPostsByUser(
	ctx context.Context,
	userID string,
	page int64,
	limit int64,
) ([]models.PostModel, int64, error) {

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user id")
	}

	skip := (page - 1) * limit

	cursor, err := s.postModel.Find(
		ctx,
		bson.M{"userId": uID},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []models.PostModel

	for cursor.Next(ctx) {
		var post models.PostModel
		if err := cursor.Decode(&post); err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	total, err := s.postModel.CountDocuments(ctx, bson.M{"userId": uID})
	if err != nil {
		return nil, 0, err
	}

	start := skip
	end := skip + limit

	if start > int64(len(posts)) {
		return []models.PostModel{}, total, nil
	}

	if end > int64(len(posts)) {
		end = int64(len(posts))
	}

	return posts[start:end], total, nil
}
