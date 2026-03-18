package user_service

import (
	"context"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

type UserService struct {
	usersRepository UserRepository
}

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
}

func NewUsersService(usersRepository UserRepository) *UserService {
	return &UserService{
		usersRepository: usersRepository,
	}
}
