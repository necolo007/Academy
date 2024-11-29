package utils

import (
	"Academy/Model"
	"crypto/rand"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("HduHelperMember")

// GenerateRandomToken 生成一个随机的token_id
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 16) // 生成128位（16字节）随机数
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateToken 生成JWT
func GenerateToken(user Model.User) (string, error) {
	tokenID, err := GenerateRandomToken()
	if err != nil {
		return "", err
	}
	claims := &Model.UserClaims{
		UserId: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Id:        tokenID,
			Issuer:    "cxr",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateToken VerifyToken 验证JWT
func ValidateToken(tokenString string) (*Model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Model.UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
