package model

import (
	//"project/internal/database"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `json:"name" validate:"required,unique" gorm:"unique,notnull"`
	Email        string `json:"email" validate:"required"`
	PasswordHash string `json:"-" validate:"required"`
	Dob          string `json:"Dob" validate:"required"`
}

type UserSignup struct {
	UserName string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Dob      string `json:"Dob" validate:"required"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgotPass struct {
	Email string `json:"email" validate:"required"`
	Dob   string `json:"dob" validate:"required"`
}

type Newpassword struct {
	Otp             string `json:"otp" validate:"required"`
	Newpassword     string `json:"newpassword" validate:"required"`
	ConfirmPassword string `json:"confirmpassword" validate:"required"`
	Email           string `json:"email" validate:"required"`
}
