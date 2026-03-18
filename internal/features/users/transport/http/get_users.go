package user_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	core_http_utils "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHttpHandler) getUsers(rw http.ResponseWriter, req *http.Request) {
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
	limit, err := core_http_utils.GetIntQueryParam(req, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' parameter: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(req, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' parameter: %w", err)
	}

	return limit, offset, nil
}
