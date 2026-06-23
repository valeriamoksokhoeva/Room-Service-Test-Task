package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"rooms_service/internal/models"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("token.Method.(*jwt.SigningMethodHMAC); !ok")
            return nil, models.ErrInternalError
        }
        return []byte(os.Getenv("SECRET_KEY")), nil
    })

    if err != nil {
        return nil, models.ErrInternalError
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, models.ErrInternalError
    }

    return claims, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			HandleError(w, models.ErrUnauthorized)
			return 
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := ParseToken(tokenString)

		if err != nil {
			HandleError(w, models.ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next(w, r.WithContext(ctx))
	}
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxRole, ok := r.Context().Value("role").(string)
		if !ok {
			HandleError(w, models.ErrUnauthorized)
			return
		}
		if ctxRole != role {
			HandleError(w, models.ErrForbidden)
			return 
		}
		next(w, r)
	}
}