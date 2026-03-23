package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	userId int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `delete from todoapp.users where id = $1;`

	cmdTag, err := r.pool.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("delete user: no such user with id %d : %w", userId, core_errors.ErrNotFound)
	}

	return nil
}
