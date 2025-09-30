package dto

import (
	"errors"
	"repo_just_repo/internal/model"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

// ValidateAndToEvent ValidateCreate проверяет поля и парсит дату
func (r *CreateEventRequest) ValidateAndToEvent() (model.Event, error) {
	if r.UserID <= 0 {
		return model.Event{}, errors.New("user_id is required and must be > 0")
	}
	if strings.TrimSpace(r.Text) == "" {
		return model.Event{}, errors.New("event text is required")
	}
	t, err := time.Parse(DateLayout, r.Date)
	if err != nil {
		return model.Event{}, errors.New("date must be in YYYY-MM-DD format")
	}
	// Нормализуем время 00:00:00 UTC (храним date-only)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return model.Event{
		UserID: r.UserID,
		Date:   t,
		Text:   r.Text,
		Active: true,
	}, nil
}

// ValidateAndToEvent Для Update: проверяем id и, при необходимости, возвращаем Event с ID
func (r *UpdateEventRequest) ValidateAndToEvent() (model.Event, error) {
	if r.ID <= 0 {
		return model.Event{}, errors.New("id is required")
	}
	ev, err := (&CreateEventRequest{UserID: r.UserID, Date: r.Date, Text: r.Text}).ValidateAndToEvent()
	if err != nil {
		return model.Event{}, err
	}
	ev.ID = r.ID
	return ev, nil
}
