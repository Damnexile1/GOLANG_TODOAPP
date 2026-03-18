package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	select id, version, full_name, phone_number
	from todoapp.users
	order by id asc
	limit $1
	offset $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		userModels = append(userModels, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)
	return userDomains, nil
}
