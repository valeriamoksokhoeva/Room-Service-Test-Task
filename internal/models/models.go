package models

import (
	"time"

	"github.com/google/uuid"
)


type StatusT string
type RoleT string 

const (
	Active StatusT = "active"
	Cancelled StatusT = "cancelled"
)

const (
	Admin RoleT = "admin"
	UserT RoleT = "user"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	Role RoleT `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Capacity int `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
}

type Schedule struct {
	ID uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"room_id"`
	DaysOfWeek []int `json:"days_of_week"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
}

type Slot struct {
	ID uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"room_id"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
}

type Booking struct {
	ID uuid.UUID `json:"id"`
	SlotID uuid.UUID `json:"slot_id"`
	UserID uuid.UUID `json:"user_id"`
	Status StatusT `json:"status"`
	ConferenceLink string `json:"link"`
	CreatedAt time.Time `json:"created_at"`
}

type Pagination struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Total int `json:"total"`
}