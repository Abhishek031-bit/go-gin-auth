package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null" json:"email" binding:"email,required"`
	Password string `gorm:"not null" json:"-" binding:"required"`
	Files    []File
}

type File struct {
	gorm.Model
	FilePath     string `gorm:"not null"`
	OriginalName string `gorm:"not null"`
	UserID       uint
}

type AuthRequest struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=8"`
}
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
