package service

import (
	"context"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	respond "rooms_service/internal/models/dto/response"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (r *RoomService) RegisterUser(ctx context.Context, req request.RegisterRequest) (models.User, error){
	if len(req.Email) < 8 || len(req.Password) < 4 {
		return models.User{}, models.ErrInvalidRequest
	}

	if r.user.CheckEmailExist(ctx, req.Email) {
		return models.User{}, models.ErrNotFound
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	
	if err != nil {
		return models.User{}, models.ErrInternalError
	}
	id := uuid.New()
	t := time.Now()
	userDB := models.UserDB{ID: id, Email: req.Email, HashedPassword: string(hashedPassword), Role: req.Role, CreatedAt: t}
	
	errorCode := r.user.AddUserToDB(ctx, userDB)
	if errorCode != nil {
		return models.User{}, errorCode
	}

	return models.User{ID: id, Email: req.Email, Role: req.Role, CreatedAt: t}, nil
}

func (r *RoomService) LoginUser(ctx context.Context, req request.LoginRequest) (respond.Token, error) {
	userdb, errorCode := r.user.GetUserByEmail(ctx, req.Email)
	if errorCode != nil {
		return respond.Token{}, errorCode
	}

	err := bcrypt.CompareHashAndPassword([]byte(userdb.HashedPassword), []byte(req.Password))
	if err != nil {
		return respond.Token{}, models.ErrNotFound
	}
	tokenResponse, err1 := GenerateToken(userdb.ID, userdb.Role)
	if err1 != nil {
		return respond.Token{}, err1
	}
	return tokenResponse, nil
}