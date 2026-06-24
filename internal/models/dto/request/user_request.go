package request

import (
	"rooms_service/internal/models"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Role models.RoleT `json:"role"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type RoomRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Capacity int `json:"capacity"`
}

type BookingRequest struct {
	SlotId uuid.UUID `json:"slot_id"`
	CreateConferenceLink bool `json:"createConferenceLink"`
}
type DummyRegisterRequest struct {
	Role models.RoleT `json:"role"`
}

type ScheduleRequest struct {
	DaysOfWeek []int `json:"days_of_week"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
}