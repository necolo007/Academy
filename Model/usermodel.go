package Model

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Role     string `gorm:"default:'user'"`
}

type UserClaims struct {
	UserId uint
	Role   string
	jwt.StandardClaims
}
