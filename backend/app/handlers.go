package app

import "backend/handlers"

type Handlers struct {
	UserHandler         *handlers.UserHandler
	PostHandler         *handlers.PostHandler
	ReelHandler         *handlers.ReelHandler
	CommentHandler      *handlers.CommentHandler
	LikeHandler         *handlers.LikeHandler
	FollowHandler       *handlers.FollowHandler
	MessageHandler      *handlers.MessageHandler
	NotificationHandler *handlers.NotificationHandler
}

func InitHandlers(services *Services) *Handlers {
	return &Handlers{
		UserHandler:         handlers.NewUserHandler(services.UserService),
		PostHandler:         handlers.NewPostHandler(services.PostService),
		ReelHandler:         handlers.NewReelHandler(services.ReelService),
		CommentHandler:      handlers.NewCommentHandler(services.CommentService),
		LikeHandler:         handlers.NewLikeHandler(services.LikeService),
		FollowHandler:       handlers.NewFollowHandler(services.FollowService),
		MessageHandler:      handlers.NewMessageHandler(services.MessageService),
		NotificationHandler: handlers.NewNotificationHandler(services.NotificationService),
	}
}
