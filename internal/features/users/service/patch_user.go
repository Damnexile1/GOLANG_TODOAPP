package user_service

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (s *UserService) PatchUser(
	ctx context.Context,
	userId int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, userId)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("failed to apply patch: %w", err)
	}
	patchedUser, err := s.usersRepository.PatchUser(ctx, userId, user)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to patch user: %w", err)
	}
	return patchedUser, nil
}
