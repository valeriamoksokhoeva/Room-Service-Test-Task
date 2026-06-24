package service_test

import (
	"context"
	"errors"
	"testing"
	"time"
	
	"rooms_service/internal/models"
	"rooms_service/internal/models/dto/request"
	"rooms_service/internal/service"
	
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUserToDB(ctx context.Context, user models.UserDB) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models.UserDB, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.UserDB), args.Error(1)
}

func (m *MockUserRepository) CheckEmailExist(ctx context.Context, email string) bool {
	args := m.Called(ctx, email)
	return args.Bool(0)
}

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) GetRooms(ctx context.Context) ([]models.Room, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Room), args.Error(1)
}

func (m *MockRoomRepository) CreateRoom(ctx context.Context, room models.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) CheckRoomExist(ctx context.Context, roomID uuid.UUID) bool {
	args := m.Called(ctx, roomID)
	return args.Bool(0)
}

type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) CreateSchedule(ctx context.Context, schedule models.Schedule) error {
	args := m.Called(ctx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) ScheduleExists(ctx context.Context, roomID uuid.UUID) bool {
	args := m.Called(ctx, roomID)
	return args.Bool(0)
}

func (m *MockScheduleRepository) GetAllSchedules(ctx context.Context) ([]models.Schedule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Schedule), args.Error(1)
}

type MockSlotRepository struct {
	mock.Mock
}

func (m *MockSlotRepository) InsertSlots(ctx context.Context, slots []models.Slot) error {
	args := m.Called(ctx, slots)
	return args.Error(0)
}

func (m *MockSlotRepository) GetAvailableSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]models.Slot, error) {
	args := m.Called(ctx, roomID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Slot), args.Error(1)
}

func (m *MockSlotRepository) CheckAvailableSlot(ctx context.Context, slotID uuid.UUID) bool {
	args := m.Called(ctx, slotID)
	return args.Bool(0)
}

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) GetAllBookings(ctx context.Context, page int, pageSize int) ([]models.Booking, int, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Booking), args.Int(1), args.Error(2)
}

func (m *MockBookingRepository) InsertBooking(ctx context.Context, booking models.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepository) MyBookings(ctx context.Context, userID uuid.UUID) ([]models.Booking, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (models.Booking, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(models.Booking), args.Error(1)
}

func (m *MockBookingRepository) UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID) error {
	args := m.Called(ctx, bookingID)
	return args.Error(0)
}

