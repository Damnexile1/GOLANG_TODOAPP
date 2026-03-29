package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Получение задачи
// @Description Возвращает задачу по идентификатору
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} GetTaskResponse "Задача успешно получена"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId path value")
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
