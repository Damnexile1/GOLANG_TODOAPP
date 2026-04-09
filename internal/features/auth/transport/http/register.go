package auth_transport_http

import (
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type RegisterRequest struct {
	Email       string  `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string  `json:"password" validate:"required,min=6,max=100" example:"password123"`
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79999999999"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя и возвращает JWT токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} RegisterResponse "Успешная регистрация"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode register request")
		return
	}

	tokens, err := h.authService.Register(ctx, req.Email, req.Password, req.FullName, req.PhoneNumber)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register user")
		return
	}

	response := RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	responseHandler.JSONResponse(response, http.StatusCreated)
}
