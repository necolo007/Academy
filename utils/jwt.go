package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
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
func GenerateToken(userID int32) (string, error) {
	token_id, err := GenerateRandomToken()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 设置2小时有效期
		"token_id": token_id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyToken 验证JWT
func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保使用的是相同的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	// 验证 token 是否有效
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}
	return false, nil
}
