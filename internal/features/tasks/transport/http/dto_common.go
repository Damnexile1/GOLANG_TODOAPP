package tasks_transport_http

import (
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id" example:"1"`
	Version      int        `json:"version" example:"1"`
	Title        string     `json:"title" example:"Buy groceries"`
	Description  *string    `json:"description" example:"Milk, eggs, bread"`
	Completed    bool       `json:"completed" example:"false"`
	Status       string     `json:"status" example:"created"`
	Deadline     *time.Time `json:"deadline" example:"2026-04-05T10:00:00Z"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-03-29T10:00:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"2026-03-29T12:00:00Z"`
	AuthorUserId int        `json:"author_user_id" example:"1"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		Status:       task.StatusKey.String(),
		Deadline:     task.Deadline,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserId: task.AuthorUserId,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))
	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}
	return dtos
}
