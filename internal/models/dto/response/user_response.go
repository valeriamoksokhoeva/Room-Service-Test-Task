package respond

import "rooms_service/internal/models"

type Token struct {
	AccessToken string `json:"token"`
}

type BookingResponse struct {
	Booking models.Booking `json:"booking"`
}

type BookingsListResponse struct {
	Bookings []models.Booking `json:"bookings"`
}

type UserResponse struct {
	User models.User `json:"user"`
}

type UserListResponse struct {
	Users []models.User `json:"users"`
}

type RoomResponse struct {
	Room models.Room `json:"room"`
}

type RoomsListResponse struct {
	Rooms []models.Room `json:"room"`
}

type SlotResponse struct {
	Slot models.Slot `json:"slot"`
}

type SlotListResponse struct {
	Slots []models.Slot `json:"slots"`
}

type ScheduleResponse struct {
	Schedule models.Schedule `json:"schedule"`
}