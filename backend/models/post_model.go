package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Content string `bson:"content" json:"content" validate:"required"`

	Images []PostImage `bson:"images" json:"images"`

	UserID primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type PostImage struct {
	URL      string `bson:"url" json:"url"`
	PublicID string `bson:"publicId" json:"publicId"`
}
