package domain

import (
	"fmt"
	"time"

	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserId int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserId int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserId: authorUserId,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserId int,
) Task {
	return Task{
		UninitializedId,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserId,
	}
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("invalid title length %d :%w", titleLen, core_errors.ErrInvalidArgument)
	}
	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("invalid description length %d :%w", descriptionLen, core_errors.ErrInvalidArgument)
		}
	}
	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("CompletedAt cant be null when completed=true :%w", core_errors.ErrInvalidArgument)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("CompletedAt is earlier than CreatedAt :%w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("CompletedAt must be null when completed=false :%w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}
