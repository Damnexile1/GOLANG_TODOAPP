package tasks_transport_http

import (
	"net/http"
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserId int     `json:"author_user_id" validate:"required,gte=-1"`
}

type CreateTaskResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserId int        `json:"author_user_id"`
}

func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	var req CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode create task request",
		)

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		req.Title,
		req.Description,
		req.AuthorUserId,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	responseHandler.JSONResponse(taskDTOFromDomain(taskDomain), http.StatusOK)
}

func taskDTOFromDomain(task domain.Task) CreateTaskResponse {
	return CreateTaskResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserId: task.AuthorUserId,
	}
}
