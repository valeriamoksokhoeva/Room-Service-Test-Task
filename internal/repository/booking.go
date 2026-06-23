package repository

import (
	"context"
	"log"
	"rooms_service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (rc *BookingRepo) InsertBooking(ctx context.Context, booking models.Booking) error {
	tx, err := rc.conn.Begin(ctx)
	if err != nil {
		return models.ErrInternalError
	}
	defer tx.Rollback(ctx)

	var slot_id uuid.UUID
	query1 := `
		SELECT id FROM slots
		WHERE id = $1 FOR UPDATE
	`
	err = tx.QueryRow(ctx, query1, booking.SlotID).Scan(&slot_id)
	if err != nil {
		return models.ErrSlotNotFound
	}

	var count int
	query2 := `
		SELECT COUNT(*) FROM bookings
		WHERE slot_id = $1 AND status = 'active'
	`

	err = tx.QueryRow(ctx, query2, booking.SlotID).Scan(&count)
	if err != nil {
		return models.ErrInternalError
	}
	if count > 0 {
		return models.ErrSlotAlreadyBooked
	}
	
	query3 := `
		INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, query3, booking.ID, booking.SlotID, booking.UserID, booking.Status, booking.ConferenceLink, booking.CreatedAt)
	if err != nil {
		return models.ErrInternalError
	}
	return tx.Commit(ctx)
}

func (rc *BookingRepo) MyBookings(ctx context.Context, user_id uuid.UUID) ([]models.Booking, error) {
	query := `
		SELECT bookings.id, bookings.slot_id, bookings.user_id, bookings.status, bookings.conference_link, bookings.created_at 
		FROM bookings
		JOIN slots ON bookings.slot_id = slots.id
		WHERE bookings.user_id = $1 AND slots.start_at >= $2
		ORDER BY slots.start_at
	`
	rows, err := rc.conn.Query(ctx, query, user_id, time.Now())
	if err != nil {
		return nil, models.ErrInternalError
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var boo models.Booking
		err = rows.Scan(&boo.ID, &boo.SlotID, &boo.UserID, &boo.Status, &boo.ConferenceLink, &boo.CreatedAt)
		if err != nil {
			return nil, models.ErrInternalError
		}
		bookings = append(bookings, boo)
	}

	if err = rows.Err(); err != nil{
		return nil, models.ErrInternalError
	}
	return bookings, nil
}

func (rc *BookingRepo) GetBookingByID(ctx context.Context, booking_id uuid.UUID) (models.Booking, error) {
	query := `
		SELECT id, slot_id, user_id, status, conference_link, created_at FROM bookings
		WHERE id = $1
	`
	
	var booking models.Booking
	err := rc.conn.QueryRow(ctx, query, booking_id).Scan(&booking.ID, &booking.SlotID, &booking.UserID,
        &booking.Status, &booking.ConferenceLink, &booking.CreatedAt)
	
	if err == pgx.ErrNoRows {
		return models.Booking{}, models.ErrBookingNotFound
	}
	if err != nil {
		log.Printf("get booking by id error: %v", err)
		return models.Booking{}, models.ErrInternalError
	}
	return booking, nil
}

func (rc *BookingRepo) UpdateBookingStatus(ctx context.Context, booking_id uuid.UUID) error {
	query := `
		UPDATE bookings
		SET status = 'cancelled'
		WHERE id = $1
	`
	_, err := rc.conn.Exec(ctx, query, booking_id)

	if err != nil {
		log.Printf("UpdateBookingStatus error: %v", err)
		return models.ErrInternalError
	}
	return nil
}

func (rc *BookingRepo) GetAllBookings(ctx context.Context, page int, pageSize int) ([]models.Booking, int, error) {
	var total int
	query1 := `
		SELECT COUNT(*) FROM bookings
	`
	err := rc.conn.QueryRow(ctx, query1).Scan(&total)
	if err != nil {
		log.Printf("GetAllBookings error: %v", err)
		return nil, 0, models.ErrInternalError
	}

	offset := (page - 1) * pageSize

	query2 := `
		SELECT id, slot_id, user_id, status, conference_link, created_at FROM bookings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := rc.conn.Query(ctx, query2, pageSize, offset)
	if err != nil {
		log.Printf("GetAllBookings error: %v", err)
		return nil, 0, models.ErrInternalError
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		err = rows.Scan(&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt)
		if err != nil {
			return nil, 0, models.ErrInternalError
		}
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}

	return bookings, total, nil
}