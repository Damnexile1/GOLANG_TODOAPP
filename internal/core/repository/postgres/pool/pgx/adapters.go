package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func (r pgxRow) Scan(dest ...any) error {

	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}
	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		if pgxErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesForeignKey)
		}
	}
	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}
