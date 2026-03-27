package tasks_service

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, taskId int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTaskById(ctx, taskId)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error getting task: %w", err)
	}
	return task, nil
}