// TestDomainService_AddRoom - тест создания комнаты
func TestDomainService_AddRoom(t *testing.T) {
	tests := []struct {
		name        string
		req         request.RoomRequest
		setupMocks  func(*MockRoomRepository)
		expectedErr error
	}{
		{
			name: "успешное создание комнаты",
			req: request.RoomRequest{
				Name:        "Переговорка А",
				Description: "Большая переговорка",
				Capacity:    10,
			},
			setupMocks: func(r *MockRoomRepository) {
				r.On("CreateRoom", mock.Anything, mock.AnythingOfType("models.Room")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "ошибка при создании комнаты",
			req: request.RoomRequest{
				Name:     "Переговорка А",
				Capacity: 10,
			},
			setupMocks: func(r *MockRoomRepository) {
				r.On("CreateRoom", mock.Anything, mock.AnythingOfType("models.Room")).Return(errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomRepo := new(MockRoomRepository)
			tt.setupMocks(mockRoomRepo)

			svc := service.NewService(
				mockRoomRepo,
				nil,
				nil, 
				nil, 
				nil, 
				nil,
			)

			room, err := svc.AddRoom(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Name, room.Name)
				assert.Equal(t, tt.req.Capacity, room.Capacity)
				assert.NotEmpty(t, room.ID)
			}

			mockRoomRepo.AssertExpectations(t)
		})
	}
}

// TestDomainService_GetRoomsList - тест получения списка комнат
func TestDomainService_GetRoomsList(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*MockRoomRepository)
		expectedLen int
		expectedErr error
	}{
		{
			name: "успешное получение списка комнат",
			setupMocks: func(r *MockRoomRepository) {
				rooms := []models.Room{
					{ID: uuid.New(), Name: "Комната 1", Capacity: 10},
					{ID: uuid.New(), Name: "Комната 2", Capacity: 5},
				}
				r.On("GetRooms", mock.Anything).Return(rooms, nil)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name: "пустой список комнат",
			setupMocks: func(r *MockRoomRepository) {
				r.On("GetRooms", mock.Anything).Return([]models.Room{}, nil)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name: "ошибка при получении комнат",
			setupMocks: func(r *MockRoomRepository) {
				r.On("GetRooms", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedLen: 0,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomRepo := new(MockRoomRepository)
			tt.setupMocks(mockRoomRepo)

			svc := service.NewService(
				mockRoomRepo, nil, nil, nil, nil, nil,
			)

			rooms, err := svc.GetRoomsList(context.Background())

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, rooms, tt.expectedLen)
			}

			mockRoomRepo.AssertExpectations(t)
		})
	}
}

// TestDomainService_AddSchedule - тест создания расписания
func TestDomainService_AddSchedule(t *testing.T) {
	roomID := uuid.New()

	tests := []struct {
		name        string
		roomID      uuid.UUID
		req         request.ScheduleRequest
		setupMocks  func(*MockRoomRepository, *MockScheduleRepository, *MockSlotRepository)
		expectedErr error
	}{
		{
			name:   "успешное создание расписания",
			roomID: roomID,
			req: request.ScheduleRequest{
				DaysOfWeek: []int{1, 2, 3, 4, 5},
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			setupMocks: func(r *MockRoomRepository, s *MockScheduleRepository, sl *MockSlotRepository) {
				r.On("CheckRoomExist", mock.Anything, roomID).Return(true)
				s.On("ScheduleExists", mock.Anything, roomID).Return(false)
				s.On("CreateSchedule", mock.Anything, mock.AnythingOfType("models.Schedule")).Return(nil)
				sl.On("InsertSlots", mock.Anything, mock.AnythingOfType("[]models.Slot")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "комната не найдена",
			roomID: roomID,
			req: request.ScheduleRequest{
				DaysOfWeek: []int{1, 2, 3},
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			setupMocks: func(r *MockRoomRepository, s *MockScheduleRepository, sl *MockSlotRepository) {
				r.On("CheckRoomExist", mock.Anything, roomID).Return(false)
			},
			expectedErr: models.ErrRoomNotFound,
		},
		{
			name:   "расписание уже существует",
			roomID: roomID,
			req: request.ScheduleRequest{
				DaysOfWeek: []int{1, 2, 3},
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			setupMocks: func(r *MockRoomRepository, s *MockScheduleRepository, sl *MockSlotRepository) {
				r.On("CheckRoomExist", mock.Anything, roomID).Return(true)
				s.On("ScheduleExists", mock.Anything, roomID).Return(true)
			},
			expectedErr: models.ErrScheduleExists,
		},
		{
			name:   "ошибка при создании расписания в БД",
			roomID: roomID,
			req: request.ScheduleRequest{
				DaysOfWeek: []int{1, 2, 3},
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			setupMocks: func(r *MockRoomRepository, s *MockScheduleRepository, sl *MockSlotRepository) {
				r.On("CheckRoomExist", mock.Anything, roomID).Return(true)
				s.On("ScheduleExists", mock.Anything, roomID).Return(false)
				s.On("CreateSchedule", mock.Anything, mock.AnythingOfType("models.Schedule")).Return(errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name:   "ошибка при вставке слотов",
			roomID: roomID,
			req: request.ScheduleRequest{
				DaysOfWeek: []int{1, 2, 3},
				StartTime:  "09:00",
				EndTime:    "18:00",
			},
			setupMocks: func(r *MockRoomRepository, s *MockScheduleRepository, sl *MockSlotRepository) {
				r.On("CheckRoomExist", mock.Anything, roomID).Return(true)
				s.On("ScheduleExists", mock.Anything, roomID).Return(false)
				s.On("CreateSchedule", mock.Anything, mock.AnythingOfType("models.Schedule")).Return(nil)
				sl.On("InsertSlots", mock.Anything, mock.AnythingOfType("[]models.Slot")).Return(errors.New("slot db error"))
			},
			expectedErr: errors.New("slot db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomRepo := new(MockRoomRepository)
			mockScheduleRepo := new(MockScheduleRepository)
			mockSlotRepo := new(MockSlotRepository)
			
			tt.setupMocks(mockRoomRepo, mockScheduleRepo, mockSlotRepo)

			svc := service.NewService(
				mockRoomRepo,     
				nil,
				mockScheduleRepo, 
				mockSlotRepo,     
				nil,
				nil,             
			)

			schedule, err := svc.AddSchedule(context.Background(), tt.req, tt.roomID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.roomID, schedule.RoomID)
				assert.Equal(t, tt.req.DaysOfWeek, schedule.DaysOfWeek)
				assert.Equal(t, tt.req.StartTime, schedule.StartTime)
				assert.Equal(t, tt.req.EndTime, schedule.EndTime)
				assert.NotEmpty(t, schedule.ID)
			}

			mockRoomRepo.AssertExpectations(t)
			mockScheduleRepo.AssertExpectations(t)
			mockSlotRepo.AssertExpectations(t)
		})
	}
}

// TestDomainService_ListSlots - тест получения слотов
func TestDomainService_ListSlots(t *testing.T) {
	roomID := uuid.New()
	today := time.Now().UTC().Truncate(24 * time.Hour)

	tests := []struct {
		name        string
		roomID      uuid.UUID
		date        time.Time
		setupMocks  func(*MockSlotRepository)
		expectedLen int
		expectedErr error
	}{
		{
			name:   "успешное получение слотов",
			roomID: roomID,
			date:   today,
			setupMocks: func(s *MockSlotRepository) {
				slots := []models.Slot{
					{ID: uuid.New(), RoomID: roomID, Start: today.Add(9 * time.Hour), End: today.Add(9*time.Hour + 30*time.Minute)},
					{ID: uuid.New(), RoomID: roomID, Start: today.Add(10 * time.Hour), End: today.Add(10*time.Hour + 30*time.Minute)},
				}
				s.On("GetAvailableSlots", mock.Anything, roomID, today).Return(slots, nil)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:   "нет доступных слотов",
			roomID: roomID,
			date:   today,
			setupMocks: func(s *MockSlotRepository) {
				s.On("GetAvailableSlots", mock.Anything, roomID, today).Return([]models.Slot{}, nil)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name:   "ошибка при получении слотов",
			roomID: roomID,
			date:   today,
			setupMocks: func(s *MockSlotRepository) {
				s.On("GetAvailableSlots", mock.Anything, roomID, today).Return(nil, errors.New("db error"))
			},
			expectedLen: 0,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSlotRepo := new(MockSlotRepository)
			tt.setupMocks(mockSlotRepo)

			svc := service.NewService(
				nil, nil, nil, mockSlotRepo, nil, nil,
			)

			slots, err := svc.ListSlots(context.Background(), tt.roomID, tt.date)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, slots, tt.expectedLen)
			}

			mockSlotRepo.AssertExpectations(t)
		})
	}
}

// TestDomainService_CreateBooking - тест создания бронирования
func TestDomainService_CreateBooking(t *testing.T) {
	slotID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name        string
		req         request.BookingRequest
		userID      uuid.UUID
		setupMocks  func(*MockSlotRepository, *MockBookingRepository)
		expectedErr error
	}{
		{
			name: "успешное создание бронирования",
			req: request.BookingRequest{
				SlotId: slotID,
			},
			userID: userID,
			setupMocks: func(s *MockSlotRepository, b *MockBookingRepository) {
				s.On("CheckAvailableSlot", mock.Anything, slotID).Return(true)
				b.On("InsertBooking", mock.Anything, mock.AnythingOfType("models.Booking")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "слот недоступен",
			req: request.BookingRequest{
				SlotId: slotID,
			},
			userID: userID,
			setupMocks: func(s *MockSlotRepository, b *MockBookingRepository) {
				s.On("CheckAvailableSlot", mock.Anything, slotID).Return(false)
			},
			expectedErr: models.ErrSlotNotFound,
		},
		{
			name: "ошибка при сохранении в БД",
			req: request.BookingRequest{
				SlotId: slotID,
			},
			userID: userID,
			setupMocks: func(s *MockSlotRepository, b *MockBookingRepository) {
				s.On("CheckAvailableSlot", mock.Anything, slotID).Return(true)
				b.On("InsertBooking", mock.Anything, mock.AnythingOfType("models.Booking")).Return(errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSlotRepo := new(MockSlotRepository)
			mockBookingRepo := new(MockBookingRepository)
			tt.setupMocks(mockSlotRepo, mockBookingRepo)

			svc := service.NewService(
				nil, nil, nil, mockSlotRepo, mockBookingRepo, nil,
			)

			booking, err := svc.CreateBooking(context.Background(), tt.req, tt.userID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, booking.UserID)
				assert.Equal(t, tt.req.SlotId, booking.SlotID)
				assert.Equal(t, models.Active, booking.Status)
				assert.NotEmpty(t, booking.ID)
			}

			mockSlotRepo.AssertExpectations(t)
			mockBookingRepo.AssertExpectations(t)
		})
	}
}


// TestDomainService_GetMyBookings - тест получения бронирований пользователя
func TestDomainService_GetMyBookings(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		setupMocks  func(*MockBookingRepository)
		expectedLen int
		expectedErr error
	}{
		{
			name:   "успешное получение бронирований",
			userID: userID,
			setupMocks: func(b *MockBookingRepository) {
				bookings := []models.Booking{
					{ID: uuid.New(), UserID: userID, Status: models.Active},
					{ID: uuid.New(), UserID: userID, Status: models.Active},
				}
				b.On("MyBookings", mock.Anything, userID).Return(bookings, nil)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:   "нет бронирований",
			userID: userID,
			setupMocks: func(b *MockBookingRepository) {
				b.On("MyBookings", mock.Anything, userID).Return([]models.Booking{}, nil)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name:   "ошибка при получении",
			userID: userID,
			setupMocks: func(b *MockBookingRepository) {
				b.On("MyBookings", mock.Anything, userID).Return(nil, errors.New("db error"))
			},
			expectedLen: 0,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBookingRepo := new(MockBookingRepository)
			tt.setupMocks(mockBookingRepo)

			svc := service.NewService(
				nil, nil, nil, nil, mockBookingRepo, nil,
			)

			bookings, err := svc.GetMyBookings(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, bookings, tt.expectedLen)
			}

			mockBookingRepo.AssertExpectations(t)
		})
	}
}

// TestDomainService_CancelMyBooking - тест отмены бронирования
func TestDomainService_CancelMyBooking(t *testing.T) {
	bookingID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name        string
		bookingID   uuid.UUID
		userID      uuid.UUID
		setupMocks  func(*MockBookingRepository)
		expectedErr error
	}{
		{
			name:      "успешная отмена",
			bookingID: bookingID,
			userID:    userID,
			setupMocks: func(b *MockBookingRepository) {
				booking := models.Booking{
					ID:     bookingID,
					UserID: userID,
					Status: models.Active,
				}
				b.On("GetBookingByID", mock.Anything, bookingID).Return(booking, nil)
				b.On("UpdateBookingStatus", mock.Anything, bookingID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:      "попытка отменить чужое бронирование",
			bookingID: bookingID,
			userID:    otherUserID,
			setupMocks: func(b *MockBookingRepository) {
				booking := models.Booking{
					ID:     bookingID,
					UserID: userID,
					Status: models.Active,
				}
				b.On("GetBookingByID", mock.Anything, bookingID).Return(booking, nil)
			},
			expectedErr: models.ErrForbidden,
		},
		{
			name:      "отмена уже отменённого (идемпотентность)",
			bookingID: bookingID,
			userID:    userID,
			setupMocks: func(b *MockBookingRepository) {
				booking := models.Booking{
					ID:     bookingID,
					UserID: userID,
					Status: models.Cancelled,
				}
				b.On("GetBookingByID", mock.Anything, bookingID).Return(booking, nil)
			},
			expectedErr: nil,
		},
		{
			name:      "бронирование не найдено",
			bookingID: uuid.New(),
			userID:    userID,
			setupMocks: func(b *MockBookingRepository) {
				b.On("GetBookingByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(models.Booking{}, models.ErrBookingNotFound)
			},
			expectedErr: models.ErrBookingNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBookingRepo := new(MockBookingRepository)
			tt.setupMocks(mockBookingRepo)

			svc := service.NewService(
				nil, nil, nil, nil, mockBookingRepo, nil,
			)

			booking, err := svc.CancelMyBooking(context.Background(), tt.bookingID, tt.userID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, models.Cancelled, booking.Status)
			}

			mockBookingRepo.AssertExpectations(t)
		})
	}
}

// TestDomainService_AllBookings - тест получения всех бронирований с пагинацией
func TestDomainService_AllBookings(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		pageSize      int
		setupMocks    func(*MockBookingRepository)
		expectedLen   int
		expectedTotal int
		expectedErr   error
	}{
		{
			name:     "успешное получение с пагинацией",
			page:     1,
			pageSize: 10,
			setupMocks: func(b *MockBookingRepository) {
				bookings := []models.Booking{
					{ID: uuid.New(), Status: models.Active},
					{ID: uuid.New(), Status: models.Active},
				}
				b.On("GetAllBookings", mock.Anything, 1, 10).Return(bookings, 25, nil)
			},
			expectedLen:   2,
			expectedTotal: 25,
			expectedErr:   nil,
		},
		{
			name:     "пустой список",
			page:     1,
			pageSize: 10,
			setupMocks: func(b *MockBookingRepository) {
				b.On("GetAllBookings", mock.Anything, 1, 10).Return([]models.Booking{}, 0, nil)
			},
			expectedLen:   0,
			expectedTotal: 0,
			expectedErr:   nil,
		},
		{
			name:     "ошибка при получении",
			page:     1,
			pageSize: 10,
			setupMocks: func(b *MockBookingRepository) {
				b.On("GetAllBookings", mock.Anything, 1, 10).Return(nil, 0, errors.New("db error"))
			},
			expectedLen:   0,
			expectedTotal: 0,
			expectedErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBookingRepo := new(MockBookingRepository)
			tt.setupMocks(mockBookingRepo)

			svc := service.NewService(
				nil, nil, nil, nil, mockBookingRepo, nil,
			)

			bookings, total, err := svc.AllBookings(context.Background(), tt.page, tt.pageSize)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, bookings, tt.expectedLen)
				assert.Equal(t, tt.expectedTotal, total)
			}

			mockBookingRepo.AssertExpectations(t)
		})
	}
}

