package user_service

import (
	"context"
	"fmt"
)

func (s *UserService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	err := s.usersRepository.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user from repo: %w", err)
	}
	return nil
}
