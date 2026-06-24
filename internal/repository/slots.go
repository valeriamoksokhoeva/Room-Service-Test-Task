package repository

import (
	"context"
	"log"
	"rooms_service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (rc *SlotRepo) InsertSlots(ctx context.Context, slots []models.Slot)  error {
    log.Printf("inserting %d slots", len(slots))  // сколько слотов пришло
    
    if len(slots) == 0 {
        log.Println("no slots to insert")
        return nil
    }

    batch := &pgx.Batch{}

    query := `
        INSERT INTO slots (id, room_id, start_at, end_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT DO NOTHING
    `

    for _, slot := range slots {
        log.Printf("queuing slot: %v %v %v", slot.RoomID, slot.Start, slot.End)
        batch.Queue(query, slot.ID, slot.RoomID, slot.Start, slot.End)
    }

    results := rc.conn.SendBatch(ctx, batch)
    defer func() {
		if err := results.Close(); err != nil {
			log.Printf("batch close error: %v", err)
		}
	}()

    for range slots {
        if _, err := results.Exec(); err != nil {
            log.Printf("batch exec error: %v", err)  // реальная ошибка
            return models.ErrInternalError
        }
    }
    return nil
}

func (rc *SlotRepo) GetAvailableSlots(ctx context.Context, room_id uuid.UUID, today time.Time) ([]models.Slot, error) {
	query := `
		SELECT id, room_id, start_at, end_at 
		FROM slots
		WHERE room_id = $1
		AND start_at >= $2   
		AND start_at < $3        
		AND NOT EXISTS (
			SELECT 1 FROM bookings
			WHERE bookings.slot_id = slots.id
			AND bookings.status = 'active'
		)
		ORDER BY start_at
	`
	var slots []models.Slot
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	log.Printf("searching slots for room: %v, from: %v, to: %v", room_id, startOfDay, endOfDay)
	rows, err := rc.conn.Query(ctx, query, room_id, startOfDay, endOfDay)
	if err != nil {
		log.Printf("get_available_slots error: %v", err)
		return nil, models.ErrInternalError
	}
	defer rows.Close()

	for rows.Next() {
		var slot models.Slot
		err := rows.Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End)
		if err != nil {
			log.Println("error rows scan")
			return nil, models.ErrInternalError
		}
		slots = append(slots, slot)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error")
		return nil, models.ErrInternalError
	}
	return slots, nil
}

func (rc *SlotRepo) CheckAvailableSlot(ctx context.Context, slot_id uuid.UUID) bool {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM slots WHERE id = $1 AND start_at >= $2
		)
	`
	var check bool
	err := rc.conn.QueryRow(ctx, query, slot_id, time.Now()).Scan(&check)
	if err != nil {
		return false
	}
	return check
}