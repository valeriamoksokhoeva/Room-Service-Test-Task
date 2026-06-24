package cron

import (
	"context"
	"rooms_service/internal/repository"
	"rooms_service/internal/service"
	"time"
	"log"
)

func StartCronGenerationSlot(ctx context.Context, schedule repository.ScheduleRepository, slot repository.SlotRepository) {
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				generateFutureSlots(ctx, schedule, slot)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func generateFutureSlots(ctx context.Context, schedule repository.ScheduleRepository, slot repository.SlotRepository) {
	schedules, err := schedule.GetAllSchedules(ctx)
	if err != nil {
		return
	}

	newDay := time.Now().UTC().AddDate(0,0,30)
	for _, s := range schedules {
		slots := service.GenerateSlots(s, newDay, 1)
		err = slot.InsertSlots(ctx, slots)
		if err != nil {
			log.Printf("cron: insert slots error: %v", err)
		}
	}
}