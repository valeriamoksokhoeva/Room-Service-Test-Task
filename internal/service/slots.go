package service

import (
	"context"
	"rooms_service/internal/models"
	"time"

	"github.com/google/uuid"
)

func (r *RoomService) ListSlots(ctx context.Context, room_id uuid.UUID, today time.Time) ([]models.Slot, error){
	slots, err := r.slot.GetAvailableSlots(ctx, room_id, today)
	if err != nil {
		return nil, err
	}
	return slots, nil
}