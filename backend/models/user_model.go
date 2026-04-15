package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	Admin UserRole = "admin"
	User  UserRole = "user"
)

type UserModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName  string             `bson:"fullname" json:"fullname" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=8"`
	Avatar    Avatar             `bson:"avatar" json:"avatar"`
	Role      string             `bson:"role" json:"role" validate:"required"`
	Bio       string             `bson:"bio" json:"bio"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type Avatar struct {
	URL      string `bson:"url" json:"url"`
	PublicID string `bson:"publicid" json:"publicid"`
}
