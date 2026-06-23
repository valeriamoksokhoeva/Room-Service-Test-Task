package repository

import (
	"context"
	"time"

	"rooms_service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepo struct {
	conn *pgxpool.Pool
}

type UserRepo struct {
	conn *pgxpool.Pool
}

type SlotRepo struct {
	conn *pgxpool.Pool
}

type ScheduleRepo struct {
	conn *pgxpool.Pool
}

type BookingRepo struct {
	conn *pgxpool.Pool
}

func NewRoomRepository(conn *pgxpool.Pool) RoomRepository {
	return &RoomRepo{conn: conn}
}

func NewUserRepository(conn *pgxpool.Pool) UserRepository {
	return &UserRepo{conn: conn}
}

func NewSlotRepository(conn *pgxpool.Pool) SlotRepository {
	return &SlotRepo{conn: conn}
}

func NewScheduleRepository(conn *pgxpool.Pool) ScheduleRepository {
	return &ScheduleRepo{conn: conn}
}

func NewBookingRepository(conn *pgxpool.Pool) BookingRepository {
	return &BookingRepo{conn: conn}
}


type UserRepository interface {
    AddUserToDB(ctx context.Context, user models.UserDB) error
    GetUserByEmail(ctx context.Context, email string) (models.UserDB, error)
    CheckEmailExist(ctx context.Context, email string) bool
}

type RoomRepository interface {
    GetRooms(ctx context.Context) ([]models.Room, error)
    CreateRoom(ctx context.Context, room models.Room) error
    CheckRoomExist(ctx context.Context, roomID uuid.UUID) bool
}

type ScheduleRepository interface {
    CreateSchedule(ctx context.Context, schedule models.Schedule) error
    ScheduleExists(ctx context.Context, roomID uuid.UUID) bool
    GetAllSchedules(ctx context.Context) ([]models.Schedule, error)
}

type SlotRepository interface {
    InsertSlots(ctx context.Context, slots []models.Slot) error
    GetAvailableSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]models.Slot, error)
    CheckAvailableSlot(ctx context.Context, slotID uuid.UUID) bool
}

type BookingRepository interface {
	GetAllBookings(ctx context.Context, page int, pageSize int) ([]models.Booking, int, error)
    InsertBooking(ctx context.Context, booking models.Booking) error
    MyBookings(ctx context.Context, userID uuid.UUID) ([]models.Booking, error)
    GetBookingByID(ctx context.Context, bookingID uuid.UUID) (models.Booking, error)
    UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID) error
}