package repository

import (
	"context"
	"database/sql"
	"repo_just_repo/internal/config"
	"repo_just_repo/internal/model"
	"time"
)

type CalendarRepository struct {
	*config.Database
}

func NewCalendarRepository(db *config.Database) *CalendarRepository {
	return &CalendarRepository{db}
}

func (r *CalendarRepository) Create(ctx context.Context, e *model.Event) (int64, error) {
	query := `
		INSERT INTO events (user_id, date, text, active, created_at, updated_at)
		VALUES ($1, $2, $3, TRUE, NOW(), NOW())
		RETURNING id
	`
	var id int64
	err := r.QueryRowContext(ctx, query, e.UserID, e.Date, e.Text).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CalendarRepository) Update(ctx context.Context, e *model.Event) error {
	query := `
		UPDATE events
		SET user_id=$1, date=$2, text=$3, updated_at=NOW()
		WHERE id=$4 AND active=TRUE
	`
	res, err := r.ExecContext(ctx, query, e.UserID, e.Date, e.Text, e.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CalendarRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE events
		SET active=FALSE, updated_at=NOW()
		WHERE id=$1 AND active=TRUE
	`
	res, err := r.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CalendarRepository) FindByID(ctx context.Context, id int64) (*model.Event, error) {
	query := `
		SELECT id, user_id, date, text, active, created_at, updated_at
		FROM events
		WHERE id=$1 AND active=TRUE
	`
	var e model.Event
	err := r.GetContext(ctx, &e, query, id)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *CalendarRepository) FindForRange(ctx context.Context, userID int64, from, to time.Time) ([]model.Event, error) {
	query := `
		SELECT id, user_id, date, text, active, created_at, updated_at
		FROM events
		WHERE user_id=$1
		  AND date >= $2
		  AND date <= $3
		  AND active=TRUE
		ORDER BY date ASC
	`
	var events []model.Event
	err := r.SelectContext(ctx, &events, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	return events, nil
}
