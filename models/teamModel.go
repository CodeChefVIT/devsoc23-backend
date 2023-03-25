package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	Id               primitive.ObjectID   `bson:"_id"`
	TeamName         *string              `json:"teamName" validate:"required"`
	TeamLeaderId     primitive.ObjectID   `bson:"leaderId"`
	TeamMembers      []primitive.ObjectID `bson:"teamMember" json:"teamMember"`
	TeamSize         int                  `bson:"teamSize" json:"teamSize"`
	ProjectId        primitive.ObjectID   `bson:"projectId,omitempty"`
	InvitedTeammates []primitive.ObjectID `json:"invitedTeammates,omitempty"`
	Round            string               `json:"round,omitempty"`
	IsFinalised      bool                 `json:"isFinalised,omitempty"`
	InviteCode       string               `bson:"inviteCode" json:"inviteCode,omitempty"`
	CreatedAt        time.Time            `json:"createdTime"`
	UpdatedAt        time.Time            `json:"updatedTime"`
}

type CreateTeamRequest struct {
	TeamName *string `bson:"teamName" json:"teamName,omitempty"`
}

type UpdateTeam struct {
	Id               primitive.ObjectID   `bson:"_id,omitempty"`
	TeamName         *string              `json:"teamName,omitempty"`
	TeamLeaderId     primitive.ObjectID   `bson:"leaderId,omitempty"`
	TeamMembers      []primitive.ObjectID `json:"teamMember,omitempty"`
	TeamSize         int                  `json:"teamSize,omitempty"`
	ProjectId        primitive.ObjectID   `bson:"projectId,omitempty"`
	InvitedTeammates *string              `json:"invitedTeammates,omitempty"`
	Round            *string              `json:"round,omitempty"`
	IsFinalised      bool                 `json:"isFinalised,omitempty"`
	InviteLink       *string              `json:"InviteLink,omitempty"`
	CreatedAt        time.Time            `json:"createdTime,omitempty"`
	UpdatedAt        time.Time            `json:"updatedTime"`
}
