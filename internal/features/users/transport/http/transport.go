package user_transport_http

import (
	"context"
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_http_server "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/server"
)

type UsersHttpHandler struct {
	usersService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit, offset *int,
	) ([]domain.User, error)
}

func NewUsersHTTPHandler(usersService UserService) *UsersHttpHandler {
	return &UsersHttpHandler{
		usersService: usersService,
	}
}

func (h *UsersHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.getUsers,
		},
	}
}
