package service

import (
	"context"
	"log"
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	"time"

	"github.com/google/uuid"
)

// ДОБАВИТЬ КРОН ДЖОБ !!!!!
func (r *RoomService) AddSchedule(ctx context.Context, req request.ScheduleRequest, room_id uuid.UUID) (models.Schedule, error) {
	if !r.room.CheckRoomExist(ctx, room_id) {
		return models.Schedule{}, models.ErrRoomNotFound
	}

	if r.schedule.ScheduleExists(ctx, room_id) {
		return models.Schedule{}, models.ErrScheduleExists
	}
	model := models.Schedule{ID: uuid.New(), RoomID: room_id, DaysOfWeek: req.DaysOfWeek, StartTime: req.StartTime, EndTime: req.EndTime}

	err := r.schedule.CreateSchedule(ctx, model)
	if err != nil {
		log.Println("error adding schedule to db")
		return models.Schedule{}, err
	}

	slots := GenerateSlots(model, time.Now(), 30)

	err = r.slot.InsertSlots(ctx, slots)
	if err != nil {
		log.Println("error with repo inserts slots")
		return models.Schedule{}, err
	}
	return model, nil
}

func GenerateSlots(req models.Schedule, today time.Time, days int) []models.Slot{
	var slots []models.Slot

	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, i)
	
		weekday := int(date.Weekday()) // 0,1,2,3,4,5,6
		if weekday == 0 {
			weekday = 7 // 1,2,3,4,5,6,7
		}

		if !containsDay(req.DaysOfWeek, weekday) {
			continue
		}

		start, _ := time.Parse("15:04", req.StartTime)
		end, _ := time.Parse("15:04", req.EndTime)

		slotStart := time.Date(
			date.Year(), date.Month(), date.Day(),
			start.Hour(), start.Minute(), 0, 0, time.UTC,
		)

		slotEnd := time.Date(
			date.Year(), date.Month(), date.Day(),
			end.Hour(), end.Minute(), 0, 0, time.UTC,
		)

		for slotStart.Add(30 * time.Minute).Before(slotEnd) || 
			slotStart.Add(30 * time.Minute).Equal(slotEnd) {
				slots = append(slots, models.Slot{
					ID: uuid.New(),
					RoomID: req.RoomID,
					Start: slotStart,
					End: slotStart.Add(30 * time.Minute),
				})

				slotStart = slotStart.Add(30 * time.Minute)
			}
	}

	return slots
}

func containsDay(days []int, day int) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}