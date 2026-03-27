package tasks_postgres_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTaskById(ctx context.Context, taskId int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	select id, version, title, description, completed, created_at, completed_at, author_user_id
	from todoapp.tasks
	where id = $1;
	`

	row := r.pool.QueryRow(ctx, query, taskId)
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id=%d not found: %w", taskId, core_errors.ErrNotFound)
		}
		return domain.Task{}, fmt.Errorf("scan row: %w", err)
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
