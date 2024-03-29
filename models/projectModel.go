package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id                 primitive.ObjectID   `bson:"_id"`
	TeamId             *string              `bson:"teamId"`
	Idea               *string              `json:"idea"`
	IdeaLink           *string              `json:"ideaLink"`
	ProjectName        *string              `json:"projectName"`
	ProjectTagLine     *string              `json:"projectTagLine"`
	ProjectStack       *string              `json:"projectStack"`
	ProjectDescription *string              `json:"projectDescription"`
	ProjectStatus      *string              `json:"projectStatus,omitempty"`
	ProjectDriveLink   *string              `json:"projectDriveLink"`
	ProjectVideoLink   *string              `json:"projectVideoLink"`
	ProjectFigmaLink   *string              `json:"projectFigmaLink"`
	ProjectGithubLink  *string              `json:"projectGithubLink"`
	ProjectTrack       *string              `json:"projectTrack"`
	ProjectTags        []string             `json:"projectTags"`
	IsFinal            bool                 `json:"isFinal,omitempty"`
	LikeCount          int                  `json:"like"`
	LikesId            []primitive.ObjectID `json:"likesId"`
}
type CreateProjectIdeaRequest struct {
	ProjectName        *string `json:"projectName"`
	ProjectDescription *string `json:"projectDescription"`
	ProjectDriveLink   *string `json:"projectDriveLink"`
	ProjectFigmaLink   *string `json:"projectFigmaLink"`
	ProjectTrack       *string `json:"projectTrack"`
}

type CreateProjectRequest struct {
	ProjectName        *string  `json:"projectName"`
	ProjectDescription *string  `json:"projectDescription"`
	ProjectVideoLink   *string  `json:"projectVideoLink"`
	ProjectTagLine     *string  `json:"projectTagLine"`
	ProjectStack       *string  `json:"projectStack"`
	ProjectDriveLink   *string  `json:"projectDriveLink"`
	ProjectFigmaLink   *string  `json:"projectFigmaLink"`
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
	ProjectDriveLink   *string            `json:"projectDriveLink"`
	ProjectFigmaLink   *string            `json:"projectFigmaLink"`
	ProjectVideoLink   *string            `json:"projectVideoLink"`
	ProjectGithubLink  *string            `json:"projectGithubLink"`
	ProjectTrack       *string            `json:"projectTrack"`
	ProjectTags        []string           `json:"projectTags"`
	IsFinal            bool               `json:"isFinal,omitempty"`
}

type CreateProjectComment struct {
	Id        primitive.ObjectID `bson:"_id"`
	ProjectId *string            `json:"projectId" validate:"required"`
	Comment   *string            `json:"comment" validate:"required"`
}

type ProjectComment struct {
	Id        primitive.ObjectID `bson:"_id"`
	ProjectId *string            `json:"projectId" validate:"required"`
	Comment   *string            `bson:"comment"`
}
