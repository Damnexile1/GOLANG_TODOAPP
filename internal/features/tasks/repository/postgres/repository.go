package tasks_postgres_postgres

import core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"

type TaskRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(pool core_postgres_pool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}
