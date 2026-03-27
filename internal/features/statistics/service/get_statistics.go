package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("from must be before to %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.StatisticsRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("error getting tasks: %w", err)
	}
	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(
			0,
			0,
			nil,
			nil,
		)
	}

	tasksCreated := len(tasks)
	tasksCompleted := 0
	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}
		completedDuration := task.CompletionDuration()
		if completedDuration != nil {
			totalCompletedDuration += *completedDuration
		}
	}

	TasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration

	if tasksCompleted > 0 && totalCompletedDuration > 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&TasksCompletedRate,
		tasksAverageCompletionTime,
	)
}
