package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReelModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content string             `bson:"content" json:"content"`
	Reel    string             `bson:"reel" json:"reel" validate:"required"`
	UserID  primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
