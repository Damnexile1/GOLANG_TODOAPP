package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	select id, version, title, description, completed, created_at, completed_at, author_user_id
	from todoapp.tasks
    `)

	var args []any
	var conditions []string

	if userId != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", len(args)+1))
		args = append(args, *userId)
	}
	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, *from)
	}
	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", len(args)+1))
		args = append(args, *to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" where " + strings.Join(conditions, " and "))
	}
	queryBuilder.WriteString(" order by id asc;")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)

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
		return nil, fmt.Errorf("rows: %w", err)
	}
	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
