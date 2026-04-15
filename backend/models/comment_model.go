package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContentType string

const (
	Post    ContentType = "post"
	Reel    ContentType = "reel"
	Comment ContentType = "comment"
)
type CommentModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`
	TargetID        primitive.ObjectID `bson:"targetId" json:"targetId" validate:"required"`
	TargetType      ContentType         `bson:"targetType" json:"targetType" validate:"required"`
	Text            string             `bson:"text" json:"text" validate:"required"`
	ParentCommentId *primitive.ObjectID            `bson:"parentCommentId,omitempty" json:"parentCommentId"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
