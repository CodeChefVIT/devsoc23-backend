package models

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `form:"firstName" validate:"required,min=2,max=16"`
	LastName     *string            `form:"lastName" validate:"min=2,max=32"`
	Email        *string            `form:"email" validate:"required,email"`
	Password     *string            `form:"password" validate:"required,min=8,max=64"`
	PhoneNumber  *string            `form:"phoneNumber" validate:"required,e164,min=10,max=13"`
	Token        *string            `form:"token,omitempty"`
	Bio          *string            `form:"bio,omitempty"`
	Gender       *string            `form:"gender,omitempty"`
	UserRole     string             `form:"userRole" default:"HACKER" validate:"eq=ADMIN|eq=HACKER|eq=VOTER"`
	ProfilePhoto *string            `form:"profilePhoto,omitempty"`
	RegNo        *string            `form:"redNo,omitempty"`
	College      *string            `form:"college,omitempty" validate:"min=10,max=128"`
	CollegeYear  *string            `form:"collegeYear,omitempty"`
	BirthDate    *string            `form:"birthData,omitempty"`
	VerifyOtp    *string            `form:"verifyOtp,omitempty"`
	Address      *string            `form:"address,omitempty"`
	QrData       *string            `form:"qrData,omitempty"`
	Image        *string            `form:"image,omitempty"`
	IsActive     bool               `form:"isActive,omitempty"`
	IsVerify     bool               `form:"isVerify,omitempty"`
	IsCanShare   bool               `form:"isCanShare,omitempty"`
	IsCheckedIn  bool               `form:"isCheckedIn,omitempty"`
	InTeam       bool               `form:"inTeam,omitempty"`
	IsBoard      bool               `form:"isBoard,omitempty" default:"false"`
	CreatedAt    time.Time          `form:"createdTime"`
	UpdatedAt    time.Time          `form:"updatedTime"`
	TeamId       *string            `form:"teamId,omitempty"`
	UserId       string             `form:"userId"`
}

type UpdateUserRequest struct {
	FirstName     *string               `form:"firstName" validate:"required,min=2,max=16"`
	LastName      *string               `form:"lastName" validate:"min=2,max=32"`
	Email         *string               `form:"email" validate:"required,email"`
	PhoneNumber   *string               `form:"phoneNumber" validate:"required,e164,min=10,max=13"`
	ProfilePhoto  *string               `form:"profilePhoto,omitempty"`
	College       *string               `form:"college,omitempty" validate:"min=10,max=128"`
	CollegeYear   *string               `form:"collegeYear,omitempty"`
	Bio           *string               `form:"bio,omitempty"`
	RegNo         *string               `form:"redNo,omitempty"`
	Gender        *string               `form:"gender,omitempty"`
	BirthDate     *string               `form:"birthData,omitempty"`
	VerifyOtp     *string               `form:"verifyOtp,omitempty"`
	Address       *string               `form:"address,omitempty"`
	QrData        *string               `form:"qrData,omitempty"`
	IsActive      bool                  `form:"isActive,omitempty"`
	IsVerify      bool                  `form:"isVerify,omitempty"`
	IsCanShare    bool                  `form:"isCanShare,omitempty"`
	IsCheckedIn   bool                  `form:"isCheckedIn,omitempty"`
	InTeam        bool                  `form:"inTeam,omitempty"`
	UpdatedAt     time.Time             `form:"updatedTime"`
	TeamId        *string               `form:"teamId,omitempty"`
	UserId        string                `form:"userId"`
	CampaignImage *multipart.FileHeader `form:"image" binding:"required"`
}

type CreateUserRequest struct {
	Id            primitive.ObjectID    `bson:"_id"`
	FirstName     *string               `form:"firstName"`
	LastName      *string               `form:"lastName"`
	Email         *string               `form:"email"`
	Password      *string               `form:"password"`
	PhoneNumber   *string               `form:"phoneNumber"`
	RegNo         *string               `form:"redNo,omitempty"`
	College       *string               `form:"college,omitempty"`
	CollegeYear   *string               `form:"collegeYear,omitempty"`
	Bio           *string               `form:"bio,omitempty"`
	Gender        *string               `form:"gender,omitempty"`
	BirthDate     *string               `form:"birthDate,omitempty"`
	IsActive      bool                  `form:"isActive,omitempty"`
	IsVerify      bool                  `form:"isVerify,omitempty"`
	IsCanShare    bool                  `form:"isCanShare,omitempty"`
	IsCheckedIn   bool                  `form:"isCheckedIn,omitempty"`
	InTeam        bool                  `form:"inTeam,omitempty"`
	IsBoard       bool                  `form:"isBoard,omitempty" default:"USER"`
	CreatedAt     time.Time             `form:"createdTime"`
	UpdatedAt     time.Time             `form:"updatedTime"`
	UserId        string                `form:"userId"`
	CampaignImage *multipart.FileHeader `form:"image"`
}

type LoginUserRequest struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=8,max=64"`
}

type ResetPasswordRequest struct {
	Oldpass *string `json:"oldpass" validate:"required,min=8,max=64"`
	Newpass *string `json:"newpass" validate:"required,min=8,max=64"`
}

type ForgetPasswordRequest struct {
	Email   string `json:"email"`
	OTP     string `json:"otp"`
	Newpass string `json:"newpass"`
}
