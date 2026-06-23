package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary      Создать бронь (только user)
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object{slotId=string,createConferenceLink=bool} true "Данные брони"
// @Success      201  {object}  object{booking=models.Booking}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /bookings/create [post]
func (h *Handler) AddBooking(w http.ResponseWriter, r *http.Request) {
	var req request.BookingRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		HandleError(w, models.ErrInvalidRequest)
		return
	}

	user_id, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		HandleError(w, models.ErrUnauthorized)
		return
	}

	booking, err := h.service.CreateBooking(r.Context(), req, user_id)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, 201, booking)
}


// @Summary      Мои брони (только user)
// @Tags         Bookings
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{bookings=[]models.Booking}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /bookings/my [get]
func (h *Handler) MyBookings(w http.ResponseWriter, r *http.Request) {
	user_id, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		HandleError(w, models.ErrUnauthorized)
		return
	}

	bookings, err := h.service.GetMyBookings(r.Context(), user_id)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, 200, bookings)
}


// @Summary      Отменить бронь (только user, только своя)
// @Tags         Bookings
// @Produce      json
// @Security     BearerAuth
// @Param        bookingId  path  string  true  "ID брони"
// @Success      200  {object}  object{booking=models.Booking}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /bookings/{bookingId}/cancel [post]
func (h *Handler) DeleteMyBooking(w http.ResponseWriter, r *http.Request) {
	user_id, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		HandleError(w, models.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	booking_id_str := vars["bookingId"]
	booking_id, err := uuid.Parse(booking_id_str)
	if booking_id_str == "" || err != nil {
		log.Printf("no booking id error: %v", err)
		HandleError(w, models.ErrInvalidRequest)
	}
	booking, err := h.service.CancelMyBooking(r.Context(), booking_id, user_id)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, 200, booking)
}


// @Summary      Список всех броней с пагинацией (только admin)
// @Tags         Bookings
// @Produce      json
// @Security     BearerAuth
// @Param        page      query  int  false  "Номер страницы (default 1)"
// @Param        pageSize  query  int  false  "Размер страницы (default 20, max 100)"
// @Success      200  {object}  object{bookings=[]models.Booking,pagination=models.Pagination}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /bookings/list [get]
func (h *Handler) ListAllBookings(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 20

	if p := r.URL.Query().Get("page"); p != "" {
		val, err := strconv.Atoi(p)
		if err != nil || val < 1 {
			HandleError(w, models.ErrInvalidRequest)
			return
		}
		page = val
	}

	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		val, err := strconv.Atoi(ps)
		if err != nil || val < 1 || val > 100 {
			HandleError(w, models.ErrInvalidRequest)
			return
		}
		pageSize = val
	}

	bookings, total, err := h.service.AllBookings(r.Context(), page, pageSize)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, 200, map[string]any{
		"bookings": bookings,
		"pagination": models.Pagination{
			Page: page,
			PageSize: pageSize,
			Total: total,
		},
	})
}