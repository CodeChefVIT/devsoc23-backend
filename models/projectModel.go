package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id                 primitive.ObjectID `bson:"_id"`
	TeamId             *string            `bson:"teamId"`
	ProjectName        *string            `json:"projectName"`
	ProjectDescription *string            `json:"projectDescription"`
	ProjectStatus      *string            `json:"projectStatus,omitempty"`
	ProjectVideoLink   *string            `json:"projectVideoLink"`
	ProjectGithubLink  *string            `json:"projectGithubLink"`
	ProjectTrack       *string            `json:"projectTrack"`
	ProjectTags        []string           `json:"projectTags"`
	IsFinal            bool               `json:"isFinal,omitempty"`
}

type CreateProjectRequest struct {
	ProjectName        *string  `json:"projectName"`
	ProjectDescription *string  `json:"projectDescription"`
	ProjectVideoLink   *string  `json:"projectVideoLink"`
	ProjectGithubLink  *string  `json:"projectGithubLink"`
	ProjectTrack       *string  `json:"projectTrack"`
	ProjectTags        []string `json:"projectTags"`
}

type UpdateProjectRequest struct {
	Id                 primitive.ObjectID `bson:"_id"`
	TeamId             *string            `bson:"teamId"`
	ProjectName        *string            `json:"projectName"`
	ProjectDescription *string            `json:"projectDescription"`
	ProjectStatus      *string            `json:"projectStatus,omitempty"`
	ProjectVideoLink   *string            `json:"projectVideoLink"`
	ProjectGithubLink  *string            `json:"projectGithubLink"`
	ProjectTrack       *string            `json:"projectTrack"`
	ProjectTags        []string           `json:"projectTags"`
	IsFinal            bool               `json:"isFinal,omitempty"`
}