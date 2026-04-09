package users_postgres_repository

import "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"

type UserModel struct {
	ID           int
	Version      int
	FullName     string
	PhoneNumber  *string
	Email        string
	PasswordHash string
	Role         int
	ManagerId    *int
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		role, _ := domain.UserRoleFromInt(user.Role)
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
			user.Email,
			user.PasswordHash,
			role,
			user.ManagerId,
		)
	}

	return userDomains
}
