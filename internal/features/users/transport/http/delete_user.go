package user_transport_http

import (
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	core_http_utils "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/utils"
)

func (h *UsersHttpHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")
	}
	err = h.usersService.DeleteUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
	}
	err = h.usersService.DeleteUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
	}

	responseHandler.NoContentResponse()
}
