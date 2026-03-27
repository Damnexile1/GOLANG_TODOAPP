package tasks_postgres_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) Create(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	insert into todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
	values ($1, $2, $3, $4, $5, $6)
	returning id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserId,
	)
	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserId,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf("%v: user with id=%d: %w",
				err, taskModel.AuthorUserId, core_errors.ErrNotFound)
		}
		return domain.Task{}, fmt.Errorf("scan error %w", err)
	}

	taskDomain := domain.Task{
		ID:           taskModel.ID,
		Version:      taskModel.Version,
		Title:        taskModel.Title,
		Description:  taskModel.Description,
		Completed:    taskModel.Completed,
		CreatedAt:    taskModel.CreatedAt,
		CompletedAt:  taskModel.CompletedAt,
		AuthorUserId: taskModel.AuthorUserId,
	}

	return taskDomain, nil
}
