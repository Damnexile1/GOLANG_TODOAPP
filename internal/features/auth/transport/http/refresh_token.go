package auth_transport_http

import (
	"net/http"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_request "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/request"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken godoc
// @Summary Обновление токена
// @Description Обновляет access token используя refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse "Новые токены"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (h *AuthHTTPHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req RefreshTokenRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode refresh token request")
		return
	}

	tokens, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	response := RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	responseHandler.JSONResponse(response, http.StatusOK)
}
