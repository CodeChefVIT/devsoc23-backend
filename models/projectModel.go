package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id                 primitive.ObjectID   `bson:"_id"`
	TeamId             *string              `bson:"teamId"`
	ProjectName        *string              `json:"projectName"`
	ProjectTagLine     *string              `json:"projectTagLine"`
	ProjectStack       *string              `json:"projectStack"`
	ProjectDescription *string              `json:"projectDescription"`
	ProjectStatus      *string              `json:"projectStatus,omitempty"`
	ProjectVideoLink   *string              `json:"projectVideoLink"`
	ProjectGithubLink  *string              `json:"projectGithubLink"`
	ProjectTrack       *string              `json:"projectTrack"`
	ProjectTags        []string             `json:"projectTags"`
	IsFinal            bool                 `json:"isFinal,omitempty"`
	LikeCount          int                  `json:"like"`
	LikesId            []primitive.ObjectID `json:"likesId"`
}

type CreateProjectRequest struct {
	ProjectName        *string  `json:"projectName"`
	ProjectDescription *string  `json:"projectDescription"`
	ProjectVideoLink   *string  `json:"projectVideoLink"`
	ProjectTagLine     *string  `json:"projectTagLine"`
	ProjectStack       *string  `json:"projectStack"`
	ProjectGithubLink  *string  `json:"projectGithubLink"`
	ProjectTrack       *string  `json:"projectTrack"`
	ProjectTags        []string `json:"projectTags"`
}

type UpdateProjectRequest struct {
	Id                 primitive.ObjectID `bson:"_id"`
	TeamId             *string            `bson:"teamId"`
	ProjectName        *string            `json:"projectName"`
	ProjectTagLine     *string            `json:"projectTagLine"`
	ProjectStack       *string            `json:"projectStack"`
	ProjectDescription *string            `json:"projectDescription"`
	ProjectStatus      *string            `json:"projectStatus,omitempty"`
	ProjectVideoLink   *string            `json:"projectVideoLink"`
	ProjectGithubLink  *string            `json:"projectGithubLink"`
	ProjectTrack       *string            `json:"projectTrack"`
	ProjectTags        []string           `json:"projectTags"`
	IsFinal            bool               `json:"isFinal,omitempty"`
}
