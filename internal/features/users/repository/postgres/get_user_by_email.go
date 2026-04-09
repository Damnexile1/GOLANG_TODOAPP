package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	select id, version, full_name, phone_number, email, password_hash, role, manager_id
	from todoapp.users
	where email = $1;
	`

	row := r.pool.QueryRow(ctx, query, email)
	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.Role,
		&userModel.ManagerId,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with email=%s not found: %w", email, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	role, _ := domain.UserRoleFromInt(userModel.Role)
	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
		userModel.Email,
		userModel.PasswordHash,
		role,
		userModel.ManagerId,
	)
	return userDomain, nil
}
