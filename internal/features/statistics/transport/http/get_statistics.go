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

type getStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

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

func toDtoFromDomain(statistics domain.Statistics) getStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return getStatisticsResponse{
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
