package handlers

import (
	"encoding/json"
	"net/http"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	respond "rooms_service/internal/models/dto/response"
)

// @Summary      Регистрация пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body object{email=string,password=string,role=string} true "Данные пользователя"
// @Success      201  {object}  object{user=models.User}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /register [post]
func (h *Handler) Registrate(w http.ResponseWriter, r *http.Request){
	var req request.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, models.ErrInvalidRequest)
		return
	}
	
	user, err := h.service.RegisterUser(r.Context(), req)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, 201, respond.UserResponse{User: user})
}

// @Summary      Авторизация по email и паролю
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body object{email=string,password=string} true "Учётные данные"
// @Success      200  {object}  object{token=string}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	err :=json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}

	token, err := h.service.LoginUser(r.Context(), req)
	if err != nil {
		HandleError(w, err)
		return
	}
	WriteJSON(w, 201, token)
}