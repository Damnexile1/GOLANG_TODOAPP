package statistics_postgres_repository

import (
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserId int
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserId,
	)
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))
	for i, t := range taskModels {
		domains[i] = domain.Task{
			ID:           t.ID,
			Version:      t.Version,
			Title:        t.Title,
			Description:  t.Description,
			Completed:    t.Completed,
			CreatedAt:    t.CreatedAt,
			CompletedAt:  t.CompletedAt,
			AuthorUserId: t.AuthorUserId,
		}
	}

	return domains
}
