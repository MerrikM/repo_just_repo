package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"repo_just_repo/internal/dto"
	"repo_just_repo/internal/service"
	"time"
)

type EventHandler struct {
	svc *service.CalendarService
}

func NewEventHandler(svc *service.CalendarService) *EventHandler {
	return &EventHandler{svc: svc}
}

// CreateEvent POST /create_event
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	event, err := h.svc.CreateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": event})
}

// UpdateEvent POST /update_event
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	event, err := h.svc.UpdateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": event})
}

// DeleteEvent POST /delete_event
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteEvent(r.Context(), req); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": "deleted"})
}

// EventsForDay GET /events_for_day?user_id=1&date=2023-12-31
func (h *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error":"invalid date"}`, http.StatusBadRequest)
		return
	}

	events, err := h.svc.EventsForDay(r.Context(), parseUserID(userID), date)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// EventsForWeek GET /events_for_week
func (h *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error":"invalid date"}`, http.StatusBadRequest)
		return
	}

	events, err := h.svc.EventsForWeek(r.Context(), parseUserID(userID), date)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// EventsForMonth GET /events_for_month
func (h *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error":"invalid date"}`, http.StatusBadRequest)
		return
	}

	events, err := h.svc.EventsForMonth(r.Context(), parseUserID(userID), date)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// helper: JSON response
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// helper: convert user_id to int64
func parseUserID(s string) int64 {
	// простая реализация (для примера)
	// в реальном коде нужна обработка ошибок
	var id int64
	_, err := fmt.Sscan(s, &id)
	if err != nil {
		return 0
	}
	return id
}
