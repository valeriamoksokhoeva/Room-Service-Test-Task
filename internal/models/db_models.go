package models

import (
	"time"

	"github.com/google/uuid"
)

type UserDB struct {
	ID uuid.UUID `db:"id"`
	Email string `db:"email"`
	HashedPassword string `db:"password"`
	Role RoleT `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}