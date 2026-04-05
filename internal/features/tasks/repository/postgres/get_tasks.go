package tasks_postgres_postgres

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userId *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	select id, version, title, description, completed, status_key, deadline, created_at, completed_at, author_user_id
	from todoapp.tasks
	%s
	order by id asc
	limit $1
	offset $2;
	`

	args := []any{limit, offset}
	if userId != nil {
		query = fmt.Sprintf(query, "where author_user_id = $3")
		args = append(args, *userId)
	} else {
		query = fmt.Sprintf(query, "")
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	defer rows.Close()
	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := taskDomainsFromModels(taskModels)
	return taskDomains, nil
}
