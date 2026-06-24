package service

import (
	"context"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	"time"

	"github.com/google/uuid"
)

func (r *RoomService) GetRoomsList(ctx context.Context) ([]models.Room, error) {
	rooms, err := r.room.GetRooms(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomService) AddRoom(ctx context.Context, req request.RoomRequest) (models.Room, error) {
	room := models.Room{ID: uuid.New(), Name: req.Name, Description: req.Description, Capacity: req.Capacity, CreatedAt: time.Now()}

	err := r.room.CreateRoom(ctx, room)
	if err != nil {
		return models.Room{}, err
	}
	return room, nil
}