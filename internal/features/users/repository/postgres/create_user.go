package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number, email, password_hash, role, manager_id) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, full_name, phone_number, email, password_hash, role, manager_id;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, user.Email, user.PasswordHash, int(user.Role), user.ManagerId)

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
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	role, _ := domain.UserRoleFromInt(userModel.Role)
	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
		userModel.Email,
		userModel.PasswordHash,
		role,
		userModel.ManagerId,
	), nil
}
