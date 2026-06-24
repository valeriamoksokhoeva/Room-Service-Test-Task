// @title           Room Booking Service
// @version         1.0
// @description     Сервис бронирования переговорок
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log"
	"os"
	"rooms_service/internal/conference"
	"rooms_service/internal/cron"
	"rooms_service/internal/db"
	"rooms_service/internal/handlers"
	"rooms_service/internal/repository"
	"rooms_service/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found (running in Docker?)")
    }
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connString := os.Getenv("DATABASE_URL")
	conn, err := db.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer conn.Close()

	room := repository.NewRoomRepository(conn)
	user := repository.NewUserRepository(conn)
	schedule := repository.NewScheduleRepository(conn)
	slot := repository.NewSlotRepository(conn)
	booking := repository.NewBookingRepository(conn)

	conf := conference.NewMockConferenceService()
	service := service.NewService(room, user, schedule, slot, booking, conf)
	handler := handlers.NewHandler(service)

	r := handlers.Router(handler)
	cron.StartCronGenerationSlot(ctx, schedule, slot)

	handlers.StartServer(r)
}