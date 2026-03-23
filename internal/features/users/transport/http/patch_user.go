package user_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    types.Nullable[string] `json:"full_name"`
	PhoneNumber types.Nullable[string] `json:"phone_number"`
}

func (h *UsersHttpHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request PatchUserRequest
	err := core_http_request.DecodeAndValidateRequest(r, &request)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	log.Debug(
		fmt.Sprintf(
			"PathUserRequest fields: \nFullName: '%v' \nPhoneNumber: '%v' \n",
			request.FullName,
			request.PhoneNumber,
		),
	)

	w.WriteHeader(http.StatusOK)

}
