package handlers

import (
	"encoding/json"
	"net/http"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	respond "rooms_service/internal/models/dto/response"
	"time"
	"log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary      Создать расписание переговорки (только admin)
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        roomId  path  string  true  "ID переговорки"
// @Param        body    body  models.Schedule  true  "Расписание"
// @Success      201  {object}  object{schedule=models.Schedule}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /rooms/{roomId}/schedule/create [post]
func (h *Handler) CreateScheduleByAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    log.Printf("ALL VARS: %+v", vars)
    room_id_str := vars["roomId"]
    log.Printf("roomId string: '%s'", room_id_str)
	room_id, err := uuid.Parse(room_id_str)
	if err != nil || room_id_str == ""{
		HandleError(w, models.ErrInvalidRequest)
		return
	}
	var req request.ScheduleRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, models.ErrInvalidRequest)
		return
	}
	
	schedule, err := h.service.AddSchedule(r.Context(), req, room_id)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, 201, respond.ScheduleResponse{Schedule: schedule})
}


// @Summary      Список доступных слотов по переговорке и дате
// @Tags         Slots
// @Produce      json
// @Security     BearerAuth
// @Param        roomId  path   string  true   "ID переговорки"
// @Param        date    query  string  true   "Дата в формате YYYY-MM-DD"
// @Success      200  {object}  object{slots=[]models.Slot}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /rooms/{roomId}/slots/list [get]
func (h *Handler) GetSlots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room_id_str := vars["roomId"]

	room_id, err := uuid.Parse(room_id_str)
	if err != nil {
       	HandleError(w, models.ErrInvalidRequest)
        return
    }

	dateStr := r.URL.Query().Get("date")
    if dateStr == "" {
        HandleError(w, models.ErrInvalidRequest)
        return
    }

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
        HandleError(w, models.ErrInvalidRequest)
        return
    }

	slots, err := h.service.ListSlots(r.Context(), room_id, date)
	if err != nil {
		HandleError(w, err)
		return  
	}
	WriteJSON(w, 200, respond.SlotListResponse{Slots: slots})
}