package statistics_service

import (
	"context"
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

type StatisticsService struct {
	StatisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(StatisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		StatisticsRepository: StatisticsRepository,
	}
}
