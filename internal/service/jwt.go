package service

import (
	"log"
	"os"
	"rooms_service/internal/models"
	respond "rooms_service/internal/models/dto/response"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func GenerateToken(userId uuid.UUID, role models.RoleT) (respond.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": userId, "role": role, "exp": time.Now().Add(24 * time.Hour).Unix()})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Printf("%v", err)
		return respond.Token{}, models.ErrInternalError
	}

	tokenResponse := respond.Token{AccessToken: tokenString}
	return tokenResponse, nil
}