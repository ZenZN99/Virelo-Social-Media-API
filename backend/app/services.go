package app

import (
	"backend/services"
	"backend/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	UserService         *services.UserService
	PostService         *services.PostService
	ReelService         *services.ReelService
	CommentService      *services.CommentService
	LikeService         *services.LikeService
	FollowService       *services.FollowService
	MessageService      *services.MessageService
	NotificationService *services.NotificationService
	TokenService        *utils.TokenService
}

func InitServices(
	userCollection *mongo.Collection,
	postCollection *mongo.Collection,
	reelCollection *mongo.Collection,
	commentCollection *mongo.Collection,
	likeCollection *mongo.Collection,
	followCollection *mongo.Collection,
	messageCollection *mongo.Collection,
	notificationCollection *mongo.Collection,
) *Services {

	tokenService := utils.NewTokenService()
	cloudinaryService := utils.NewCloudinaryService()

	userService := services.NewUserService(
		userCollection,
		tokenService,
		cloudinaryService,
	)

	postService := services.NewPostService(
		postCollection,
		followCollection,
		notificationCollection,
		cloudinaryService,
	)

	reelService := services.NewReelService(
		reelCollection,
		followCollection,
		notificationCollection,
		cloudinaryService,
	)

	commentService := services.NewCommentService(
		commentCollection,      
		notificationCollection, 
		postCollection,         
		reelCollection,         
	)

	likeService := services.NewLikeService(
		likeCollection,
		postCollection,
		reelCollection,
		commentCollection,
		notificationCollection,
	)

	followService := services.NewFollowService(
		followCollection,
		notificationCollection,
	)

	messageService := services.NewMessageService(
		messageCollection,
		cloudinaryService,
	)

	notificationService := services.NewNotificationService(
		notificationCollection,
	)

	return &Services{
		UserService:         userService,
		PostService:         postService,
		ReelService:         reelService,
		CommentService:      commentService,
		LikeService:         likeService,
		FollowService:       followService,
		MessageService:      messageService,
		NotificationService: notificationService,
		TokenService:        tokenService,
	}
}
