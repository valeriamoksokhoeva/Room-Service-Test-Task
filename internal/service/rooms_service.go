package service

import (
	"context"
	"rooms_service/internal/conference"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	respond "rooms_service/internal/models/dto/response"
	"rooms_service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type RoomService struct {
	room repository.RoomRepository
	user repository.UserRepository
	schedule repository.ScheduleRepository
	slot repository.SlotRepository
	booking repository.BookingRepository
	conference conference.ConferenceService
}

func NewService(room repository.RoomRepository, user repository.UserRepository, schedule repository.ScheduleRepository, slot repository.SlotRepository, booking repository.BookingRepository, c conference.ConferenceService) DomainService {
	return &RoomService{room: room, user: user, schedule: schedule, slot: slot, booking: booking, conference: c}
}

type DomainService interface {
	RegisterUser(ctx context.Context, req request.RegisterRequest) (models.User, error)
	LoginUser(ctx context.Context, req request.LoginRequest) (respond.Token, error)
	GetRoomsList(ctx context.Context) ([]models.Room, error)
	AddRoom(ctx context.Context, req request.RoomRequest) (models.Room, error) 
	AddSchedule(ctx context.Context, req request.ScheduleRequest, room_id uuid.UUID) (models.Schedule, error)
	ListSlots(ctx context.Context, room_id uuid.UUID, today time.Time) ([]models.Slot, error)
	CreateBooking(ctx context.Context, req request.BookingRequest, user_id uuid.UUID) (models.Booking, error)
	GetMyBookings(ctx context.Context, user_id uuid.UUID) ([]models.Booking, error)
	CancelMyBooking(ctx context.Context, booking_id uuid.UUID, user_id uuid.UUID) (models.Booking, error)
	AllBookings(ctx context.Context, page int, pageSize int) ([]models.Booking, int,  error)
}