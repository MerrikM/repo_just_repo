package service

import (
	"context"
	"database/sql"
	_ "errors"
	"repo_just_repo/internal/dto"
	"repo_just_repo/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для CalendarRepository
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindByID(ctx context.Context, id int64) (*model.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepo) Create(ctx context.Context, ev *model.Event) (int64, error) {
	args := m.Called(ctx, ev)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) Update(ctx context.Context, ev *model.Event) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepo) FindForRange(ctx context.Context, userID int64, from, to time.Time) ([]model.Event, error) {
	args := m.Called(ctx, userID, from, to)
	return args.Get(0).([]model.Event), args.Error(1)
}

func TestCreateEvent_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewCalendarService(mockRepo)

	req := dto.CreateEventRequest{
		UserID: 1,
		Date:   "2025-10-01",
		Text:   "Test event",
	}

	ev, _ := req.ValidateAndToEvent()
	mockRepo.On("Create", mock.Anything, &ev).Return(int64(42), nil)

	id, err := svc.CreateEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)
	mockRepo.AssertExpectations(t)
}

func TestCreateEvent_InvalidDate(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewCalendarService(mockRepo)

	req := dto.CreateEventRequest{
		UserID: 1,
		Date:   "bad-date",
		Text:   "Test",
	}

	id, err := svc.CreateEvent(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
}

func TestUpdateEvent_NotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewCalendarService(mockRepo)

	req := dto.UpdateEventRequest{
		ID:     999,
		UserID: 1,
		Date:   "2025-10-01",
		Text:   "Update",
	}
	ev, _ := req.ValidateAndToEvent()

	mockRepo.On("Update", mock.Anything, &ev).Return(sql.ErrNoRows)

	_, err := svc.UpdateEvent(context.Background(), req)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDeleteEvent_InvalidID(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewCalendarService(mockRepo)

	err := svc.DeleteEvent(context.Background(), dto.DeleteEventRequest{ID: 0})
	assert.Error(t, err)
	assert.Equal(t, "id is required", err.Error())
}

func TestDeleteEvent_NotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewCalendarService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(123)).Return(sql.ErrNoRows)

	err := svc.DeleteEvent(context.Background(), dto.DeleteEventRequest{ID: 123})
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestEventsForDay(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewCalendarService(mockRepo)

	date := time.Date(2025, 10, 1, 12, 0, 0, 0, time.UTC)
	expectedEvents := []model.Event{
		{ID: 1, UserID: 1, Date: date, Text: "Meeting"},
	}

	from := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	to := from
	mockRepo.On("FindForRange", mock.Anything, int64(1), from, to).Return(expectedEvents, nil)

	events, err := svc.EventsForDay(context.Background(), 1, date)

	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "Meeting", events[0].Text)
	mockRepo.AssertExpectations(t)
}
