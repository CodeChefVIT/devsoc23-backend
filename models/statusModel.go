package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId primitive.ObjectID `bson:"userId,omitempty"`
	TeamId *string            `bson:"teamId,omitempty"`
	InHall bool               `bson:"inHall,omitempty"`
	Time   struct {
		Num      int    `bson:"num,omitempty"`
		IsTime   string `bson:"timeof,omitempty"`
		IfStatus string `bson:"status,omitempty"`
	} `bson:"time,omitempty"`
	UpdatedAt time.Time `json:"updatedTime"`
	CreatedAt time.Time `json:"CreatedTime"`
}
