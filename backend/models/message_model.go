package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SenderID   primitive.ObjectID `bson:"senderId" json:"senderId" validate:"required"`
	ReceiverID primitive.ObjectID `bson:"receiverId" json:"receiverId" validate:"required"`
	Content    string             `bson:"content" json:"content"`
	Image      string             `bson:"image" json:"image"`
	IsRead     bool               `bson:"isRead" json:"isRead" default:"false"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
