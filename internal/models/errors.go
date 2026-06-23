package models

type ErrorCode string

const (
	INVALID_REQUEST ErrorCode = "INVALID_REQUEST"
	UNAUTHORIZED ErrorCode = "UNAUTHORIZED"
	NOT_FOUND ErrorCode = "NOT_FOUND"
	ROOM_NOT_FOUND ErrorCode = "ROOM_NOT_FOUND"
	SLOT_NOT_FOUND ErrorCode = "SLOT_NOT_FOUND"
	SLOT_ALREADY_BOOKED ErrorCode = "SLOT_ALREADY_BOOKED"
	BOOKING_NOT_FOUND ErrorCode = "BOOKING_NOT_FOUND"
	FORBIDDEN ErrorCode = "FORBIDDEN"
	SCHEDULE_EXISTS ErrorCode = "SCHEDULE_EXISTS"
	INTERNAL_ERROR ErrorCode = "INTERNAL_ERROR"

)
type ErrorResponse struct {
    Error Error `json:"error"`
} 
type Error struct {
	Code ErrorCode `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrInvalidRequest = Error{
		Code:    INVALID_REQUEST,
		Message: "invalid request",
	}
	
	ErrUnauthorized = Error{
		Code:    UNAUTHORIZED,
		Message: "unauthorized access",
	}
	
	ErrNotFound = Error{
		Code:    NOT_FOUND,
		Message: "resource not found",
	}
	
	ErrForbidden = Error{
		Code:    FORBIDDEN,
		Message: "access forbidden",
	}
	
	ErrInternalError = Error{
		Code:    INTERNAL_ERROR,
		Message: "internal server error",
	}
	
	ErrRoomNotFound = Error{
		Code:    ROOM_NOT_FOUND,
		Message: "room not found",
	}
	
	ErrSlotNotFound = Error{
		Code:    SLOT_NOT_FOUND,
		Message: "slot not found",
	}
	
	ErrSlotAlreadyBooked = Error{
		Code:    SLOT_ALREADY_BOOKED,
		Message: "slot is already booked",
	}
	
	ErrBookingNotFound = Error{
		Code:    BOOKING_NOT_FOUND,
		Message: "booking not found",
	}
	
	ErrScheduleExists = Error{
		Code:    SCHEDULE_EXISTS,
		Message: "schedule already exists",
	}
)
