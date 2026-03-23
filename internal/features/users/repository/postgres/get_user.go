package users_postgres_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userId int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	select id, version, full_name, phone_number
	from todoapp.users
	where id = $1;
	`

	row := r.pool.QueryRow(ctx, query, userId)
	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d not found: %w", userId, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
	return userDomain, nil
}
