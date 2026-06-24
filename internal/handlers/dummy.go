package handlers

import (
	"encoding/json"
	"net/http"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	"rooms_service/internal/service"

	"github.com/google/uuid"
)

const (
    AdminUUID = "a0000000-0000-0000-0000-000000000001"
    UserUUID  = "b0000000-0000-0000-0000-000000000002"
)

// @Summary      Получить тестовый JWT по роли
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body object{role=string} true "Роль: admin или user"
// @Success      200  {object}  object{token=string}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /dummyLogin [post]
func (h *Handler)DummyLogin(w http.ResponseWriter, r *http.Request) {
	var dummy request.DummyRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&dummy)

	if err != nil {
		HandleError(w, models.ErrInvalidRequest)
		return
	}
	var userId uuid.UUID
	switch dummy.Role {
	case models.Admin:
		userId, _ = uuid.Parse(AdminUUID)
	case models.UserT:
		userId, _ = uuid.Parse(UserUUID)
	default:
		HandleError(w, models.ErrInvalidRequest)
		return
	}

	tokenResponse, errCode := service.GenerateToken(userId, dummy.Role)
	if errCode != nil {
		HandleError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, tokenResponse)
}

