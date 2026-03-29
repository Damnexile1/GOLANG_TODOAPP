package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type queryParams struct {
	userId *int
	from   *time.Time
	to     *time.Time
}

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"10"`
	TasksCompleted             int      `json:"tasks_completed" example:"7"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"0.7"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"16h10m9.862461s"`
}

// GetStatistics godoc
// @Summary Получение статистики
// @Description Возвращает статистику по задачам с возможностью фильтрации по пользователю и диапазону дат
// @Tags statistics
// @Produce json
// @Param user_id query int false "ID пользователя для фильтрации статистики" example(1)
// @Param from query string false "Начальная дата периода" example(2026-03-01)
// @Param to query string false "Конечная дата периода" example(2026-03-29)
// @Success 200 {object} GetStatisticsResponse "Статистика успешно получена"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)
	var params queryParams
	params, err := getUserIdFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "error getting user_id/from/to query params")
		return
	}
	domainStatistics, err := h.statisticsService.GetStatistics(ctx, params.userId, params.from, params.to)
	if err != nil {
		responseHandler.ErrorResponse(err, "error getting statistics")

		return
	}
	response := toDtoFromDomain(domainStatistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDtoFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIdFromToQueryParams(req *http.Request) (queryParams, error) {
	const (
		useridQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userId, err := core_http_request.GetIntQueryParam(req, useridQueryParamKey)
	if err != nil {
		return queryParams{nil, nil, nil}, fmt.Errorf("get 'userId' parameter: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(req, fromQueryParamKey)
	if err != nil {
		return queryParams{nil, nil, nil}, fmt.Errorf("get 'from' parameter: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(req, toQueryParamKey)
	if err != nil {
		return queryParams{nil, nil, nil}, fmt.Errorf("get 'to' parameter: %w", err)
	}

	return queryParams{userId, from, to}, nil
}
