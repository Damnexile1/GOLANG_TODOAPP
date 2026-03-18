package user_service

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

func (s *UserService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit cannot be less than 0: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset cannot be less than 0: %w", core_errors.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}
