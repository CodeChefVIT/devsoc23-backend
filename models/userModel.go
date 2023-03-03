package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=16"`
	LastName     *string            `json:"lastName" validate:"min=2,max=32"`
	Email        *string            `json:"email" validate:"required,email"`
	Password     *string            `json:"password" validate:"required,min=8,max=64"`
	PhoneNumber  *string            `json:"phoneNumber" validate:"required,e164,min=10,max=13"`
	Token        *string            `json:"token,omitempty"`
	UserRole     string             `json:"userRole" default:"HACKER" validate:"eq=ADMIN|eq=HACKER|eq=VOTER"`
	ProfilePhoto *string            `json:"profilePhoto,omitempty"`
	College      *string            `json:"college,omitempty" validate:"min=10,max=128"`
	CollegeYear  *string            `json:"collegeYear,omitempty"`
	BirthDate    *string            `json:"birthData,omitempty"`
	VerifyOtp    *string            `json:"verifyOtp,omitempty"`
	Address      *string            `json:"address,omitempty"`
	QrData       *string            `json:"qrData,omitempty"`
	IsActive     bool               `json:"isActive,omitempty"`
	IsVerify     bool               `json:"isVerify,omitempty"`
	IsCanShare   bool               `json:"isCanShare,omitempty"`
	IsCheckedIn  bool               `json:"isCheckedIn,omitempty"`
	InTeam       bool               `json:"inTeam,omitempty"`
	IsBoard      bool               `json:"isBoard,omitempty" default:"false"`
	CreatedAt    time.Time          `json:"createdTime"`
	UpdatedAt    time.Time          `json:"updatedTime"`
	TeamId       *string            `json:"teamId,omitempty"`
	UserId       string             `json:"userId"`
}

type CreateUserRequest struct {
	Id          primitive.ObjectID `bson:"_id"`
	FirstName   *string            `json:"firstName" validate:"required,min=2,max=16"`
	LastName    *string            `json:"lastName" validate:"min=2,max=32"`
	Email       *string            `json:"email" validate:"required,email"`
	Password    *string            `json:"password" validate:"required,min=8,max=64"`
	PhoneNumber *string            `json:"phoneNumber" validate:"required,e164,min=10,max=13"`
	College     *string            `json:"college,omitempty" validate:"min=10,max=128"`
	CollegeYear *string            `json:"collegeYear,omitempty"`
	BirthDate   *string            `json:"birthDate,omitempty"`
	IsActive    bool               `json:"isActive,omitempty"`
	IsVerify    bool               `json:"isVerify,omitempty"`
	IsCanShare  bool               `json:"isCanShare,omitempty"`
	IsCheckedIn bool               `json:"isCheckedIn,omitempty"`
	InTeam      bool               `json:"inTeam,omitempty"`
	IsBoard     bool               `json:"isBoard,omitempty" default:"USER"`
	CreatedAt   time.Time          `json:"createdTime"`
	UpdatedAt   time.Time          `json:"updatedTime"`
	UserId      string             `json:"userId"`
}

type LoginUserRequest struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=8,max=64"`
}
