package dto

type CreateEventRequest struct {
	UserID int64  `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`   // YYYY-MM-DD
	Text   string `json:"event" form:"event"` // текст события
}

type UpdateEventRequest struct {
	ID     int64  `json:"id" form:"id"`
	UserID int64  `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
	Text   string `json:"event" form:"event"`
}

type DeleteEventRequest struct {
	ID int64 `json:"id" form:"id"`
}

// EventsQuery параметры для GET /events_for_*
type EventsQuery struct {
	UserID int64  `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
}
