package handlers

import (
	"encoding/json"
	"net/http"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	respond "rooms_service/internal/models/dto/response"
)

// @Summary      Список переговорок
// @Tags         Rooms
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{rooms=[]models.Room}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /rooms/list [get]
func (s *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := s.service.GetRoomsList(r.Context())
	if err != nil {
		HandleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, respond.RoomsListResponse{Rooms: rooms})
}


// @Summary      Создать переговорку (только admin)
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object{name=string,description=string,capacity=int} true "Данные переговорки"
// @Success      201  {object}  object{room=models.Room}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /rooms/create [post]
func (s *Handler) CreateRoomByAdmin(w http.ResponseWriter, r *http.Request) {
	var req request.RoomRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, models.ErrInvalidRequest)
		return
	}
	room, err := s.service.AddRoom(r.Context(), req)
	if err != nil {
		HandleError(w, err)
		return
	}
	WriteJSON(w, 201, respond.RoomResponse{Room: room})
}