package user_transport_http

import (
	"net/http"

	core_http_server "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/server"
)

type UsersHttpHandler struct {
	usersService UserService
}

type UserService interface {
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
	}
}
