package repository

import (
	"context"

	"rooms_service/internal/models"
)

func (rc *UserRepo) CheckEmailExist(ctx context.Context, email string) bool {
	queryExists := `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE email = $1
		)
	`
	var check bool
	err1 := rc.conn.QueryRow(ctx, queryExists, email).Scan(&check)
	if err1 != nil {
		return false
	}
	return check
}
func (rc *UserRepo) AddUserToDB(ctx context.Context, user models.UserDB) error {
	query := `
		INSERT INTO users (id, email, password, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := rc.conn.Exec(ctx, query, user.ID, user.Email, user.HashedPassword, user.Role, user.CreatedAt)
	if err != nil {
		return models.ErrInternalError
	}
	return nil
}

func (rc *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.UserDB, error) {
    query := `
        SELECT id, email, password, role, created_at FROM users
        WHERE email = $1
    `
    var userdb models.UserDB
    err := rc.conn.QueryRow(ctx, query, email).Scan(
        &userdb.ID,
        &userdb.Email,
        &userdb.HashedPassword,
        &userdb.Role,
        &userdb.CreatedAt,
    )
    if err != nil {
        return models.UserDB{}, models.ErrInternalError
    }
    return userdb, nil
}