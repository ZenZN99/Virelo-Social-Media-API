package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FollowModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Follower    primitive.ObjectID `bson:"follower" json:"follower"`
	Following   primitive.ObjectID `bson:"following" json:"following"`
	IsFollowing bool               `bson:"isFollowing" json:"isFollowing" default:"false"`
}
