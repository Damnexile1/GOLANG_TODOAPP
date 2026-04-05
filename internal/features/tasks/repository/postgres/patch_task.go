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
	    status_key = $4,
	    deadline = $5,
	    completed_at = $6,
	    version = version + 1
	where id = $7 and version = $8
	returning 
		id,
		version,
		title,
		description,
		completed,
		status_key,
		deadline,
		created_at,
		completed_at,
		author_user_id
	;
    `

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		int(task.StatusKey),
		task.Deadline,
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
		&taskModel.StatusKey,
		&taskModel.Deadline,
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
	taskDomain := taskDomainFromModel(taskModel)
	return taskDomain, nil
}
