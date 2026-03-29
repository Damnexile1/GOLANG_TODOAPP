package user_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers godoc
// @Summary Получение списка пользователей
// @Description Возвращает список пользователей с поддержкой пагинации через limit и offset
// @Tags users
// @Produce json
// @Param limit query int false "Лимит количества пользователей" example(10)
// @Param offset query int false "Смещение для пагинации" example(0)
// @Success 200 {array} UserDTOResponse "Список пользователей успешно получен"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /users [get]
func (h *UsersHttpHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetParams(req)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to parse query parameters",
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetParams(req *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_request.GetIntQueryParam(req, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' parameter: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(req, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' parameter: %w", err)
	}

	return limit, offset, nil
}
