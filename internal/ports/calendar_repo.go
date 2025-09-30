package ports

import (
	"context"
	"repo_just_repo/internal/model"
	"time"
)

type CalendarRepository interface {
	Create(ctx context.Context, e *model.Event) (int64, error)
	Update(ctx context.Context, e *model.Event) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.Event, error)
	FindForRange(ctx context.Context, userID int64, from, to time.Time) ([]model.Event, error)
}
