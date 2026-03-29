package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Получение списка задач
// @Description Возвращает список задач с возможностью фильтрации по user_id и пагинацией через limit и offset
// @Tags tasks
// @Produce json
// @Param user_id query int false "ID пользователя для фильтрации задач" example(1)
// @Param limit query int false "Лимит количества задач" example(10)
// @Param offset query int false "Смещение для пагинации" example(0)
// @Success 200 {array} TaskDTOResponse "Список задач успешно получен"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)
	userId, limit, offset, err := getUserIdLimitOffsetParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to parse user_id/limit/offset params",
		)
		return
	}
	domainTasks, err := h.tasksService.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(domainTasks))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIdLimitOffsetParams(req *http.Request) (*int, *int, *int, error) {
	const (
		useridQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	userId, err := core_http_request.GetIntQueryParam(req, useridQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userId' parameter: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(req, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' parameter: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(req, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' parameter: %w", err)
	}

	return userId, limit, offset, nil
}
