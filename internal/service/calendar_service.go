package service

import (
	"context"
	"database/sql"
	"errors"
	"repo_just_repo/internal/dto"
	"repo_just_repo/internal/model"
	"repo_just_repo/internal/ports"
	"time"
)

type CalendarService struct {
	repo ports.CalendarRepository
}

func NewCalendarService(repo ports.CalendarRepository) *CalendarService {
	return &CalendarService{repo: repo}
}

// CreateEvent Создать событие
func (s *CalendarService) CreateEvent(ctx context.Context, req dto.CreateEventRequest) (int64, error) {
	ev, err := req.ValidateAndToEvent()
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, &ev)
}

// UpdateEvent Обновить событие
func (s *CalendarService) UpdateEvent(ctx context.Context, req dto.UpdateEventRequest) (*model.Event, error) {
	ev, err := req.ValidateAndToEvent()
	if err != nil {
		return nil, err
	}
	err = s.repo.Update(ctx, &ev)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &ev, err
}

// DeleteEvent Удалить событие
func (s *CalendarService) DeleteEvent(ctx context.Context, req dto.DeleteEventRequest) error {
	if req.ID <= 0 {
		return errors.New("id is required")
	}
	err := s.repo.Delete(ctx, req.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

// EventsForDay События на день
func (s *CalendarService) EventsForDay(ctx context.Context, userID int64, date time.Time) ([]model.Event, error) {
	from := normalizeDate(date)
	to := from
	return s.repo.FindForRange(ctx, userID, from, to)
}

// EventsForWeek События на неделю
func (s *CalendarService) EventsForWeek(ctx context.Context, userID int64, date time.Time) ([]model.Event, error) {
	from := normalizeDate(date)
	to := from.AddDate(0, 0, 6)
	return s.repo.FindForRange(ctx, userID, from, to)
}

// EventsForMonth События на месяц
func (s *CalendarService) EventsForMonth(ctx context.Context, userID int64, date time.Time) ([]model.Event, error) {
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return s.repo.FindForRange(ctx, userID, from, to)
}

// normalizeDate приводит время к 00:00:00 UTC
func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// ErrNotFound используется для бизнес-логики (например, при попытке удалить несуществующее событие)
var ErrNotFound = errors.New("event not found")
