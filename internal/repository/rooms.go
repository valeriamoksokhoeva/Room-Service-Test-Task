package repository

import (
	"context"
	"log"
	"rooms_service/internal/models"

	"github.com/google/uuid"
)

func (rc *RoomRepo) GetRooms(ctx context.Context) ([]models.Room, error) {
	query := `
		SELECT id, name, description, capacity, created_at FROM rooms
	`
	var rooms []models.Room
	rows, err := rc.conn.Query(ctx, query)
	if err != nil {
		log.Println("error with Query")
		return nil, models.ErrInternalError
	}
	defer rows.Close()

	
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)
		if err != nil {
			log.Println("error with rows.Scan()")
			return nil, models.ErrInternalError
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil{
		log.Println("error with rows.err()")
		return nil, models.ErrInternalError
	}
	return rooms, nil
}

func (rc *RoomRepo) CreateRoom(ctx context.Context, room models.Room) error {
	query := `
		INSERT INTO rooms (id, name, description, capacity, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := rc.conn.Exec(ctx, query, room.ID, room.Name, room.Description, room.Capacity, room.CreatedAt)
	if err != nil {
		return models.ErrInternalError
	}
	return nil
}

func (rc *RoomRepo) CheckRoomExist(ctx context.Context, room_id uuid.UUID) bool {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM rooms
			WHERE id = $1
		)
	`

	var check bool
	err1 := rc.conn.QueryRow(ctx, query, room_id).Scan(&check)
	if err1 != nil {
		return false
	}
	return check
}