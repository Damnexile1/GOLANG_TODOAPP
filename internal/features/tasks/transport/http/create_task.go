package tasks_transport_http

import (
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Buy groceries"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Milk, eggs, bread"`
	AuthorUserId int     `json:"author_user_id" validate:"required,gte=-1" example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary Создание задачи
// @Description Создает новую задачу в системе
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Тело запроса для создания задачи"
// @Success 200 {object} CreateTaskResponse "Задача успешно создана"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks [post]
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

	responseHandler.JSONResponse(CreateTaskResponse(taskDTOFromDomain(taskDomain)), http.StatusOK)
}
