package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationType string

const (
	POST_    NotificationType = "post"
	REEL_    NotificationType = "reel"
	LIKE_    NotificationType = "like"
	COMMENT_ NotificationType = "comment"
	FOLLOW_  NotificationType = "follow"
)

type NotificationModel struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SenderID         primitive.ObjectID `bson:"senderId" json:"senderId" validate:"required"`
	ReceiverID       primitive.ObjectID `bson:"receiverId" json:"receiverId" validate:"required"`
	NotificationType NotificationType   `bson:"notificationType" json:"notificationType" validate:"required"`
	TargetID         primitive.ObjectID `bson:"targetId" json:"targetId" validate:"required"`
	IsRead           bool               `bson:"isRead" json:"isRead" default:"false"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
