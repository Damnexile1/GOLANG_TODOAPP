package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/auth/jwt"
	core_http_server "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	Register(ctx context.Context, email, password, fullName string, phoneNumber *string) (*jwt.TokenPair, error)
	Login(ctx context.Context, email, password string) (*jwt.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error)
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.RefreshToken,
		},
	}
}
