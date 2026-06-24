package conference

import (
	"fmt"
	"log"
	"math/rand"
	"rooms_service/internal/models"
	"time"

	"github.com/google/uuid"
)

type ConferenceService interface {
	CreateLink(bookingID uuid.UUID) (string, error)
}

type MockConferenceService struct{}

func NewMockConferenceService() *MockConferenceService {
	return &MockConferenceService{}
}

func (s *MockConferenceService) CreateLink(bookingID uuid.UUID) (string, error) {
	time.Sleep(100 * time.Millisecond)

	if rand.Intn(10) == 0 {
		log.Println("conference service unavailable")
		return "", models.ErrNotFound
	}

	return fmt.Sprintf("https://meet.example.com/%s", bookingID.String()), nil
}