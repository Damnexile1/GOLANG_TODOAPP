package tasks_postgres_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	taskId int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	update todoapp.tasks
	set 
	    title = $1,
	    description = $2,
	    completed = $3,
	    completed_at = $4,
	    version = version + 1
	where id = $5 and version = $6
	returning 
		id,
		version,
		title,
		description,
		completed,
		created_at,
		completed_at,
		author_user_id
	;
    `

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
		task.Version,
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id = %d concurrently accessed: %w",
				taskId,
				core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("scan failed: %w", err)
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
