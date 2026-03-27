package tasks_postgres_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, taskId int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `delete from todoapp.tasks where id = $1;`
	cmdTag, err := r.pool.Exec(ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("delete task: no such task with id %d : %w", taskId, core_errors.ErrNotFound)
	}

	return nil
}
