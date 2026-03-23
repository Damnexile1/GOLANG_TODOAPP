package user_service

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (s *UserService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetUser from repo: %w", err)
	}
	return user, nil
}
