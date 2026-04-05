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
	StatusKey   TaskStatus
	Deadline    *time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserId int
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
	Deadline    Nullable[time.Time]
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	statusKey TaskStatus,
	deadline *time.Time,
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
		StatusKey:    statusKey,
		Deadline:     deadline,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserId: authorUserId,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	deadline *time.Time,
	authorUserId int,
) Task {
	return Task{
		ID:           UninitializedId,
		Version:      UninitializedVersion,
		Title:        title,
		Description:  description,
		Completed:    false,
		StatusKey:    TaskStatusCreated,
		Deadline:     deadline,
		CreatedAt:    time.Now(),
		CompletedAt:  nil,
		AuthorUserId: authorUserId,
	}
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
	deadline Nullable[time.Time],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
		Deadline:    deadline,
	}
}

func (t *Task) UpdateStatus() {
	if t.StatusKey == TaskStatusCompleted || t.StatusKey == TaskStatusFailed {
		return
	}
	if t.Deadline != nil && time.Now().After(*t.Deadline) && !t.Completed {
		t.StatusKey = TaskStatusFailed
		return
	}
	if t.Completed {
		t.StatusKey = TaskStatusCompleted
	} else {
		t.StatusKey = TaskStatusCreated
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

	if !t.StatusKey.IsValid() {
		return fmt.Errorf("invalid status key %d :%w", t.StatusKey, core_errors.ErrInvalidArgument)
	}

	if t.Deadline != nil && t.Deadline.Before(t.CreatedAt) {
		return fmt.Errorf("deadline is earlier than CreatedAt :%w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (tp *TaskPatch) Validate() error {
	if tp.Title.Set && tp.Title.Value == nil {
		return fmt.Errorf("title cant be patched to null:%w", core_errors.ErrInvalidArgument)
	}
	if tp.Completed.Set && tp.Completed.Value == nil {
		return fmt.Errorf("completed cant be patched to null:%w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := t.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}
	tmp := *t
	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	if patch.Deadline.Set {
		tmp.Deadline = patch.Deadline.Value
	}
	if patch.Completed.Set {
		completed := *patch.Completed.Value
		if completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
			tmp.StatusKey = TaskStatusCompleted
		} else {
			tmp.CompletedAt = nil
			tmp.StatusKey = TaskStatusCreated
		}
		tmp.Completed = completed
	}
	tmp.UpdateStatus()

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}

	*t = tmp

	return nil
}
