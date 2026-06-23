package repository

import (
	"context"
	"log"
	"rooms_service/internal/models"

	"github.com/google/uuid"
)

func (rc *ScheduleRepo) CreateSchedule(ctx context.Context,req models.Schedule) error {
	query := `
		INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := rc.conn.Exec(ctx, query, req.ID, req.RoomID, req.DaysOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return models.ErrInternalError
	}

	return nil
}

func (rc *ScheduleRepo) ScheduleExists(ctx context.Context, room_id uuid.UUID) bool {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM schedules
			WHERE room_id = $1
		)
	`

	var check bool
	err1 := rc.conn.QueryRow(ctx, query, room_id).Scan(&check)
	if err1 != nil {
		return false
	}
	return check
}

func (rc *ScheduleRepo) GetAllSchedules(ctx context.Context) ([]models.Schedule, error) {
	query := `
		SELECT id, room_id, days_of_week, start_time, end_time FROM schedules
	`
	var schedules []models.Schedule
	rows, err := rc.conn.Query(ctx, query)
	if err != nil {
		log.Println("error doing query to get all schedules")
		return nil, models.ErrInternalError
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(&schedule.ID, &schedule.RoomID, &schedule.DaysOfWeek, &schedule.StartTime, &schedule.EndTime)
		if err != nil {
			log.Println("error rows.scan")
			return nil, models.ErrInternalError
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		log.Println("rows.err()")
		return nil, models.ErrInternalError
	}

	return schedules, nil
}