package handlers


import (
	"encoding/json"
	"net/http"
	"rooms_service/internal/models"
    "fmt"
	
)
func HandleError(w http.ResponseWriter, err error) {
    switch err {
    case models.ErrSlotNotFound:
        WriteError(w, 404, models.ErrSlotNotFound)
    case models.ErrRoomNotFound:
        WriteError(w, 404, models.ErrRoomNotFound)
    case models.ErrBookingNotFound:
        WriteError(w, 404, models.ErrBookingNotFound)
    case models.ErrSlotAlreadyBooked:
        WriteError(w, 409, models.ErrSlotAlreadyBooked)
    case models.ErrScheduleExists:
        WriteError(w, 409, models.ErrScheduleExists)
    case models.ErrInvalidRequest:
        WriteError(w, 400, models.ErrInvalidRequest)
    case models.ErrUnauthorized:
        WriteError(w, 401, models.ErrUnauthorized)
    case models.ErrForbidden:
        WriteError(w, 403, models.ErrForbidden)
    default:
        WriteError(w, 500, models.ErrInternalError)
    }
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        // Ошибка кодирования — уже отправили заголовки
        // Можно только логировать
        fmt.Printf("Failed to encode JSON: %v", err)
    }
}

 func WriteError(w http.ResponseWriter, status int, err models.Error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]any{
        "error": err,
    })
}
