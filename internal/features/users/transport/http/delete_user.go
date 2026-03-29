package user_transport_http

import (
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаляет пользователя из системы по идентификатору
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 204 "Пользователь успешно удален"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /users/{id} [delete]
func (h *UsersHttpHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")
	}
	err = h.usersService.DeleteUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
	}

	responseHandler.NoContentResponse()
}
