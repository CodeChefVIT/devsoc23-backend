package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeLine struct {
	Id          primitive.ObjectID `bson:"_id"`
	Title       *string            `json:"title" validate:"required"`
	SubTitle    *string            `json:"subTitle"`
	Description *string            `json:"description"`
	Venue       *string            `json:"venue"`
	Date        *string            `json:"date" validate:"required"`
	StartTime   *string            `json:"startTime" validate:"required"`
	EndTime     *string            `json:"endTime" validate:"required"`
}
