package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_errors "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/errors"
	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	userId int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	update todoapp.users 
	set 
		full_name=$1,
		phone_number=$2, 
		version=version+1 
	where id=$3 and version=$4 
	returning 
		id, 
		version, 
		full_name, 
		phone_number;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		userId,
		user.Version,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id =`%d` concurrently accessed: %w",
				userId,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("scan failed: %w", err)
	}
	userDomain := domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		FullName:    userModel.FullName,
		PhoneNumber: userModel.PhoneNumber,
	}

	return userDomain, nil

}
