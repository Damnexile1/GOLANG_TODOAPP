package users_postgres_repository

import (
	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
)

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(pool core_postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
