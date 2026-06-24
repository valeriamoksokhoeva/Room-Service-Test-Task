package service

import (
	"context"
	"log"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	"time"

	"github.com/google/uuid"
)

func (r *RoomService) CreateBooking(ctx context.Context, req request.BookingRequest, user_id uuid.UUID) (models.Booking, error) {
	if !r.slot.CheckAvailableSlot(ctx, req.SlotId) {
		return models.Booking{}, models.ErrSlotNotFound
	}

	booking := models.Booking{
		ID: uuid.New(),
		SlotID: req.SlotId,
		UserID: user_id,
		Status: models.Active,
		ConferenceLink: "",
		CreatedAt: time.Now(),
	}

	if req.CreateConferenceLink {
		link, err := r.conference.CreateLink(booking.ID)
		if err != nil {
			log.Printf("conference service error: %v", err)
		} else {
			booking.ConferenceLink = link
		}
	}
	err := r.booking.InsertBooking(ctx, booking)
	if err != nil {
		return models.Booking{}, err
	}
	return booking, nil
}

func (r *RoomService) GetMyBookings(ctx context.Context, user_id uuid.UUID) ([]models.Booking, error) {
	bookings, err := r.booking.MyBookings(ctx, user_id)
	if err != nil {
		log.Println("error with db my bookings")
		return nil, err
	}

	return bookings, nil
}

func (r *RoomService) CancelMyBooking(ctx context.Context, booking_id uuid.UUID, user_id uuid.UUID) (models.Booking, error) {
	booking, err := r.booking.GetBookingByID(ctx, booking_id)
	if err != nil {
		return models.Booking{}, err
	}

	if booking.UserID != user_id {
		return models.Booking{}, models.ErrForbidden
	}

	if booking.Status == models.Cancelled {
		return booking, nil
	}

	err = r.booking.UpdateBookingStatus(ctx, booking_id)
	if err != nil {
		return models.Booking{}, err
	}
	booking.Status = models.Cancelled
	return booking, nil
}

func (r *RoomService) AllBookings(ctx context.Context, page int, pageSize int) ([]models.Booking, int,  error) {
	bookings, total, err := r.booking.GetAllBookings(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}