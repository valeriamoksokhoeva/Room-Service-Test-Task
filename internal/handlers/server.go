package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rooms_service/internal/service"
	"syscall"
	"time"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "rooms_service/docs"
	"github.com/gorilla/mux"
)

type Handler struct {
    service service.DomainService
}

func NewHandler(s service.DomainService) Handler {
	return Handler{service: s}
}
func Router(h Handler) *mux.Router  {
    m := mux.NewRouter()

	m.HandleFunc("/register", h.Registrate).Methods("POST")
	m.HandleFunc("/login", h.Login).Methods("POST")
	m.HandleFunc("/dummyLogin", h.DummyLogin).Methods("POST")


	m.HandleFunc("/rooms/list", AuthMiddleware(h.ListRooms)).Methods("GET")
	m.HandleFunc("/rooms/create", AuthMiddleware(RequireRole("admin", h.CreateRoomByAdmin))).Methods("POST")

	m.HandleFunc("/rooms/{roomId}/schedule/create", AuthMiddleware(RequireRole("admin", h.CreateScheduleByAdmin))).Methods("POST")

	m.HandleFunc("/rooms/{roomId}/slots/list", AuthMiddleware(h.GetSlots)).Methods("GET")

	m.HandleFunc("/bookings/create", AuthMiddleware(RequireRole("user", h.AddBooking))).Methods("POST")
	m.HandleFunc("/bookings/my", AuthMiddleware(RequireRole("user", h.MyBookings))).Methods("GET")
	m.HandleFunc("/bookings/{bookingId}/cancel", AuthMiddleware(RequireRole("user", h.DeleteMyBooking))).Methods("POST")
	m.HandleFunc("/bookings/list", AuthMiddleware(RequireRole("admin", h.ListAllBookings))).Methods("GET")

	m.HandleFunc("/_info", h.HealthCheck).Methods("GET")
	m.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return m
}

func StartServer(r *mux.Router) {
	server := &http.Server{
		Addr: ":8080",
		Handler: r,

	}

	go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Errorf("Server failed: %v", err)
        }
    }()

	quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // 8. Даем время завершить текущие запросы
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped")
}

// @Summary      Health check
// @Tags         System
// @Success      200
// @Router       /_info [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}