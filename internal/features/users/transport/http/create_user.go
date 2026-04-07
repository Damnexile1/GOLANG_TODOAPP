package user_transport_http

import (
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name"    validate:"required,min=3,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79999999999"`
}

type CreateUserResponse UserDTOResponse

// CreateUser 	godoc
// @Summary 	Создание пользователя
// @Description СОздать нового пользователя в системе
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		request body 		CreateUserRequest true "CreateUser тело запроса"
// @Success 	201 	{object} 	CreateUserResponse "Успешно созданный пользователь"
// @Failure 	400 	{object} 	core_http_response.ErrorResponse "Bad Request"
// @Failure 	500 	{object} 	core_http_response.ErrorResponse "Internal Server Error"
// @Router 		/users [post]
func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoice CreateUser handler")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to validate create user request")

		return
	}
	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	// Временно создаем пользователя с дефолтными значениями
	// TODO: добавить email и password в CreateUserRequest после реализации auth
	return domain.NewUserUninitialized(
		dto.FullName,
		dto.PhoneNumber,
		"",                  // email - будет заполняться через auth
		"",                  // password_hash - будет заполняться через auth
		domain.UserRoleUser, // по умолчанию обычный пользователь
		nil,                 // manager_id
	)
}
