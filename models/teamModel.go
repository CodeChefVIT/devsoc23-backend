package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	Id               primitive.ObjectID `bson:"_id"`
	TeamName         *string            `json:"teamName" validate:"required"`
	TeamLeader       *string            `json:"teamLeader"`
	TeamMembers      *string            `json:"teamMember"`
	TeamSize         *string            `json:"teamSize"`
	ProjectId        *string            `json:"projectId,omitempty"`
	InvitedTeammates *string            `json:"invitedTeammates,omitempty"`
	Round            *string            `json:"round,omitempty"`
	IsFinalised      bool               `json:"isFinalised,omitempty"`
	InviteLink       *string            `json:"InviteLink,omitempty"`
	CreatedAt        time.Time          `json:"createdTime"`
	UpdatedAt        time.Time          `json:"updatedTime"`
	TeamId           *string            `json:"teamId"`
}
