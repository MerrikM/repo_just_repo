package model

import (
	"time"
)

// Event Domain-модель (используется в сервисном слое и при доступе к БД через sqlx)
type Event struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Date      time.Time `db:"date" json:"date"`
	Text      string    `db:"text" json:"event"`
	Active    bool      `db:"active" json:"active"` // true — событие активноs
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
