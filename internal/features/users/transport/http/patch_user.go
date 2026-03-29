package user_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    types.Nullable[string] `json:"full_name"`
	PhoneNumber types.Nullable[string] `json:"phone_number"`
}

type PatchUserResponse UserDTOResponse

// Validate PatchUser godoc
// @Summary Частичное обновление пользователя
// @Description Частично обновляет данные пользователя по идентификатору
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body PatchUserRequest true "Тело запроса для частичного обновления пользователя"
// @Success 200 {object} PatchUserResponse "Пользователь успешно обновлен"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /users/{id} [patch]
func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("required field 'full_name' is required")
		}
	}
	fullNameLen := len([]rune(*r.FullName.Value))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf("required field 'full_name' is between 3 and 100")
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("required field 'phone_number' is between 10 and 15")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("required field 'phone_number' must start with '+'")
			}
		}
	}

	return nil
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

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(r PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		r.FullName.ToDomain(),
		r.PhoneNumber.ToDomain(),
	)
}
